package rpc

import (
	"context"
	_ "embed"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

type Client struct {
	requests  int
	rateLimit int
	prevStop  time.Time
	client    *ethclient.Client
}

var mpContract = common.HexToAddress("0x213073989821f738a7ba3520c3d31a1f9ad31bbd")
var saleTopic = common.HexToHash("0xae3392a96856e8c1881402157f65e69336cb9e04ffba578babad5b29909def82")

func New() *Client {
	client, err := ethclient.Dial("https://api.roninchain.com/rpc")
	if err != nil {
		log.Fatal(err)
	}
	return &Client{rateLimit: 1000, prevStop: time.Now(), client: client}
}

func (rpc *Client) Close() {
	rpc.client.Close()
}

func (rpc *Client) GetClient() *ethclient.Client {
	if rpc.requests < rpc.rateLimit {
		rpc.requests++
		return rpc.client
	}
	currentTime := time.Now()
	log.Println("RPC limit reached")
	time.Sleep(rpc.prevStop.Add(5 * time.Minute).Sub(currentTime))
	log.Println("RPC resuming")
	rpc.prevStop = currentTime
	rpc.requests = 0
	return rpc.client
}

func GetBlocksFromLogs(logs []types.Log) []uint64 {
	var blocks []uint64
	for _, lg := range logs {
		blocks = append(blocks, lg.BlockNumber)
	}
	return blocks
}

func GetBlocksTimestamp(rpc *Client, blockNumbers []uint64) map[uint64]time.Time {
	// Map to keep track of block numbers recorded
	blocks := make(map[uint64]time.Time)
	for i, blockNumber := range blockNumbers {
		// Check if block number is already recorded
		if _, found := blocks[blockNumber]; !found {
			// Get the block data by its number
			block, err := rpc.GetClient().BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
			if err != nil {
				log.Fatal(err)
			}
			blocks[blockNumber] = time.Unix(int64(block.Time()), 0)
		}
		log.Println("Fetched", i, "out of", len(blockNumbers), " block timestamps")
	}
	return blocks
}

func GetLatestBlockNumber(rpc *Client) uint64 {
	blockNumber, err := rpc.GetClient().BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return blockNumber
}

func GetFilters(start uint64, end uint64) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(start)),
		ToBlock:   big.NewInt(int64(end)),
		Addresses: []common.Address{mpContract},
		Topics:    [][]common.Hash{{saleTopic}},
	}
}

func GetLogs(rpc *Client, filter ethereum.FilterQuery) []types.Log {
	logs, err := rpc.GetClient().FilterLogs(context.Background(), filter)
	if err != nil {
		log.Fatal("GetLogs:", err)
	}
	return logs
}

func GetSalesFromLogs(rpc *Client, blocksTimestamp map[uint64]time.Time, logs []types.Log) []Sale {
	var sales []Sale
	for _, lg := range logs {
		tx, _, err := rpc.GetClient().TransactionByHash(context.Background(), lg.TxHash)
		if err != nil {
			log.Println("GetSalesFromLogs:", err)
		}
		tkId, sale := GetSaleFromTx(tx)
		if tkId != "0x00000000000000000000000032950db2a7164ae833121501c797d79e7b79d74c" {
			continue
		}
		sale.TxHash = lg.TxHash.String()
		sale.Timestamp = blocksTimestamp[lg.BlockNumber]
		sales = append(sales, sale)
	}
	return sales
}

func GetSaleFromTx(tx *types.Transaction) (string, Sale) {
	tokenID := hexutil.Encode(tx.Data()[324:324+32])
	id := big.NewInt(0).SetBytes(tx.Data()[388 : 388+32]).Uint64()
	startPrice, _ := WeiToEther(big.NewInt(0).SetBytes(tx.Data()[452 : 452+32])).Float64()
	endPrice, _ := WeiToEther(big.NewInt(0).SetBytes(tx.Data()[516 : 516+32])).Float64()
	return tokenID, Sale{ID: int(id), StartPrice: startPrice, EndPrice: endPrice}
}

type Sale struct {
	TxHash     string
	ID         int
	StartPrice float64
	EndPrice   float64
	Timestamp  time.Time
}

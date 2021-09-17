package main

import (
	"axie-sales-sub/firebase"
	"axie-sales-sub/graphql"
	"axie-sales-sub/rpc"
	"github.com/shanemaglangit/agp"
	"log"
	"time"
)

func main() {
	rpcClient := rpc.New()
	defer rpcClient.Close()
	fb := firebase.New()
	gql := graphql.New()

	startBlock := rpc.GetLatestBlockNumber(rpcClient)
	endBlock := rpc.GetLatestBlockNumber(rpcClient)

	for startBlock <= endBlock {
		log.Println("Processing block", startBlock, "to", endBlock)
		filter := rpc.GetFilters(startBlock, endBlock)
		logs := rpc.GetLogs(rpcClient, filter)

		log.Println("Fetched", len(logs), "logs")
		blocks := rpc.GetBlocksFromLogs(logs)

		log.Println("Fetched", len(blocks), "blocks")
		blockTimestamp := rpc.GetBlocksTimestamp(rpcClient, blocks)

		log.Println("Fetched", len(blockTimestamp), "blockTimestamp")
		sales := rpc.GetSalesFromLogs(rpcClient, blockTimestamp, logs)

		log.Println("Fetched", len(sales), "sales")
		ids := GetIDSlice(sales)

		log.Println("Extracted", len(ids), "ids")
		axieMap, err := gql.GetAxiesInfo(ids)
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}

		log.Println("Mapped", len(axieMap), "axies")
		firebaseSales := ParseToFirebaseSales(sales, axieMap)

		log.Println("Parsed", len(firebaseSales), "firebase sales")
		fb.SaveSalesMultiple(firebaseSales)
		startBlock = endBlock
		endBlock = rpc.GetLatestBlockNumber(rpcClient)
		for startBlock == endBlock {
			log.Println("Reached end: Sleep for 5s")
			time.Sleep(time.Second * 5)
			endBlock = rpc.GetLatestBlockNumber(rpcClient)
		}
	}
}

func ParseToFirebaseSales(sales []rpc.Sale, axies map[int]graphql.Axie) []firebase.Sale {
	var firebaseSales []firebase.Sale
	for _, sale := range sales {
		axie := axies[sale.ID]
		firebaseSale := firebase.Sale{
			TxHash:     sale.TxHash,
			ID:         sale.ID,
			StartPrice: sale.StartPrice,
			EndPrice:   sale.EndPrice,
			Timestamp:  sale.Timestamp,
		}

		if axie.Stage != 4 {
			firebaseSales = append(firebaseSales, firebaseSale)
			continue
		}

		genes, err := agp.ParseHexDecode(axie.Genes)
		if err != nil {
			firebaseSales = append(firebaseSales, firebaseSale)
			continue
		}

		firebaseSale.Class = axie.Class
		firebaseSale.Purity = axie.Pureness
		firebaseSale.NumMystic = axie.NumMystic
		firebaseSale.BreedCount = axie.BreedCount

		firebaseSale.HP = axie.Stats.HP
		firebaseSale.Speed = axie.Stats.Speed
		firebaseSale.Skill = axie.Stats.Skill
		firebaseSale.Morale = axie.Stats.Morale

		firebaseSale.EyesD = genes.Eyes.D.PartId
		firebaseSale.EyesR1 = genes.Eyes.R1.PartId
		firebaseSale.EyesR2 = genes.Eyes.R2.PartId

		firebaseSale.EarsD = genes.Ears.D.PartId
		firebaseSale.EarsR1 = genes.Ears.R1.PartId
		firebaseSale.EarsR2 = genes.Ears.R2.PartId

		firebaseSale.HornD = genes.Horn.D.PartId
		firebaseSale.HornR1 = genes.Horn.R1.PartId
		firebaseSale.HornR2 = genes.Horn.R2.PartId

		firebaseSale.MouthD = genes.Mouth.D.PartId
		firebaseSale.MouthR1 = genes.Mouth.R1.PartId
		firebaseSale.MouthR2 = genes.Mouth.R2.PartId

		firebaseSale.BackD = genes.Back.D.PartId
		firebaseSale.BackR1 = genes.Back.R1.PartId
		firebaseSale.BackR2 = genes.Back.R2.PartId

		firebaseSale.TailD = genes.Tail.D.PartId
		firebaseSale.TailR1 = genes.Tail.R1.PartId
		firebaseSale.TailR2 = genes.Tail.R2.PartId

		firebaseSale.GeneQuality = genes.GeneQuality

		firebaseSales = append(firebaseSales, firebaseSale)
	}
	return firebaseSales
}

func GetIDSlice(sales []rpc.Sale) []int {
	var ids []int
	for _, sale := range sales {
		ids = append(ids, sale.ID)
	}
	return ids
}

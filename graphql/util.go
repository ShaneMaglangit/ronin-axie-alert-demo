package graphql

import (
	"math/big"
)

func AxieAliasMapToSlice(axieMap map[int]Axie, axieRes AxieRes) map[int]Axie {
	for _, axie := range axieRes.Data {
		id, _ := big.NewInt(0).SetString(axie.ID, 10)
		axieMap[int(id.Uint64())] = axie
	}
	return axieMap
}
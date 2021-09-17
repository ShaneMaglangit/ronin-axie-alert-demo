package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type GraphQL struct {
	requests  int
	rateLimit int
	prevStop  time.Time
}

func New() *GraphQL {
	return &GraphQL{rateLimit: 1000, prevStop: time.Now()}
}

func (graphql *GraphQL) getClient() *http.Client {
	if graphql.requests < graphql.rateLimit {
		graphql.requests++
		return &http.Client{Timeout: 30 * time.Second}
	}
	currentTime := time.Now()
	log.Println("GraphQL Limit Reached")
	time.Sleep(graphql.prevStop.Add(5 * time.Minute).Sub(currentTime))
	log.Println("GraphQL Limit Resuming")
	graphql.prevStop = currentTime
	graphql.requests = 0
	return &http.Client{Timeout: 15 * time.Second}
}

func (graphql *GraphQL) GetAxiesInfo(ids []int) (map[int]Axie, error) {
	axieMap := make(map[int]Axie)
	endpoint := "https://graphql-gateway.axieinfinity.com/graphql"

	for i := 0; i < len(ids); i += 30 {
		end := i + 30
		if end > len(ids) {
			end = len(ids)
		}

		sub := ids[i:end]

		query := "query { "
		for _, id := range sub {
			query += fmt.Sprintf(" axie%d: axie(axieId: %d) { ...AxieDetail } ", id, id)
		}
		query += "} fragment AxieDetail on Axie { id birthDate stage breedCount matronId sireId class pureness numMystic stats { hp speed skill morale } genes }"
		body, _ := json.Marshal(map[string]interface{}{"query": query})

		resp, err := graphql.getClient().Post(endpoint, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return axieMap, err
		}

		var axieRes AxieRes
		if err = json.NewDecoder(resp.Body).Decode(&axieRes); err != nil {
			return axieMap, err
		}
		axieMap = AxieAliasMapToSlice(axieMap, axieRes)
	}

	return axieMap, nil
}

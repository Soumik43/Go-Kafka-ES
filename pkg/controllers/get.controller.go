package controller

import (
	"bytes"
	"encoding/json"
	config "go-kafka-es/pkg/config"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	es, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{config.GetConfig("ES_ADDRESS")},
			Username:  config.GetConfig("ES_USERNAME"),
			Password:  config.GetConfig("ES_PASSWORD"),
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.Search(
		es.Search.WithContext(r.Context()),
		es.Search.WithIndex(config.GetConfig("ES_INDEX")),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var resMap map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(resMap["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(resMap["took"].(float64)),
	)
	for _, hit := range resMap["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

}

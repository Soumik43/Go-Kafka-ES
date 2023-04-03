package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	config "go-kafka-es/pkg/config"
	"log"
	"net/http"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/segmentio/kafka-go"
)

type ElasticMessage struct {
	Message string `json:"message"`
}

func WriteHandler(w http.ResponseWriter, r *http.Request) {

	es, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{config.GetConfig("ES_ADDRESS")},
			Username:  config.GetConfig("ES_USERNAME"),
			Password:  config.GetConfig("ES_PASSWORD"),
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v\n", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.GetConfig("KAFKA_ADDRESS")},
		Topic:   config.GetConfig("KAFKA_TOPIC"),
	})
	defer reader.Close()

	idIndex := 1
	for {
		msg, err := reader.ReadMessage(r.Context())
		if err != nil {
			log.Fatalf("Failed to read message: %v\n", err)
		}

		log.Printf("Received message: %s\n", string(msg.Value))

		kafkaMessage := ElasticMessage{
			Message: string(msg.Value),
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(kafkaMessage); err != nil {
			log.Fatalf("Failed to encode message: %v\n", err)
		}

		req := esapi.IndexRequest{
			Index:      config.GetConfig("ES_INDEX"),
			DocumentID: strconv.Itoa(idIndex),
			Body:       &buf,
			Refresh:    "true",
		}

		res, err := req.Do(r.Context(), es)
		if err != nil {
			log.Fatalf("Failed to perform request 1: %v\n", err)
			continue
		}
		if res.IsError() {
			log.Fatalf("Failed to perform request 2: %v\n", res.String())
			continue
		}
		res.Body.Close()

		fmt.Printf("Indexed message with ID %d\n", idIndex)

		// Increment index for ID (test)
		// Logic will fail for concurrent requests (maybe atomic increment is a better solution)
		idIndex++
	}
}

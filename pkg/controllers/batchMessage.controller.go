package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	config "go-kafka-es/pkg/config"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/segmentio/kafka-go"
)

type JsonStruct struct {
	MessagesToFetch []string `json:"messagesToFetch"`
}

type BatchIndividualMessage struct {
	Message string `json:"message"`
}

func batchAPIReq(bi esutil.BulkIndexer, messages *[]string, messageIndex *int, countSuccessful *uint64) {
	batchMessages, batchMessagesErr := json.Marshal(JsonStruct{MessagesToFetch: *messages})
	if batchMessagesErr != nil {
		log.Printf("Failed to marshal batch messages: %v\n", batchMessagesErr)
	}
	biErr := bi.Add(
		context.Background(),
		esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: strconv.Itoa(*messageIndex),
			Body:       bytes.NewReader(batchMessages),
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				atomic.AddUint64(countSuccessful, 1)
				fmt.Println("Successful batch", *countSuccessful)
			},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: %s", err)
				} else {
					log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
				}
			},
		},
	)
	if biErr != nil {
		log.Printf("Failed to add message: %v\n", biErr)
	}
}

func BatchMessageHandler(w http.ResponseWriter, r *http.Request) {

	es, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses:     []string{config.GetConfig("ES_ADDRESS")},
			Username:      config.GetConfig("ES_USERNAME"),
			Password:      config.GetConfig("ES_PASSWORD"),
			RetryOnStatus: []int{502, 503, 504, 429},
			MaxRetries:    5,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v\n", err)
	}

	conn, err := kafka.DialLeader(r.Context(), "tcp", config.GetConfig("KAFKA_ADDRESS"), config.GetConfig("KAFKA_TOPIC"), 0)
	if err != nil {
		log.Printf("Error connecting to Kafka: %v", err)
	}
	defer conn.Close()

	bi, err := esutil.NewBulkIndexer(
		esutil.BulkIndexerConfig{
			Index:  config.GetConfig("ES_INDEX"),
			Client: es,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create BulkIndexer: %v\n", err)
	}
	defer bi.Close(r.Context())

	messages := make([]string, 100)
	// Create a batch
	// TODO: future => how to reduce batch size
	batch := conn.ReadBatch(1, 10e6)
	messageIndex := 1
	var countSuccessful uint64
	for {
		m, err := batch.ReadMessage()
		if err != nil {
			log.Printf("error reading message %v\n", err)
			break
		}
		batchMessage := BatchIndividualMessage{
			Message: string(m.Value),
		}
		message, _ := json.Marshal(batchMessage)
		biErr := bi.Add(
			r.Context(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.Itoa(messageIndex),
				Body:       bytes.NewReader(message),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
					fmt.Println("Successful batch", countSuccessful)
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if biErr != nil {
			log.Printf("Failed to add message: %v\n", biErr)
		}
		messageIndex++
	}
	if batch.Offset() < 100 {
		batchAPIReq(bi, &messages, &messageIndex, &countSuccessful)
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}

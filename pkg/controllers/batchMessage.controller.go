package controller

import (
	"fmt"
	config "go-kafka-es/pkg/config"
	"log"
	"net/http"

	"github.com/segmentio/kafka-go"
)

func BatchMessageHandler(w http.ResponseWriter, r *http.Request) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.GetConfig("KAFKA_ADDRESS")},
		Topic:    config.GetConfig("KAFKA_TOPIC"),
		MinBytes: 1e6,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	conn, err := kafka.DialLeader(r.Context(), "tcp", config.GetConfig("KAFKA_ADDRESS"), config.GetConfig("KAFKA_TOPIC"), 0)
	if err != nil {
		log.Printf("Error connecting to Kafka: %v", err)
	}
	batch := conn.ReadBatch(1e3, 10e6)
	for {
		m, err := batch.ReadMessage()
		if err != nil {
			log.Printf("error reading message %v", err)
			break
		}
		fmt.Println(string(m.Value))
	}
	fmt.Println(batch.Offset())
	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}

	
	// messageBuffer := make([]kafka.Message, 10)
	// for {

	// 	for ind := range messageBuffer {
	// 		message, err := reader.FetchMessage(context.Background())
	// 		if err != nil {
	// 			log.Fatalf("Error while reading kafka message %s", err.Error())
	// 		}

	// 		fmt.Printf("Recieved message Topic -%s, offset - %d , partition -%d, value - %v\n",
	// 			message.Topic, message.Offset, message.Partition, message.Value)
	// 		// time.Sleep(2 * time.Second)
	// 		messageBuffer[ind] = message
	// 	}
	// 	err := reader.CommitMessages(r.Context(), messageBuffer...)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(messageBuffer)
	// 	fmt.Println("Commit")
	// }
}

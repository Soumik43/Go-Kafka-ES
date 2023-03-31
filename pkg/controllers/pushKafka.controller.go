package controller

import (
	"net/http"

	"github.com/segmentio/kafka-go"
)

func PushKafkaHandler(w http.ResponseWriter, r *http.Request) {

	data := make([]byte, r.ContentLength)
	r.Body.Read(data)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "myESTopic",
	})
	defer writer.Close()

	err := writer.WriteMessages(r.Context(), kafka.Message{
		Value: data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message pushed to kafka server!"))
}

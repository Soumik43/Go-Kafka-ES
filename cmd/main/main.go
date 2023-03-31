package main

import (
	route "go-kafka-es/pkg/routes"
	"log"
	"net/http"
)

func main() {
	route.PushRoute()
	route.WriteRoute()
	route.GetRoute()
	route.BatchMessageRoute()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

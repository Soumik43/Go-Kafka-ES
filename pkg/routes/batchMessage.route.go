package routes

import (
	controller "go-kafka-es/pkg/controllers"
	"net/http"
)

func BatchMessageRoute() {
	http.HandleFunc("/getBatch", controller.BatchMessageHandler)
}

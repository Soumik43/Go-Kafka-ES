package routes

import (
	controller "go-kafka-es/pkg/controllers"
	"net/http"
)

func WriteRoute() {
	http.HandleFunc("/write", controller.WriteHandler)
}

package routes

import (
	controller "go-kafka-es/pkg/controllers"
	"net/http"
)

func GetRoute() {
	http.HandleFunc("/get", controller.GetHandler)
}

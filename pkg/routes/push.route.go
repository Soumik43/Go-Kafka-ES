package routes

import (
	controller "go-kafka-es/pkg/controllers"
	"net/http"
)

func PushRoute() {
	http.HandleFunc("/push", controller.PushHandler)
}

package utils

import (
	"encoding/json"
	"fmt"
	controller "go-kafka-es/pkg/controllers"
)

func elasticJsonStruct(doc controller.ElasticMessage) string {

	jsonStructType := &controller.ElasticMessage{
		Message: doc.Message,
	}

	fmt.Println("\ndocStruct:", jsonStructType)

	b, err := json.Marshal(jsonStructType)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(b)
}

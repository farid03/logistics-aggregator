package main

import (
	"log"
	storage "logistics-aggregator/src/go"
	"logistics-aggregator/src/go/handler"
	"net/http"
)

func init() {
	storage.ConnectDB()
}

func main() {
	log.Printf("Server started")
	router := handler.NewRouter()
	log.Fatal(http.ListenAndServe(":12030", router))
}

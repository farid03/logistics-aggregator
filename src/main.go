package main

import (
	"log"
	sw "logistics-aggregator/src/go"
	"net/http"
)

func init() {
	sw.ConnectDB()
}

func main() {
	log.Printf("Server started")

	router := sw.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

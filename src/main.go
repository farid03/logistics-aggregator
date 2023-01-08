package main

import (
	"log"
	sw "logistics-aggregator/src/go"
	"net/http"
)

func main() {
	log.Printf("Server started")

	router := sw.NewRouter()
	sw.ConnectDB()

	log.Fatal(http.ListenAndServe(":8080", router))

}

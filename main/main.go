package main

import (
	"net/http"
	"log"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}

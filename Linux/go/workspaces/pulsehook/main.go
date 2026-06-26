package main

import (
	"log"
	"net/http"
)

func main() {
	dispatcher := NewDispatcher(4, 1024)
	dispatcher.Start()
	server := NewServer(dispatcher)
	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", server.Routes()); err != nil {
		log.Fatal(err)
	}
}

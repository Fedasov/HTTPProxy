package main

import (
	"Proxy/pkg/repository/mongodb"
	handler "Proxy/proxy/http"
	"log"
	"net/http"
)

func main() {
	collection := handler.InitDB()

	proxy := &handler.DataBase{
		Repo: mongodb.NewRequestRepository(collection),
	}

	err := http.ListenAndServe(":8080", proxy)
	if err != nil {
		log.Fatalf("Error starting proxy: %v", err)
	}
}

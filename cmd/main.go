package main

import (
	"context"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/jackc/pgx/stdlib"

	hand "Proxy/internal/pkg/http"
	repo "Proxy/internal/repository/mongodb"

	_ "Proxy/docs"
)

// @title API Proxy
// @version 1.0
// @description API server for Proxy

// @host localhost:8000
// @BasePath /
func main() {
	client, collection := InitDB()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Disabling MongoDB: %v", err)
		}
	}()

	router := setupRouter(collection)
	Server(router)
}

func InitDB() (*mongo.Client, *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Client().ApplyURI("mongodb://mongo-container:27017")

	client, err := mongo.Connect(ctx, options)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB does not respond: %v", err)
	}

	collection := client.Database("web").Collection("requests")

	return client, collection
}

func InitHandler(collection *mongo.Collection) *hand.Handler {
	requestRepo := repo.NewRequestRepository(collection)

	return hand.NewHandler(requestRepo)
}

func setupRouter(collection *mongo.Collection) http.Handler {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := setupLogRouter(collection)
	router.PathPrefix("/").Handler(api)

	return router
}

func setupLogRouter(collection *mongo.Collection) http.Handler {
	router := mux.NewRouter() //.PathPrefix("/pkg/v1").Subrouter()

	handler := InitHandler(collection)

	router.HandleFunc("/requests/{id}", handler.GetRequestByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/repeat/{id}", handler.RepeatRequest).Methods("POST", "OPTIONS")
	router.HandleFunc("/requests", handler.GetAllRequests).Methods("GET", "OPTIONS")

	return router
}

func Server(router http.Handler) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000", "http://localhost:8080"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Csrf-Token", "Content-Type", "AuthToken"},
		ExposedHeaders:   []string{"X-Csrf-Token", "AuthToken"},
	})

	corsHandler := c.Handler(router)

	fmt.Printf("The server is running on http://localhost:%d\n", 8000)
	fmt.Printf("Swagger is running on http://localhost:%d/swagger/index.html\n", 8000)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 8000), corsHandler)
	if err != nil {
		log.Fatalf("Error when starting the server: %v", err)
		return
	}
}

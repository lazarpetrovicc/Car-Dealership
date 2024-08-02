package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/lazarpetrovicc/Car-Dealership/handlers"
	"github.com/lazarpetrovicc/Car-Dealership/routers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupResponse sets up CORS headers for all responses.
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	// If the request method is OPTIONS, return early without further processing
	if req.Method == "OPTIONS" {
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
		(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
		return
	}
	// Set CORS headers
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found") // Log a message if .env file is not found
	}

	// Get MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set") // Exit if MongoDB URI is not set
	}

	// Set up a context with a timeout for MongoDB operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err) // Exit if connection fails
	}

	// Ping MongoDB to ensure the connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err) // Exit if ping fails
	}

	// Initialize the car handler with the MongoDB client and database name
	handlers.InitCarHandler(client, "carDealershipDB")

	// Initialize the router with the routes
	router := routers.InitRoutes()

	// Set up the HTTP server with CORS headers and the router for routing requests
	server := &http.Server{
		Addr: ":8000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Set CORS headers for the response
			setupResponse(&w, req)

			// Handle OPTIONS requests by returning early
			if req.Method == "OPTIONS" {
				return
			}

			// Use the router to handle the request
			router.ServeHTTP(w, req)
		}),
	}

	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Start the HTTP server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8000: %v\n", err)
		}
	}()
	log.Println("Server is ready to handle requests at :8000")

	// Block until a signal is received
	<-stop

	// Create a context with a timeout for server shutdown
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")

	// Disconnect the MongoDB client
	if err := client.Disconnect(ctxShutDown); err != nil {
		log.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}

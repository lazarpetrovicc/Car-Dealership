package routers

import (
	"github.com/gorilla/mux"
	"github.com/lazarpetrovicc/Car-Dealership/handlers"
)

// InitRoutes initializes the routes for car-related operations.
func InitRoutes() *mux.Router {
	carRouter := mux.NewRouter()

	// CRUD operations on cars

	// GET /cars/{status}
	// Fetch cars by their status (e.g., available, reserved, sold).
	carRouter.HandleFunc("/cars/{status}", handlers.GetCarsByStatus).Methods("GET")

	// POST /cars
	// Create a new car.
	carRouter.HandleFunc("/cars", handlers.CreateCar).Methods("POST")

	// PUT /cars/{id}
	// Update an existing car by its ID.
	carRouter.HandleFunc("/cars/{id}", handlers.UpdateCar).Methods("PUT")

	// DELETE /cars/{id}
	// Delete a car by its ID.
	carRouter.HandleFunc("/cars/{id}", handlers.DeleteCar).Methods("DELETE")

	// Actions on cars

	// POST /cars/{id}/reserve
	// Reserve a car by its ID.
	carRouter.HandleFunc("/cars/{id}/reserve", handlers.ReserveCar).Methods("POST")

	// POST /cars/{id}/sell
	// Sell a car to a customer by its ID.
	carRouter.HandleFunc("/cars/{id}/sell", handlers.SellCar).Methods("POST")

	// POST /cars/{id}/cancel-reservation
	// Cancel a reservation of a car by its ID.
	carRouter.HandleFunc("/cars/{id}/cancel-reservation", handlers.CancelReservation).Methods("POST")

	// Endpoint to fetch car image

	// GET /cars/image/{id}
	// Fetch the image of a car by its ID.
	carRouter.HandleFunc("/cars/image/{id}", handlers.GetCarImage).Methods("GET")

	return carRouter
}

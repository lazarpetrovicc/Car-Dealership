package services

import (
	"github.com/lazarpetrovicc/Car-Dealership/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// NewCarServiceInterface initializes and returns a new instance of the carService that satisfies the IcarService interface.
func NewCarServiceInterface(client *mongo.Client, dbName string) IcarService {
	return NewCarService(client, dbName)
}

// IcarService defines the interface for car-related operations.
type IcarService interface {
	// GetCarsByStatus retrieves cars from the database based on their status.
	// Returns a slice of cars and any error encountered.
	GetCarsByStatus(status string) ([]models.Car, error)

	// GetCarImage retrieves the image data associated with a car by its picture ID.
	// Returns the image data as a byte slice and any error encountered.
	GetCarImage(pictureID string) ([]byte, error)

	// CreateCar adds a new available car to the database and uploads its image to GridFS.
	// Returns the result of the insertion operation and any error encountered.
	CreateCar(car *models.Car, fileData []byte, fileName string) (interface{}, error)

	// UpdateCar modifies an existing car's details and updates its image in GridFS. Only available cars can be updated, and their status cannot be changed through updating.
	// Returns the result of the update operation and any error encountered.
	UpdateCar(id primitive.ObjectID, car *models.Car, fileData []byte, fileName string) (interface{}, error)

	// DeleteCar removes a car from the database and deletes its associated image from GridFS. Only available cars can be deleted.
	// Returns the result of the deletion operation and any error encountered.
	DeleteCar(id primitive.ObjectID) (interface{}, error)

	// ReserveCar changes the status of a car to "reserved" and associates a customer with it. Only available cars can be reserved.
	// Returns the result of the update operation and any error encountered.
	ReserveCar(id primitive.ObjectID, customer models.Customer) (interface{}, error)

	// CancelReservation updates the status of a reserved car back to "available" and clears customer information.
	// Returns the result of the update operation and any error encountered.
	CancelReservation(id primitive.ObjectID) (interface{}, error)

	// SellCar updates the status of a car to "sold" and associates a customer with it. Only available cars can be sold.
	// Returns the result of the update operation and any error encountered.
	SellCar(id primitive.ObjectID, customer models.Customer) (interface{}, error)

	// SetGridFSBucket sets the GridFS bucket used for storing car images.
	SetGridFSBucket(bucket *gridfs.Bucket)
}

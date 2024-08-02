package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Constants for car statuses
const (
	CarStatusAvailable = "available"
	CarStatusReserved  = "reserved"
	CarStatusSold      = "sold"
)

// Car represents a car in the dealership.
type Car struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`                                      // Unique identifier for the car
	Make     string             `bson:"make" json:"make" validate:"required"`                                   // Manufacturer of the car
	Model    string             `bson:"model" json:"model" validate:"required"`                                 // Model of the car
	Year     int                `bson:"year" json:"year" validate:"required,min=1900"`                          // Year of manufacture
	Price    float64            `bson:"price" json:"price" validate:"required,min=1"`                           // Price of the car
	Status   string             `bson:"status" json:"status" validate:"required,oneof=available reserved sold"` // Current status of the car
	Customer *Customer          `bson:"customer,omitempty" json:"customer,omitempty"`                           // Customer associated with the car (if any)
	Picture  string             `bson:"picture" json:"picture" validate:"required"`                             // GridFS file ID for the car's image
}

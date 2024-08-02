package models

// Customer represents a customer in the dealership system.
type Customer struct {
	FullName    string `bson:"fullName" json:"fullName" validate:"required"`              // Full name of the customer
	Email       string `bson:"email" json:"email" validate:"required,email"`              // Email address of the customer
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber" validate:"required,number"` // Phone number of the customer
}

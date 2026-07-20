package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/lazarpetrovicc/Car-Dealership/models"
	"github.com/lazarpetrovicc/Car-Dealership/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	validate   *validator.Validate
	carService services.IcarService
)

// SetCarService sets the carService variable for testing purposes
func SetCarService(service services.IcarService) {
	carService = service
}

// SetValidator sets the global validator instance for testing purposes
func SetValidator(v *validator.Validate) {
	validate = v
}

// InitCarHandler initializes the car handler with the given MongoDB client and database name
func InitCarHandler(client *mongo.Client, dbName string) {
	validate = validator.New()
	carService = services.NewCarServiceInterface(client, dbName)
}

// HealthCheck returns a simple readiness response.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetCarsByStatus retrieves cars by their status and returns them in JSON format
func GetCarsByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	// Check if the status is one of the valid constants
	switch status {
	case models.CarStatusAvailable, models.CarStatusReserved, models.CarStatusSold:
		cars, err := carService.GetCarsByStatus(status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, cars)
	default:
		http.Error(w, "Invalid status provided", http.StatusBadRequest)
	}
}

// GetCarImage retrieves a car's image by its ID and returns it in JPEG format
func GetCarImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pictureID := vars["id"]

	// Retrieve the car image data from the service
	fileData, err := carService.GetCarImage(pictureID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set headers and write the image data to the response
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(fileData)))
	w.Write(fileData)
}

// CreateCar handles the creation of a new available car and saves its details in the database
func CreateCar(w http.ResponseWriter, r *http.Request) {
	var car models.Car

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	// Extract form values
	car.Make = r.FormValue("make")
	car.Model = r.FormValue("model")
	year, _ := strconv.Atoi(r.FormValue("year"))
	car.Year = year
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	car.Price = price
	car.Status = models.CarStatusAvailable // Ensure the status remains available

	// Retrieve the file from the form
	file, handler, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, "picture is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file into a byte slice
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read uploaded picture", http.StatusBadRequest)
		return
	}

	car.Picture = string(fileData)

	// Validate the car struct
	if err := validate.Struct(car); err != nil {
		log.Println("Validation errors: ", err)
		handleValidationErrors(w, err)
		return
	}

	// Save the car in the database
	result, err := carService.CreateCar(&car, fileData, handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, result)
}

// UpdateCar handles updating the details of an existing available car in the database. Only available cars can be updated, and their status cannot be changed through updating.
func UpdateCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	carID := vars["id"]
	id, err := primitive.ObjectIDFromHex(carID)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	var car models.Car

	// Parse the multipart form data
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	// Extract form values
	car.Make = r.FormValue("make")
	car.Model = r.FormValue("model")
	year, _ := strconv.Atoi(r.FormValue("year"))
	car.Year = year
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	car.Price = price
	car.Status = models.CarStatusAvailable // Ensure the status remains available

	// Retrieve the file from the form, if present
	var fileData []byte
	var fileName string
	file, handler, err := r.FormFile("picture")
	if err == nil {
		defer file.Close()
		fileData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		fileName = handler.Filename
		car.Picture = string(fileData)
	}

	// Validate the car struct
	if err := validate.Struct(car); err != nil {
		log.Println("Validation errors: ", err)
		handleValidationErrors(w, err)
		return
	}

	// Update the car in the database
	result, err := carService.UpdateCar(id, &car, fileData, fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, result)
}

// DeleteCar handles deleting a car from the database by its ID. Only available cars can be deleted.
func DeleteCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	// Delete the car from the database
	result, err := carService.DeleteCar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, result)
}

// ReserveCar handles reserving a car by a customer. Only available cars can be reserved.
func ReserveCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid customer data", http.StatusBadRequest)
		return
	}

	// Validate the customer struct
	if err := validate.Struct(customer); err != nil {
		log.Println("Validation errors: ", err)
		handleValidationErrors(w, err)
		return
	}

	// Reserve the car for the customer
	result, err := carService.ReserveCar(id, customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, result)
}

// CancelReservation handles canceling a car reservation
func CancelReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	// Cancel the car reservation
	result, err := carService.CancelReservation(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, result)
}

// SellCar handles selling a car to a customer. Only available cars can be sold.
func SellCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid customer data", http.StatusBadRequest)
		return
	}

	// Validate the customer struct
	if err := validate.Struct(customer); err != nil {
		log.Println("Validation errors: ", err)
		handleValidationErrors(w, err)
		return
	}

	// Sell the car to the customer
	result, err := carService.SellCar(id, customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, result)
}

// handleValidationErrors formats and returns validation errors in JSON format
func handleValidationErrors(w http.ResponseWriter, err error) {
	validationErrors := make(map[string]string)

	// Check if the error is a validation error
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, err := range errs {
			field := err.Field()
			switch err.Tag() {
			case "required":
				validationErrors[field] = field + " is required"
			case "email":
				validationErrors[field] = field + " is not a valid email address"
			case "min":
				validationErrors[field] = field + " must be at least " + err.Param()
			case "gt":
				validationErrors[field] = field + " must be greater than " + err.Param()
			case "oneof":
				validationErrors[field] = field + " must be one of " + err.Param()
			default:
				validationErrors[field] = field + " validation failed"
			}
		}
	}

	writeJSONResponse(w, http.StatusBadRequest, validationErrors)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

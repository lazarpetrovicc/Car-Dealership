package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/lazarpetrovicc/Car-Dealership/models"
	"github.com/lazarpetrovicc/Car-Dealership/services"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Constant for the testing database name
const testDbName = "carDealershipDB_test"

// setupTestDB initializes the test database, connects to MongoDB, and returns the client and database instances.
func setupTestDB(t *testing.T) (*mongo.Client, *mongo.Database) {
	// Load environment variables from .env file located in the parent directory
	err := godotenv.Load("../.env")
	if err != nil {
		t.Log("No .env file found") // Log a message if .env file is not found
	}

	// Set up a context with a timeout for MongoDB operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get MongoDB URI from environment variable
	mongoURI := "mongodb://localhost:27017" // Default to local MongoDB if not set
	if uri := os.Getenv("MONGO_TEST_URI"); uri != "" {
		mongoURI = uri
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping MongoDB to ensure the connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database(testDbName)

	// Clear existing data
	clearCollection(t, db)

	return client, db
}

// clearCollection clears the specified collections in the database.
func clearCollection(t *testing.T, db *mongo.Database) {
	// Clear the "cars" collection
	err := db.Collection("cars").Drop(context.Background())
	if err != nil && err != mongo.ErrNoDocuments {
		t.Fatalf("Failed to clear cars collection: %v", err)
	}

	// Clear the "fs.files" collection
	err = db.Collection("fs.files").Drop(context.Background())
	if err != nil && err != mongo.ErrNoDocuments {
		t.Fatalf("Failed to clear fs.files collection: %v", err)
	}

	// Clear the "fs.chunks" collection
	err = db.Collection("fs.chunks").Drop(context.Background())
	if err != nil && err != mongo.ErrNoDocuments {
		t.Fatalf("Failed to clear fs.chunks collection: %v", err)
	}
}

// TestGetCarsByStatusService tests the retrieval of cars by their status.
func TestGetCarsByStatusService(t *testing.T) {
	client, db := setupTestDB(t)
	defer func() {
		clearCollection(t, db)
		client.Disconnect(context.Background())
	}()

	service := services.NewCarServiceInterface(client, testDbName)
	var serviceInterface services.IcarService = service

	// Insert test data
	cars := []interface{}{
		models.Car{ID: primitive.NewObjectID(), Status: "available"},
		models.Car{ID: primitive.NewObjectID(), Status: "sold"},
		models.Car{ID: primitive.NewObjectID(), Status: "available"},
	}
	_, err := db.Collection("cars").InsertMany(context.Background(), cars)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Test GetCarsByStatus
	status := "available"
	result, err := serviceInterface.GetCarsByStatus(status)
	if err != nil {
		t.Fatalf("GetCarsByStatus failed: %v", err)
	}

	// Verify the result
	assert.Equal(t, 2, len(result), "Expected 2 cars with status 'available'")
	for _, car := range result {
		assert.Equal(t, status, car.Status, "Car status does not match")
	}
}

// TestGetCarImageService tests retrieving a car image from GridFS.
func TestGetCarImageService(t *testing.T) {
	client, db := setupTestDB(t)
	defer func() {
		clearCollection(t, db)
		client.Disconnect(context.Background())
	}()

	// Create a GridFS bucket
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		t.Fatalf("Failed to create GridFS bucket: %v", err)
	}

	service := services.NewCarServiceInterface(client, testDbName)
	service.SetGridFSBucket(bucket)

	// Prepare test image data
	fileData := []byte("test image data for car image")
	fileName := "testCarImage.jpg"

	// Upload the test image to GridFS
	uploadStream, err := bucket.OpenUploadStream(fileName)
	if err != nil {
		t.Fatalf("Failed to open upload stream: %v", err)
	}

	fileID := uploadStream.FileID.(primitive.ObjectID)
	_, err = uploadStream.Write(fileData)
	if err != nil {
		t.Fatalf("Failed to write data to upload stream: %v", err)
	}

	err = uploadStream.Close()
	if err != nil {
		t.Fatalf("Failed to close upload stream: %v", err)
	}

	// Test GetCarImage
	result, err := service.GetCarImage(fileID.Hex())
	if err != nil {
		t.Fatalf("GetCarImage failed: %v", err)
	}

	// Verify the result
	assert.Equal(t, fileData, result, "Image data does not match")
}

// TestCreateCarService tests creating a new car entry.
func TestCreateCarService(t *testing.T) {
	client, db := setupTestDB(t)
	defer func() {
		clearCollection(t, db)
		client.Disconnect(context.Background())
	}()

	service := services.NewCarServiceInterface(client, testDbName)
	var serviceInterface services.IcarService = service

	// Prepare test data
	fileData := []byte("test image data")
	fileName := "testImage.jpg"
	car := &models.Car{
		Make:    "Toyota",
		Model:   "Corolla",
		Year:    2022,
		Price:   20000,
		Status:  models.CarStatusAvailable,
		Picture: fileName,
	}

	// Test CreateCar
	result, err := serviceInterface.CreateCar(car, fileData, fileName)
	if err != nil {
		t.Fatalf("CreateCar failed: %v", err)
	}

	// Cast result to *mongo.InsertOneResult
	insertResult, ok := result.(*mongo.InsertOneResult)
	if !ok {
		t.Fatalf("Expected *mongo.InsertOneResult, got %T", result)
	}

	// Verify insertion
	insertedID := insertResult.InsertedID.(primitive.ObjectID)
	var insertedCar models.Car
	err = db.Collection("cars").FindOne(context.Background(), bson.M{"_id": insertedID}).Decode(&insertedCar)
	if err != nil {
		t.Fatalf("Failed to find inserted car: %v", err)
	}

	assert.Equal(t, car.Make, insertedCar.Make, "Car Make does not match")
	assert.Equal(t, car.Model, insertedCar.Model, "Car Model does not match")
	assert.Equal(t, car.Year, insertedCar.Year, "Car Year does not match")
	assert.Equal(t, car.Price, insertedCar.Price, "Car Price does not match")
	assert.Equal(t, car.Status, insertedCar.Status, "Car Status does not match")
}

// TestUpdateCarService tests updating an existing car entry.
func TestUpdateCarService(t *testing.T) {
	client, db := setupTestDB(t)
	defer func() {
		clearCollection(t, db)
		client.Disconnect(context.Background())
	}()

	service := services.NewCarServiceInterface(client, testDbName)
	var serviceInterface services.IcarService = service

	// Create a car to update
	fileData := []byte("test updated image data")
	fileName := "updatedImage.jpg"
	car := &models.Car{
		Make:    "Toyota",
		Model:   "Corolla",
		Year:    2022,
		Price:   20000,
		Status:  models.CarStatusAvailable,
		Picture: fileName,
	}

	result, err := serviceInterface.CreateCar(car, fileData, fileName)
	if err != nil {
		t.Fatalf("CreateCar failed: %v", err)
	}

	// Cast result to *mongo.InsertOneResult
	insertResult, ok := result.(*mongo.InsertOneResult)
	if !ok {
		t.Fatalf("Expected *mongo.InsertOneResult, got %T", result)
	}

	carID := insertResult.InsertedID.(primitive.ObjectID)

	// Prepare updated car data
	updatedCar := &models.Car{
		Make:    "Toyota",
		Model:   "Corolla Updated",
		Year:    2023,
		Price:   21000,
		Status:  models.CarStatusAvailable,
		Picture: fileName,
	}

	// Test UpdateCar
	updateResult, err := serviceInterface.UpdateCar(carID, updatedCar, fileData, fileName)
	if err != nil {
		t.Fatalf("UpdateCar failed: %v", err)
	}

	// Cast result to *mongo.UpdateResult
	updateResultCast, ok := updateResult.(*mongo.UpdateResult)
	if !ok {
		t.Fatalf("Expected *mongo.UpdateResult, got %T", updateResult)
	}

	if updateResultCast.ModifiedCount == 0 {
		t.Fatalf("Expected 1 car to be updated, got 0")
	}

	// Verify update
	var updatedCarResult models.Car
	err = db.Collection("cars").FindOne(context.Background(), bson.M{"_id": carID}).Decode(&updatedCarResult)
	if err != nil {
		t.Fatalf("Failed to find updated car: %v", err)
	}

	assert.Equal(t, updatedCar.Make, updatedCarResult.Make, "Car Make does not match")
	assert.Equal(t, updatedCar.Model, updatedCarResult.Model, "Car Model does not match")
	assert.Equal(t, updatedCar.Year, updatedCarResult.Year, "Car Year does not match")
	assert.Equal(t, updatedCar.Price, updatedCarResult.Price, "Car Price does not match")
	assert.Equal(t, updatedCar.Status, updatedCarResult.Status, "Car Status does not match")
}

// TestDeleteCarService tests deleting a car entry.
func TestDeleteCarService(t *testing.T) {
	client, db := setupTestDB(t)
	defer func() {
		clearCollection(t, db)
		client.Disconnect(context.Background())
	}()

	service := services.NewCarServiceInterface(client, testDbName)
	var serviceInterface services.IcarService = service

	// Create a car to delete
	fileData := []byte("test image data")
	fileName := "testImage.jpg"
	car := &models.Car{
		Make:    "Toyota",
		Model:   "Corolla",
		Year:    2022,
		Price:   20000,
		Status:  models.CarStatusAvailable,
		Picture: fileName,
	}

	result, err := serviceInterface.CreateCar(car, fileData, fileName)
	if err != nil {
		t.Fatalf("CreateCar failed: %v", err)
	}

	// Cast result to *mongo.InsertOneResult
	insertResult, ok := result.(*mongo.InsertOneResult)
	if !ok {
		t.Fatalf("Expected *mongo.InsertOneResult, got %T", result)
	}

	carID := insertResult.InsertedID.(primitive.ObjectID)

	// Test DeleteCar
	deleteResult, err := serviceInterface.DeleteCar(carID)
	if err != nil {
		t.Fatalf("DeleteCar failed: %v", err)
	}

	// Cast result to *mongo.DeleteResult
	deleteResultCast, ok := deleteResult.(*mongo.DeleteResult)
	if !ok {
		t.Fatalf("Expected *mongo.DeleteResult, got %T", deleteResult)
	}

	if deleteResultCast.DeletedCount == 0 {
		t.Fatalf("Expected 1 car to be deleted, got 0")
	}

	// Verify deletion
	err = db.Collection("cars").FindOne(context.Background(), bson.M{"_id": carID}).Decode(&models.Car{})
	if err != mongo.ErrNoDocuments {
		t.Fatalf("Expected no documents, but found one: %v", err)
	}
}

// TestReserveCarService tests reserving a car.
func TestReserveCarService(t *testing.T) {
	client, db := setupTestDB(t)
	defer func() {
		clearCollection(t, db)
		client.Disconnect(context.Background())
	}()

	service := services.NewCarServiceInterface(client, testDbName)
	var serviceInterface services.IcarService = service

	// Create a car to reserve
	fileData := []byte("test image data")
	fileName := "testImage.jpg"
	car := &models.Car{
		Make:    "Toyota",
		Model:   "Corolla",
		Year:    2022,
		Price:   20000,
		Status:  models.CarStatusAvailable,
		Picture: fileName,
	}

	result, err := serviceInterface.CreateCar(car, fileData, fileName)
	if err != nil {
		t.Fatalf("CreateCar failed: %v", err)
	}

	// Cast result to *mongo.InsertOneResult
	insertResult, ok := result.(*mongo.InsertOneResult)
	if !ok {
		t.Fatalf("Expected *mongo.InsertOneResult, got %T", result)
	}

	carID := insertResult.InsertedID.(primitive.ObjectID)

	// Prepare customer data
	customer := models.Customer{
		FullName:    "John Doe",
		Email:       "john.doe@example.com",
		PhoneNumber: "1234567890",
	}

	// Test ReserveCar
	reserveResult, err := serviceInterface.ReserveCar(carID, customer)
	if err != nil {
		t.Fatalf("ReserveCar failed: %v", err)
	}

	// Cast result to *mongo.UpdateResult
	reserveResultCast, ok := reserveResult.(*mongo.UpdateResult)
	if !ok {
		t.Fatalf("Expected *mongo.UpdateResult, got %T", reserveResult)
	}

	if reserveResultCast.ModifiedCount == 0 {
		t.Fatalf("Expected 1 car to be reserved, got 0")
	}

	// Verify reservation
	var reservedCar models.Car
	err = db.Collection("cars").FindOne(context.Background(), bson.M{"_id": carID}).Decode(&reservedCar)
	if err != nil {
		t.Fatalf("Failed to find reserved car: %v", err)
	}

	assert.Equal(t, models.CarStatusReserved, reservedCar.Status, "Car Status does not match")
	assert.Equal(t, customer.FullName, reservedCar.Customer.FullName, "Customer FullName does not match")
	assert.Equal(t, customer.Email, reservedCar.Customer.Email, "Customer Email does not match")
	assert.Equal(t, customer.PhoneNumber, reservedCar.Customer.PhoneNumber, "Customer PhoneNumber does not match")
}

// TestReserveAndCancelCarService tests reserving and then canceling a reservation for a car.
func TestReserveAndCancelCarService(t *testing.T) {
	client, db := setupTestDB(t)
	defer func() {
		clearCollection(t, db)
		client.Disconnect(context.Background())
	}()

	service := services.NewCarServiceInterface(client, testDbName)
	var serviceInterface services.IcarService = service

	// Create a car to reserve
	fileData := []byte("test image data")
	fileName := "testImage.jpg"
	car := &models.Car{
		Make:    "Toyota",
		Model:   "Corolla",
		Year:    2022,
		Price:   20000,
		Status:  models.CarStatusAvailable,
		Picture: fileName,
	}

	result, err := serviceInterface.CreateCar(car, fileData, fileName)
	if err != nil {
		t.Fatalf("CreateCar failed: %v", err)
	}

	// Cast result to *mongo.InsertOneResult
	insertResult, ok := result.(*mongo.InsertOneResult)
	if !ok {
		t.Fatalf("Expected *mongo.InsertOneResult, got %T", result)
	}

	carID := insertResult.InsertedID.(primitive.ObjectID)

	// Prepare customer data
	customer := models.Customer{
		FullName:    "John Doe",
		Email:       "john.doe@example.com",
		PhoneNumber: "1234567890",
	}

	// Test ReserveCar
	reserveResult, err := serviceInterface.ReserveCar(carID, customer)
	if err != nil {
		t.Fatalf("ReserveCar failed: %v", err)
	}

	// Cast result to *mongo.UpdateResult
	reserveResultCast, ok := reserveResult.(*mongo.UpdateResult)
	if !ok {
		t.Fatalf("Expected *mongo.UpdateResult, got %T", reserveResult)
	}

	if reserveResultCast.ModifiedCount == 0 {
		t.Fatalf("Expected 1 car to be reserved, got 0")
	}

	// Verify reservation
	var reservedCar models.Car
	err = db.Collection("cars").FindOne(context.Background(), bson.M{"_id": carID}).Decode(&reservedCar)
	if err != nil {
		t.Fatalf("Failed to find reserved car: %v", err)
	}

	assert.Equal(t, models.CarStatusReserved, reservedCar.Status, "Car Status does not match")
	assert.Equal(t, customer.FullName, reservedCar.Customer.FullName, "Customer FullName does not match")
	assert.Equal(t, customer.Email, reservedCar.Customer.Email, "Customer Email does not match")
	assert.Equal(t, customer.PhoneNumber, reservedCar.Customer.PhoneNumber, "Customer PhoneNumber does not match")

	// Test CancelReservation
	cancelResult, err := serviceInterface.CancelReservation(carID)
	if err != nil {
		t.Fatalf("CancelReservation failed: %v", err)
	}

	// Cast result to *mongo.UpdateResult
	cancelResultCast, ok := cancelResult.(*mongo.UpdateResult)
	if !ok {
		t.Fatalf("Expected *mongo.UpdateResult, got %T", cancelResult)
	}

	if cancelResultCast.ModifiedCount == 0 {
		t.Fatalf("Expected 1 reservation to be canceled, got 0")
	}

	// Verify cancellation
	var canceledCar models.Car
	err = db.Collection("cars").FindOne(context.Background(), bson.M{"_id": carID}).Decode(&canceledCar)
	if err != nil {
		t.Fatalf("Failed to find canceled car: %v", err)
	}

	assert.Equal(t, models.CarStatusAvailable, canceledCar.Status, "Car Status does not match")
	assert.Nil(t, canceledCar.Customer, "Customer information should be cleared")
}

// TestSellCarService tests selling a car.
func TestSellCarService(t *testing.T) {
	client, db := setupTestDB(t)
	defer func() {
		clearCollection(t, db)
		client.Disconnect(context.Background())
	}()

	service := services.NewCarServiceInterface(client, testDbName)
	var serviceInterface services.IcarService = service

	// Create a car to sell
	fileData := []byte("test image data")
	fileName := "testImage.jpg"
	car := &models.Car{
		Make:    "Honda",
		Model:   "Civic",
		Year:    2023,
		Price:   22000,
		Status:  models.CarStatusAvailable,
		Picture: fileName,
	}

	result, err := serviceInterface.CreateCar(car, fileData, fileName)
	if err != nil {
		t.Fatalf("CreateCar failed: %v", err)
	}

	// Cast result to *mongo.InsertOneResult
	insertResult, ok := result.(*mongo.InsertOneResult)
	if !ok {
		t.Fatalf("Expected *mongo.InsertOneResult, got %T", result)
	}

	carID := insertResult.InsertedID.(primitive.ObjectID)

	// Prepare customer data
	customer := &models.Customer{
		FullName:    "John Doe",
		Email:       "john.doe@example.com",
		PhoneNumber: "1234567890",
	}

	// Test SellCar
	sellResult, err := serviceInterface.SellCar(carID, *customer)
	if err != nil {
		t.Fatalf("SellCar failed: %v", err)
	}

	// Cast result to *mongo.UpdateResult
	sellResultCast, ok := sellResult.(*mongo.UpdateResult)
	if !ok {
		t.Fatalf("Expected *mongo.UpdateResult, got %T", sellResult)
	}

	if sellResultCast.ModifiedCount == 0 {
		t.Fatalf("Expected 1 car to be sold, got 0")
	}

	// Verify sale
	var soldCar models.Car
	err = db.Collection("cars").FindOne(context.Background(), bson.M{"_id": carID}).Decode(&soldCar)
	if err != nil {
		t.Fatalf("Failed to find sold car: %v", err)
	}

	assert.Equal(t, models.CarStatusSold, soldCar.Status, "Car Status does not match")
	assert.Equal(t, customer.FullName, soldCar.Customer.FullName, "Customer FullName does not match")
	assert.Equal(t, customer.Email, soldCar.Customer.Email, "Customer Email does not match")
	assert.Equal(t, customer.PhoneNumber, soldCar.Customer.PhoneNumber, "Customer PhoneNumber does not match")
}

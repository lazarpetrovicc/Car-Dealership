package services

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/lazarpetrovicc/Car-Dealership/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// carService provides methods to manage cars and their associated images.
type carService struct {
	carCollection *mongo.Collection // MongoDB collection for storing cars
	gridFSBucket  *gridfs.Bucket    // GridFS bucket for storing car images
}

// NewCarService initializes a new instance of carService.
func NewCarService(client *mongo.Client, dbName string) *carService {
	db := client.Database(dbName)
	carCollection := db.Collection("cars")
	bucket, _ := gridfs.NewBucket(db)
	return &carService{
		carCollection: carCollection,
		gridFSBucket:  bucket,
	}
}

// SetGridFSBucket sets the GridFS bucket used for storing car images.
func (s *carService) SetGridFSBucket(bucket *gridfs.Bucket) {
	s.gridFSBucket = bucket
}

// GetCarsByStatus retrieves cars from the database based on their status.
// Returns a slice of cars and any error encountered.
func (s *carService) GetCarsByStatus(status string) ([]models.Car, error) {
	var cars []models.Car
	cursor, err := s.carCollection.Find(context.Background(), bson.M{"status": status})
	if err != nil {
		log.Printf("Error finding cars by status '%s': %v", status, err)
		return nil, err
	}
	if err = cursor.All(context.Background(), &cars); err != nil {
		log.Printf("Error decoding cars by status '%s': %v", status, err)
		return nil, err
	}
	return cars, nil
}

// GetCarImage retrieves the image data for a specific car based on its picture ID.
// Returns a byte slice containing the image data and any error encountered.
func (s *carService) GetCarImage(pictureID string) ([]byte, error) {
	oid, err := primitive.ObjectIDFromHex(pictureID)
	if err != nil {
		log.Printf("Error converting pictureID '%s' to ObjectID: %v", pictureID, err)
		return nil, err
	}

	dStream, err := s.gridFSBucket.OpenDownloadStream(oid)
	if err != nil {
		log.Printf("Error opening download stream for pictureID '%s': %v", pictureID, err)
		return nil, err
	}
	defer dStream.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, dStream)
	if err != nil {
		log.Printf("Error copying data from download stream for pictureID '%s': %v", pictureID, err)
		return nil, err
	}

	return buf.Bytes(), nil
}

// CreateCar inserts a new available car document into the database and uploads its image to GridFS.
// Returns the MongoDB InsertOneResult and any error encountered.
func (s *carService) CreateCar(car *models.Car, fileData []byte, fileName string) (interface{}, error) {
	// Ensure that the car status is available
	car.Status = models.CarStatusAvailable

	// Upload the image to GridFS
	uploadStream, err := s.gridFSBucket.OpenUploadStream(fileName)
	if err != nil {
		log.Printf("Error opening upload stream for file '%s': %v", fileName, err)
		return nil, err
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(fileData)
	if err != nil {
		log.Printf("Error writing file '%s' to upload stream: %v", fileName, err)
		return nil, err
	}

	car.Picture = uploadStream.FileID.(primitive.ObjectID).Hex()
	// Insert the car document into the collection
	result, err := s.carCollection.InsertOne(context.Background(), car)
	if err != nil {
		log.Printf("Error inserting car into collection: %v", err)
		return nil, err
	}
	return result, nil
}

// UpdateCar updates an existing available car document in the database and updates its image in GridFS. Only available cars can be updated, and their status cannot be changed through updating.
// Returns the MongoDB UpdateOneResult and any error encountered.
func (s *carService) UpdateCar(id primitive.ObjectID, car *models.Car, fileData []byte, fileName string) (interface{}, error) {
	if fileData != nil {
		// Find the existing available car to get the current picture ID
		var existingCar models.Car
		err := s.carCollection.FindOne(context.Background(), bson.M{"_id": id, "status": models.CarStatusAvailable}).Decode(&existingCar)
		if err != nil {
			log.Printf("Error finding existing car with ID '%s': %v", id.Hex(), err)
			return nil, err
		}

		// Delete the old photo from GridFS if it exists
		if existingCar.Picture != "" {
			oldPictureID, err := primitive.ObjectIDFromHex(existingCar.Picture)
			if err != nil {
				log.Printf("Error converting old picture ID '%s' to ObjectID: %v", existingCar.Picture, err)
				return nil, err
			}
			err = s.gridFSBucket.Delete(oldPictureID)
			if err != nil {
				log.Printf("Error deleting old picture with ID '%s': %v", oldPictureID.Hex(), err)
				return nil, err
			}
		}

		// Upload the new photo to GridFS
		uploadStream, err := s.gridFSBucket.OpenUploadStream(fileName)
		if err != nil {
			log.Printf("Error opening upload stream for file '%s': %v", fileName, err)
			return nil, err
		}
		defer uploadStream.Close()

		_, err = uploadStream.Write(fileData)
		if err != nil {
			log.Printf("Error writing file '%s' to upload stream: %v", fileName, err)
			return nil, err
		}
		car.Picture = uploadStream.FileID.(primitive.ObjectID).Hex()
	}

	// Ensure that the updated car status remains "available"
	car.Status = models.CarStatusAvailable

	update := bson.D{{Key: "$set", Value: car}}
	result, err := s.carCollection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		log.Printf("Error updating car with ID '%s': %v", id.Hex(), err)
		return nil, err
	}
	return result, nil
}

// DeleteCar removes a car document from the database and deletes its associated image from GridFS. Only available cars can be deleted.
// Returns the MongoDB DeleteOneResult and any error encountered.
func (s *carService) DeleteCar(id primitive.ObjectID) (interface{}, error) {
	var car models.Car
	err := s.carCollection.FindOne(context.Background(), bson.M{"_id": id, "status": models.CarStatusAvailable}).Decode(&car)
	if err != nil {
		log.Printf("Error finding car with ID '%s' for deletion: %v", id.Hex(), err)
		return nil, err
	}

	// Delete the associated image from GridFS if it exists
	if car.Picture != "" {
		pictureID, err := primitive.ObjectIDFromHex(car.Picture)
		if err != nil {
			log.Printf("Error converting picture ID '%s' to ObjectID: %v", car.Picture, err)
			return nil, err
		}
		err = s.gridFSBucket.Delete(pictureID)
		if err != nil {
			log.Printf("Error deleting picture with ID '%s': %v", pictureID.Hex(), err)
			return nil, err
		}
	}

	// Delete the car document from the collection
	result, err := s.carCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Printf("Error deleting car with ID '%s': %v", id.Hex(), err)
		return nil, err
	}
	return result, nil
}

// ReserveCar updates the status of a car to "reserved" and assigns a customer to it. Only available cars can be reserved.
// Returns the MongoDB UpdateOneResult and any error encountered.
func (s *carService) ReserveCar(id primitive.ObjectID, customer models.Customer) (interface{}, error) {
	result, err := s.carCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id, "status": models.CarStatusAvailable},
		bson.D{{Key: "$set", Value: bson.M{"status": models.CarStatusReserved, "customer": customer}}},
	)
	if err != nil {
		log.Printf("Error reserving car with ID '%s': %v", id.Hex(), err)
		return nil, err
	}
	return result, nil
}

// CancelReservation updates the status of a reserved car back to "available" and clears the customer information.
// Returns the MongoDB UpdateOneResult and any error encountered.
func (s *carService) CancelReservation(id primitive.ObjectID) (interface{}, error) {
	result, err := s.carCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id, "status": models.CarStatusReserved},
		bson.D{{Key: "$set", Value: bson.M{"status": models.CarStatusAvailable, "customer": nil}}},
	)
	if err != nil {
		log.Printf("Error canceling reservation for car with ID '%s': %v", id.Hex(), err)
		return nil, err
	}
	return result, nil
}

// SellCar updates the status of a car to "sold" and assigns a customer to it. Only available cars can be sold.
// Returns the MongoDB UpdateOneResult and any error encountered.
func (s *carService) SellCar(id primitive.ObjectID, customer models.Customer) (interface{}, error) {
	result, err := s.carCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id, "status": models.CarStatusAvailable},
		bson.D{{Key: "$set", Value: bson.M{"status": models.CarStatusSold, "customer": customer}}},
	)
	if err != nil {
		log.Printf("Error selling car with ID '%s': %v", id.Hex(), err)
		return nil, err
	}
	return result, nil
}

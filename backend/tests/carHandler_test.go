package tests

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/lazarpetrovicc/Car-Dealership/handlers"
	"github.com/lazarpetrovicc/Car-Dealership/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// MockCarService is a mock implementation of the IcarService interface
type MockCarService struct {
	GetCarsByStatusFunc   func(status string) ([]models.Car, error)
	GetCarImageFunc       func(pictureID string) ([]byte, error)
	CreateCarFunc         func(car *models.Car, fileData []byte, fileName string) (interface{}, error)
	UpdateCarFunc         func(id primitive.ObjectID, car *models.Car, fileData []byte, fileName string) (interface{}, error)
	DeleteCarFunc         func(id primitive.ObjectID) (interface{}, error)
	ReserveCarFunc        func(id primitive.ObjectID, customer models.Customer) (interface{}, error)
	CancelReservationFunc func(id primitive.ObjectID) (interface{}, error)
	SellCarFunc           func(id primitive.ObjectID, customer models.Customer) (interface{}, error)
	SetGridFSBucketFunc   func(bucket *gridfs.Bucket)
}

// Implementing the IcarService interface methods using function fields in MockCarService
func (m *MockCarService) GetCarsByStatus(status string) ([]models.Car, error) {
	return m.GetCarsByStatusFunc(status)
}

func (m *MockCarService) GetCarImage(pictureID string) ([]byte, error) {
	return m.GetCarImageFunc(pictureID)
}

func (m *MockCarService) CreateCar(car *models.Car, fileData []byte, fileName string) (interface{}, error) {
	return m.CreateCarFunc(car, fileData, fileName)
}

func (m *MockCarService) UpdateCar(id primitive.ObjectID, car *models.Car, fileData []byte, fileName string) (interface{}, error) {
	return m.UpdateCarFunc(id, car, fileData, fileName)
}

func (m *MockCarService) DeleteCar(id primitive.ObjectID) (interface{}, error) {
	return m.DeleteCarFunc(id)
}

func (m *MockCarService) ReserveCar(id primitive.ObjectID, customer models.Customer) (interface{}, error) {
	return m.ReserveCarFunc(id, customer)
}

func (m *MockCarService) CancelReservation(id primitive.ObjectID) (interface{}, error) {
	return m.CancelReservationFunc(id)
}

func (m *MockCarService) SellCar(id primitive.ObjectID, customer models.Customer) (interface{}, error) {
	return m.SellCarFunc(id, customer)
}

func (m *MockCarService) SetGridFSBucket(bucket *gridfs.Bucket) {
	if m.SetGridFSBucketFunc != nil {
		m.SetGridFSBucketFunc(bucket)
	}
}

// Helper function to create a new multipart form request
// method: HTTP method (e.g., "POST", "PUT")
// url: request URL
// fields: map of form fields
// fileField: name of the file form field
// fileContent: content of the file to be uploaded
func newMultipartRequest(method, url string, fields map[string]string, fileField string, fileContent []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Adding form fields
	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}

	// Adding file field if specified
	if fileField != "" {
		part, err := writer.CreateFormFile(fileField, "file.jpg")
		if err != nil {
			return nil, err
		}
		_, err = part.Write(fileContent)
		if err != nil {
			return nil, err
		}
	}

	// Closing the writer to finalize the form
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	// Creating the HTTP request with the multipart form data
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

func TestGetCarsByStatus(t *testing.T) {
	mockCarService := &MockCarService{
		GetCarsByStatusFunc: func(status string) ([]models.Car, error) {
			if status == models.CarStatusAvailable {
				return []models.Car{
					{Make: "Toyota", Model: "Corolla", Year: 2020},
					{Make: "Honda", Model: "Civic", Year: 2021},
				}, nil
			}
			return nil, nil
		},
	}

	// Set the mock service in the handler
	handlers.SetCarService(mockCarService)

	t.Run("valid status", func(t *testing.T) {
		// Creating a request with a valid status
		req := httptest.NewRequest("GET", "/cars/status/available", nil)
		req = mux.SetURLVars(req, map[string]string{"status": "available"})
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.GetCarsByStatus(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusOK, rr.Code)

		var result []models.Car
		json.NewDecoder(rr.Body).Decode(&result)
		expected := []models.Car{
			{Make: "Toyota", Model: "Corolla", Year: 2020},
			{Make: "Honda", Model: "Civic", Year: 2021},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("invalid status", func(t *testing.T) {
		// Creating a request with an invalid status
		req := httptest.NewRequest("GET", "/cars/status/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"status": "invalid"})
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.GetCarsByStatus(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid status provided\n", rr.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Simulating a service error
		mockCarService.GetCarsByStatusFunc = func(status string) ([]models.Car, error) {
			if status == models.CarStatusReserved {
				return nil, assert.AnError
			}
			return nil, nil
		}

		// Creating a request with a status that causes an error
		req := httptest.NewRequest("GET", "/cars/status/reserved", nil)
		req = mux.SetURLVars(req, map[string]string{"status": "reserved"})
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.GetCarsByStatus(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, assert.AnError.Error()+"\n", rr.Body.String())
	})
}

func TestGetCarImage(t *testing.T) {
	mockCarService := &MockCarService{
		GetCarImageFunc: func(pictureID string) ([]byte, error) {
			if pictureID == "valid-id" {
				return []byte("fake image data"), nil
			}
			return nil, assert.AnError
		},
	}

	handlers.SetCarService(mockCarService)

	t.Run("valid image ID", func(t *testing.T) {
		// Creating a request with a valid image ID
		req, err := http.NewRequest("GET", "/cars/images/valid-id", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "valid-id"})

		rr := httptest.NewRecorder()
		handlers.GetCarImage(rr, req)

		// Checking the response status and headers
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "image/jpeg", rr.Header().Get("Content-Type"))
		assert.Equal(t, strconv.Itoa(len("fake image data")), rr.Header().Get("Content-Length"))
		assert.Equal(t, "fake image data", rr.Body.String())
	})

	t.Run("invalid image ID", func(t *testing.T) {
		// Creating a request with an invalid image ID
		req, err := http.NewRequest("GET", "/cars/images/invalid-id", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-id"})

		rr := httptest.NewRecorder()
		handlers.GetCarImage(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, assert.AnError.Error()+"\n", rr.Body.String())
	})
}

func TestCreateCar(t *testing.T) {
	validate := validator.New()
	handlers.SetValidator(validate)

	mockCarService := &MockCarService{
		CreateCarFunc: func(car *models.Car, fileData []byte, fileName string) (interface{}, error) {
			if car.Make == "Toyota" {
				return map[string]string{"message": "Car created successfully"}, nil
			}
			return nil, assert.AnError
		},
	}

	handlers.SetCarService(mockCarService)

	t.Run("valid car data", func(t *testing.T) {
		// Creating a multipart request with valid car data
		fileContent := []byte("fake image data")
		req, err := newMultipartRequest("POST", "/cars", map[string]string{
			"make":  "Toyota",
			"model": "Corolla",
			"year":  "2020",
			"price": "20000",
		}, "picture", fileContent)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.CreateCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusOK, rr.Code)

		var result map[string]string
		json.NewDecoder(rr.Body).Decode(&result)
		expected := map[string]string{"message": "Car created successfully"}
		assert.Equal(t, expected, result)
	})

	t.Run("invalid car data", func(t *testing.T) {
		// Creating a multipart request with invalid car data (empty make)
		fileContent := []byte("fake image data")
		req, err := newMultipartRequest("POST", "/cars", map[string]string{
			"make":  "",
			"model": "Corolla",
			"year":  "2020",
			"price": "20000",
		}, "picture", fileContent)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.CreateCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Make is required")
	})

	t.Run("service error", func(t *testing.T) {
		// Simulating a service error
		fileContent := []byte("fake image data")
		req, err := newMultipartRequest("POST", "/cars", map[string]string{
			"make":  "Honda",
			"model": "Civic",
			"year":  "2020",
			"price": "20000",
		}, "picture", fileContent)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.CreateCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, assert.AnError.Error()+"\n", rr.Body.String())
	})
}

func TestUpdateCar(t *testing.T) {
	validate := validator.New()
	handlers.SetValidator(validate)

	mockCarService := &MockCarService{
		UpdateCarFunc: func(id primitive.ObjectID, car *models.Car, fileData []byte, fileName string) (interface{}, error) {
			if id.Hex() == "60c72b2f9b1e8b3e0c6fc1c1" {
				return map[string]string{"message": "Car updated successfully"}, nil
			}
			return nil, assert.AnError
		},
	}

	handlers.SetCarService(mockCarService)

	t.Run("valid car data", func(t *testing.T) {
		// Creating a multipart request with valid car data
		fileContent := []byte("fake image data")
		req, err := newMultipartRequest("PUT", "/cars/60c72b2f9b1e8b3e0c6fc1c1", map[string]string{
			"make":  "Toyota",
			"model": "Corolla",
			"year":  "2020",
			"price": "20000",
		}, "picture", fileContent)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60c72b2f9b1e8b3e0c6fc1c1"})
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.UpdateCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusOK, rr.Code)

		var result map[string]string
		json.NewDecoder(rr.Body).Decode(&result)
		expected := map[string]string{"message": "Car updated successfully"}
		assert.Equal(t, expected, result)
	})

	t.Run("invalid car ID", func(t *testing.T) {
		// Creating a multipart request with an invalid car ID
		fileContent := []byte("fake image data")
		req, err := newMultipartRequest("PUT", "/cars/invalid", map[string]string{
			"make":  "Toyota",
			"model": "Corolla",
			"year":  "2020",
			"price": "20000",
		}, "picture", fileContent)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.UpdateCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid car ID\n", rr.Body.String())
	})

	t.Run("invalid car data", func(t *testing.T) {
		// Creating a multipart request with invalid car data (empty make)
		fileContent := []byte("fake image data")
		req, err := newMultipartRequest("PUT", "/cars/60c72b2f9b1e8b3e0c6fc1c1", map[string]string{
			"make":  "",
			"model": "Corolla",
			"year":  "2020",
			"price": "20000",
		}, "picture", fileContent)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60c72b2f9b1e8b3e0c6fc1c1"})
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.UpdateCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Make is required")
	})

	t.Run("service error", func(t *testing.T) {
		// Simulating a service error
		fileContent := []byte("fake image data")
		req, err := newMultipartRequest("PUT", "/cars/60c72b2f9b1e8b3e0c6fc1c2", map[string]string{
			"make":  "Toyota",
			"model": "Corolla",
			"year":  "2020",
			"price": "20000",
		}, "picture", fileContent)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60c72b2f9b1e8b3e0c6fc1c2"})
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.UpdateCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, assert.AnError.Error()+"\n", rr.Body.String())
	})
}

func TestDeleteCar(t *testing.T) {
	// Mocking the car service with a DeleteCar function
	mockCarService := &MockCarService{
		DeleteCarFunc: func(id primitive.ObjectID) (interface{}, error) {
			if id.Hex() == "60c72b2f9b1e8b3e0c6fc1c1" {
				return map[string]string{"message": "Car deleted successfully"}, nil
			}
			return nil, assert.AnError
		},
	}

	// Setting the mock service in the handler
	handlers.SetCarService(mockCarService)

	t.Run("valid car ID", func(t *testing.T) {
		// Creating a request with a valid car ID
		req := httptest.NewRequest("DELETE", "/cars/60c72b2f9b1e8b3e0c6fc1c1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "60c72b2f9b1e8b3e0c6fc1c1"})
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.DeleteCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusOK, rr.Code)
		var result map[string]string
		json.NewDecoder(rr.Body).Decode(&result)
		expected := map[string]string{"message": "Car deleted successfully"}
		assert.Equal(t, expected, result)
	})

	t.Run("invalid car ID", func(t *testing.T) {
		// Creating a request with an invalid car ID
		req := httptest.NewRequest("DELETE", "/cars/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.DeleteCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid car ID\n", rr.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Creating a request with a car ID that triggers a service error
		req := httptest.NewRequest("DELETE", "/cars/60c72b2f9b1e8b3e0c6fc1c2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "60c72b2f9b1e8b3e0c6fc1c2"})
		rr := httptest.NewRecorder()

		// Calling the handler
		handlers.DeleteCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, assert.AnError.Error()+"\n", rr.Body.String())
	})
}

func TestReserveCar(t *testing.T) {
	// Setting up the validator and mock service
	validate := validator.New()
	handlers.SetValidator(validate)

	mockCarService := &MockCarService{
		ReserveCarFunc: func(id primitive.ObjectID, customer models.Customer) (interface{}, error) {
			if customer.FullName == "John Doe" {
				return map[string]string{"message": "Car reserved successfully"}, nil
			}
			return nil, assert.AnError
		},
	}

	handlers.SetCarService(mockCarService)

	t.Run("valid reservation data", func(t *testing.T) {
		// Creating a valid customer object and request
		customer := models.Customer{
			FullName:    "John Doe",
			Email:       "john.doe@example.com",
			PhoneNumber: "1234567890",
		}
		body, _ := json.Marshal(customer)
		req, err := http.NewRequest("POST", "/cars/60d5f60e4f1c000088aa828e/reserve", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60d5f60e4f1c000088aa828e"})

		rr := httptest.NewRecorder()
		handlers.ReserveCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusOK, rr.Code)
		var result map[string]string
		json.NewDecoder(rr.Body).Decode(&result)
		expected := map[string]string{"message": "Car reserved successfully"}
		assert.Equal(t, expected, result)
	})

	t.Run("invalid reservation data", func(t *testing.T) {
		// Creating a customer object with invalid data
		customer := models.Customer{
			FullName:    "Name",
			Email:       "invalid-email",
			PhoneNumber: "12345",
		}
		body, _ := json.Marshal(customer)
		req, err := http.NewRequest("POST", "/cars/60d5f60e4f1c000088aa828e/reserve", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60d5f60e4f1c000088aa828e"})

		rr := httptest.NewRecorder()
		handlers.ReserveCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Email is not a valid email address")
	})

	t.Run("service error", func(t *testing.T) {
		// Creating a valid customer object with an ID that triggers a service error
		customer := models.Customer{
			FullName:    "Jane Doe",
			Email:       "jane.doe@example.com",
			PhoneNumber: "0987654321",
		}
		body, _ := json.Marshal(customer)
		req, err := http.NewRequest("POST", "/cars/60d5f60e4f1c000088aa828e/reserve", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60d5f60e4f1c000088aa828e"})

		rr := httptest.NewRecorder()
		handlers.ReserveCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, assert.AnError.Error()+"\n", rr.Body.String())
	})

	t.Run("invalid car ID", func(t *testing.T) {
		// Creating a valid customer object with an invalid car ID
		customer := models.Customer{
			FullName:    "John Doe",
			Email:       "john.doe@example.com",
			PhoneNumber: "1234567890",
		}
		body, _ := json.Marshal(customer)
		req, err := http.NewRequest("POST", "/cars/invalid-id/reserve", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-id"})

		rr := httptest.NewRecorder()
		handlers.ReserveCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid car ID\n", rr.Body.String())
	})
}

func TestCancelReservation(t *testing.T) {
	// Mocking the car service with a CancelReservation function
	mockCarService := &MockCarService{
		CancelReservationFunc: func(id primitive.ObjectID) (interface{}, error) {
			if id.Hex() == "60d5f60e4f1c000088aa828e" {
				return map[string]string{"message": "Reservation cancelled successfully"}, nil
			}
			return nil, assert.AnError
		},
	}

	// Setting the mock service in the handler
	handlers.SetCarService(mockCarService)

	t.Run("valid cancel reservation", func(t *testing.T) {
		// Creating a request with a valid reservation ID
		req, err := http.NewRequest("DELETE", "/cars/60d5f60e4f1c000088aa828e/cancel", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60d5f60e4f1c000088aa828e"})

		rr := httptest.NewRecorder()
		handlers.CancelReservation(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusOK, rr.Code)
		var result map[string]string
		json.NewDecoder(rr.Body).Decode(&result)
		expected := map[string]string{"message": "Reservation cancelled successfully"}
		assert.Equal(t, expected, result)
	})

	t.Run("service error", func(t *testing.T) {
		// Creating a request with a reservation ID that triggers a service error
		req, err := http.NewRequest("DELETE", "/cars/60d5f60e4f1c000088aa828f/cancel", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60d5f60e4f1c000088aa828f"})

		rr := httptest.NewRecorder()
		handlers.CancelReservation(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, assert.AnError.Error()+"\n", rr.Body.String())
	})

	t.Run("invalid car ID", func(t *testing.T) {
		// Creating a request with an invalid car ID
		req, err := http.NewRequest("DELETE", "/cars/invalid-id/cancel", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-id"})

		rr := httptest.NewRecorder()
		handlers.CancelReservation(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid car ID\n", rr.Body.String())
	})
}

func TestSellCar(t *testing.T) {
	// Setting up the validator and mock service
	validate := validator.New()
	handlers.SetValidator(validate)

	mockCarService := &MockCarService{
		SellCarFunc: func(id primitive.ObjectID, customer models.Customer) (interface{}, error) {
			if id.Hex() == "60d5f60e4f1c000088aa828e" {
				return map[string]string{"message": "Car sold successfully"}, nil
			}
			return nil, assert.AnError
		},
	}

	handlers.SetCarService(mockCarService)

	t.Run("valid sell car request", func(t *testing.T) {
		// Creating a valid customer object and request
		customer := models.Customer{
			FullName:    "John Doe",
			Email:       "john.doe@example.com",
			PhoneNumber: "1234567890",
		}
		customerJSON, _ := json.Marshal(customer)
		req, err := http.NewRequest("POST", "/cars/60d5f60e4f1c000088aa828e/sell", bytes.NewBuffer(customerJSON))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60d5f60e4f1c000088aa828e"})

		rr := httptest.NewRecorder()
		handlers.SellCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusOK, rr.Code)
		var result map[string]string
		json.NewDecoder(rr.Body).Decode(&result)
		expected := map[string]string{"message": "Car sold successfully"}
		assert.Equal(t, expected, result)
	})

	t.Run("invalid customer data", func(t *testing.T) {
		// Creating a customer object with invalid data
		customer := models.Customer{
			FullName:    "Name",
			Email:       "invalid-email",
			PhoneNumber: "123",
		}
		customerJSON, _ := json.Marshal(customer)
		req, err := http.NewRequest("POST", "/cars/60d5f60e4f1c000088aa828e/sell", bytes.NewBuffer(customerJSON))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60d5f60e4f1c000088aa828e"})

		rr := httptest.NewRecorder()
		handlers.SellCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Email is not a valid email address")
	})

	t.Run("service error", func(t *testing.T) {
		// Creating a valid customer object with an ID that triggers a service error
		customer := models.Customer{
			FullName:    "John Doe",
			Email:       "john.doe@example.com",
			PhoneNumber: "1234567890",
		}
		customerJSON, _ := json.Marshal(customer)
		req, err := http.NewRequest("POST", "/cars/60d5f60e4f1c000088aa828f/sell", bytes.NewBuffer(customerJSON))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "60d5f60e4f1c000088aa828f"})

		rr := httptest.NewRecorder()
		handlers.SellCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, assert.AnError.Error()+"\n", rr.Body.String())
	})

	t.Run("invalid car ID", func(t *testing.T) {
		// Creating a valid customer object with an invalid car ID
		customer := models.Customer{
			FullName:    "John Doe",
			Email:       "john.doe@example.com",
			PhoneNumber: "1234567890",
		}
		customerJSON, _ := json.Marshal(customer)
		req, err := http.NewRequest("POST", "/cars/invalid-id/sell", bytes.NewBuffer(customerJSON))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-id"})

		rr := httptest.NewRecorder()
		handlers.SellCar(rr, req)

		// Checking the response status and body
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid car ID\n", rr.Body.String())
	})
}

package serviceAdapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServiceUsecase is a mock implementation of the ServiceUsecase interface.
type MockServiceUsecase struct {
	mock.Mock
}

func (m *MockServiceUsecase) CreateService(service *Entities.Service) error {
	args := m.Called(service)
	return args.Error(0)
}

func (m *MockServiceUsecase) GetAll() (*[]Entities.Service, error) {
	args := m.Called()
	return args.Get(0).(*[]Entities.Service), args.Error(1)
}

func (m *MockServiceUsecase) GetByID(id *string) (*Entities.Service, error) {
	args := m.Called(id)
	return args.Get(0).(*Entities.Service), args.Error(1)
}

func (m *MockServiceUsecase) UpdateService(service *Entities.Service) error {
	args := m.Called(service)
	return args.Error(0)
}

func (m *MockServiceUsecase) DeleteService(id *string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockServiceUsecase) GetPackageByServiceID(serviceID *string) (*[]Entities.Package, error) {
	args := m.Called(serviceID)
	return args.Get(0).(*[]Entities.Package), args.Error(1)
}

func TestCreateService(t *testing.T) {
	// Set test mode to bypass real cloud storage.
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")

	mockUsecase := new(MockServiceUsecase)
	handler := NewServiceHandler(mockUsecase)
	app := fiber.New()
	app.Post("/services", handler.CreateService)

	// Prepare input using model.CreateServiceInput (which the handler expects)
	input := model.CreateServiceInput{
		Name:        "Test Service",
		Description: "Test Description",
		Rate:        5,
		Categories:  []string{"550e8400-e29b-41d4-a716-446655440000"},
		Packages: []model.PackageInput{
			{
				Name:        "Basic Package",
				Item:        "Item 1",
				Price:       100,
				Description: "Basic package description",
			},
		},
	}

	// Build a multipart form request:
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add the JSON field
	jsonField, err := writer.CreateFormField("json")
	assert.NoError(t, err)
	jsonBytes, err := json.Marshal(input)
	assert.NoError(t, err)
	_, err = jsonField.Write(jsonBytes)
	assert.NoError(t, err)

	// Add an attachment file field (simulate a file upload)
	fileField, err := writer.CreateFormFile("attachments", "test.txt")
	assert.NoError(t, err)
	_, err = fileField.Write([]byte("dummy file content"))
	assert.NoError(t, err)

	// Close the multipart writer to set the terminating boundary.
	err = writer.Close()
	assert.NoError(t, err)

	// Since the handler converts the input into an Entities.Service and later appends attachments,
	// we match on key fields (and ensure at least one attachment was added).
	mockUsecase.
		On("CreateService", mock.MatchedBy(func(s *Entities.Service) bool {
			return s.Name == input.Name &&
				s.Description == input.Description &&
				s.Rate == input.Rate &&
				len(s.Categories) == len(input.Categories) &&
				len(s.Packages) == len(input.Packages) &&
				len(s.Attachments) > 0
		})).
		Return(nil).Once()

	req := httptest.NewRequest("POST", "/services", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Read and log the response body (including error, if any)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Logf("TestCreateService response: %s", string(bodyBytes))

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestGetAllServices(t *testing.T) {
	mockUsecase := new(MockServiceUsecase)
	handler := NewServiceHandler(mockUsecase)

	app := fiber.New()
	app.Get("/services", handler.GetAllServices)

	tests := []struct {
		name           string
		mockServices   []Entities.Service
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			mockServices:   []Entities.Service{{Name: "Service 1"}, {Name: "Service 2"}},
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Internal Error",
			mockServices:   []Entities.Service{},
			mockError:      errors.New("database error"),
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.On("GetAll").Return(&tt.mockServices, tt.mockError).Once()

			req := httptest.NewRequest("GET", "/services", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Read and log the response body
			bodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			t.Logf("TestGetAllServices (%s) response: %s", tt.name, string(bodyBytes))

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetByServiceID(t *testing.T) {
	mockUsecase := new(MockServiceUsecase)
	handler := NewServiceHandler(mockUsecase)

	app := fiber.New()
	app.Get("/services/:id", handler.GetByServiceID)

	validID := "550e8400-e29b-41d4-a716-446655440000"
	service := &Entities.Service{ID: uuid.MustParse(validID), Name: "Test Service"}

	tests := []struct {
		name           string
		serviceID      string
		mockService    *Entities.Service
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			serviceID:      validID,
			mockService:    service,
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Not Found",
			serviceID:      validID,
			mockService:    &Entities.Service{},
			mockError:      errors.New("service not found"),
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.On("GetByID", &tt.serviceID).Return(tt.mockService, tt.mockError).Once()

			req := httptest.NewRequest("GET", "/services/"+tt.serviceID, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Read and log the response body
			bodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			t.Logf("TestGetByServiceID (%s) response: %s", tt.name, string(bodyBytes))

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetPackagebyServiceID(t *testing.T) {
	mockUsecase := new(MockServiceUsecase)
	handler := NewServiceHandler(mockUsecase)

	app := fiber.New()
	// Updated route to match the handler: /services/:id/packages
	app.Get("/services/:id/packages", handler.GetPackagesbyServiceID)

	validID := "550e8400-e29b-41d4-a716-446655440000"
	packages := []Entities.Package{{ID: uuid.MustParse(validID), Name: "Test Package"}}

	tests := []struct {
		name           string
		serviceID      string
		mockPackages   *[]Entities.Package
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			serviceID:      validID,
			mockPackages:   &packages,
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Not Found",
			serviceID:      validID,
			mockPackages:   &[]Entities.Package{},
			mockError:      errors.New("packages not found"),
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.
				On("GetPackageByServiceID", &tt.serviceID).
				Return(tt.mockPackages, tt.mockError).Once()

			req := httptest.NewRequest("GET", "/services/"+tt.serviceID+"/packages", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Read and log the response body
			bodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			t.Logf("TestGetPackagebyServiceID (%s) response: %s", tt.name, string(bodyBytes))

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUpdateService(t *testing.T) {
	mockUsecase := new(MockServiceUsecase)
	handler := NewServiceHandler(mockUsecase)
	app := fiber.New()
	// Using PUT as defined in the handler
	app.Put("/services/:id", handler.UpdateService)

	// Prepare a service update payload.
	serviceToUpdate := Entities.Service{
		Name:        "Updated Service",
		Description: "Updated Description",
		Rate:        10,
	}
	bodyBytesJSON, err := json.Marshal(serviceToUpdate)
	assert.NoError(t, err)

	id := uuid.New().String()
	// Expect that the service passed to UpdateService will have its ID set from the URL.
	mockUsecase.
		On("UpdateService", mock.MatchedBy(func(s *Entities.Service) bool {
			return s.ID.String() == id && s.Name == serviceToUpdate.Name
		})).
		Return(nil).Once()

	req := httptest.NewRequest("PUT", "/services/"+id, bytes.NewReader(bodyBytesJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Read and log the response body
	bodyResp, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Logf("TestUpdateService response: %s", string(bodyResp))

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestDeleteService(t *testing.T) {
	mockUsecase := new(MockServiceUsecase)
	handler := NewServiceHandler(mockUsecase)
	app := fiber.New()
	app.Delete("/services/:id", handler.DeleteService)

	id := uuid.New().String()
	mockUsecase.On("DeleteService", &id).Return(nil).Once()

	req := httptest.NewRequest("DELETE", "/services/"+id, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Read and log the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Logf("TestDeleteService response: %s", string(bodyBytes))

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

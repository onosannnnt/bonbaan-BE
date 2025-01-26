package serviceAdapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestCreateService(t *testing.T) {
	mockUsecase := new(MockServiceUsecase)
	handler := NewServiceHandler(mockUsecase)

	app := fiber.New()
	app.Post("/services", handler.CreateService)

	tests := []struct {
		name           string
		inputService   Entities.Service
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			inputService:   Entities.Service{Name: "Test Service"},
			mockError:      nil,
			expectedStatus: fiber.StatusCreated,
		},
		{
			name:           "Service Already Exists",
			inputService:   Entities.Service{Name: "Test Service"},
			mockError:      errors.New("service already exists"),
			expectedStatus: fiber.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.On("CreateService", &tt.inputService).Return(tt.mockError).Once()

			body, _ := json.Marshal(tt.inputService)
			req := httptest.NewRequest("POST", "/services", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockUsecase.AssertExpectations(t)
		})
	}
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
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockUsecase.AssertExpectations(t)
		})
	}
}

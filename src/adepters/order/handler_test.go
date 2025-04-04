package orderAdepter

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderUsecase ...

type MockOrderUsecase struct {
	mock.Mock
}


func (m *MockOrderUsecase) Insert(order *Entities.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderUsecase) GetAll() ([]*Entities.Order, error) {
	args := m.Called()
	return args.Get(0).([]*Entities.Order), args.Error(1)
}

func (m *MockOrderUsecase) GetOne(id *string) (*Entities.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*Entities.Order), args.Error(1)
}

func (m *MockOrderUsecase) Update(id *string, order *Entities.Order) error {
	args := m.Called(id, order)
	return args.Error(0)
}

func (m *MockOrderUsecase) Delete(id *string) error {
	args := m.Called(id)
	return args.Error(0)
}


func TestInsert(t *testing.T) {
    mockUsecase := new(MockOrderUsecase)
    handler := NewOrderHandler(mockUsecase)

    app := fiber.New()
    app.Post("/orders", handler.Insert)
    mockDetail := []byte(`{"detail": "hihi"}`)
    var orderDetail model.JSONB
    if err := json.Unmarshal(mockDetail, &orderDetail); err != nil {
        t.Fatalf("Failed to unmarshal mock detail: %v", err)
    }

    tests := []struct {
        name              string
        inputOrder        model.OrderInsertRequest
        mockError         error
        expectedStatus    int
        expectInsertCall  bool
    }{
        {
            name: "Success",
            inputOrder: model.OrderInsertRequest{
                CancellationReason: "reason",
                OrderDetail:        orderDetail,
                Note:               "note",
                Deadline:           "2023-12-31",
                UserID:             "123e4567-e89b-12d3-a456-426614174000",
                ServiceID:          "123e4567-e89b-12d3-a456-426614174001",
            },
            mockError:      nil,
            expectedStatus: fiber.StatusOK,
            expectInsertCall: true,
        },
        {
            name: "Invalid Date Format",
            inputOrder: model.OrderInsertRequest{
                Deadline: "invalid-date",
            },
            mockError:      nil,
            expectedStatus: fiber.StatusBadRequest,
            expectInsertCall: false,
        },
        {
            name: "Invalid UUID Format",
            inputOrder: model.OrderInsertRequest{
                Deadline: "2023-12-31",
                UserID:   "invalid-uuid",
            },
            mockError:      nil,
            expectedStatus: fiber.StatusBadRequest,
            expectInsertCall: false,
        },
        {
            name: "Insert Error",
            inputOrder: model.OrderInsertRequest{
                CancellationReason: "reason",
                OrderDetail:        orderDetail,
                Note:               "note",
                Deadline:           "2023-12-31",
                UserID:             "123e4567-e89b-12d3-a456-426614174000",
                ServiceID:          "123e4567-e89b-12d3-a456-426614174001",
            },
            mockError:      errors.New("insert error"),
            expectedStatus: fiber.StatusInternalServerError, // Updated status code
            expectInsertCall: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.expectInsertCall {
                mockUsecase.On("Insert", mock.AnythingOfType("*Entities.Order")).Return(tt.mockError).Once()
            }

            body, _ := json.Marshal(tt.inputOrder)
            req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
            req.Header.Set("Content-Type", "application/json")

            resp, err := app.Test(req)
            if err != nil {
                t.Fatalf("Failed to test request: %v", err)
            }

            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            mockUsecase.AssertExpectations(t)
            mockUsecase.ExpectedCalls = nil // Reset expectations for the next test case
            mockUsecase.Calls = nil         // Clear previous calls
        })
    }
}

func TestGetAll(t *testing.T) {
	mockUsecase := new(MockOrderUsecase)
	handler := NewOrderHandler(mockUsecase)

	app := fiber.New()
	app.Get("/orders", handler.GetAll)

	tests := []struct {
		name           string
		mockOrders     []*Entities.Order
		mockError      error
		expectedStatus int
	}{
		{
			name: "Success",
			mockOrders: []*Entities.Order{
				{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")},
				{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")},
			},
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Internal Error",
			mockOrders:     nil,
			mockError:      errors.New("internal error"),
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.On("GetAll").Return(tt.mockOrders, tt.mockError).Once()

			req := httptest.NewRequest("GET", "/orders", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetByID(t *testing.T) {
	mockUsecase := new(MockOrderUsecase)
	handler := NewOrderHandler(mockUsecase)

	app := fiber.New()
	app.Get("/orders/:id", handler.GetByID)

	tests := []struct {
		name           string
		orderID        string
		mockOrder      *Entities.Order
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			orderID:        "123e4567-e89b-12d3-a456-426614174000",
			mockOrder:      &Entities.Order{ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")},
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Internal Error",
			orderID:        "123e4567-e89b-12d3-a456-426614174000",
			mockOrder:      nil,
			mockError:      errors.New("internal error"),
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockOrder != nil {
				mockUsecase.On("GetOne", &tt.orderID).Return(tt.mockOrder, tt.mockError).Once()
			} else {
				mockUsecase.On("GetOne", &tt.orderID).Return(tt.mockOrder, tt.mockError).Once()
			}

			req := httptest.NewRequest("GET", "/orders/"+tt.orderID, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUpdate(t *testing.T) {
	mockUsecase := new(MockOrderUsecase)
	handler := NewOrderHandler(mockUsecase)

	app := fiber.New()
	app.Put("/orders/:id", handler.Update)

	tests := []struct {
		name           string
		orderID        string
		inputOrder     model.OrderInsertRequest
		mockError      error
		expectedStatus int
	}{
		{
			name:    "Success",
			orderID: "123e4567-e89b-12d3-a456-426614174000",
			inputOrder: model.OrderInsertRequest{
				StatusID: "123e4567-e89b-12d3-a456-426614174002",
			},
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "Invalid UUID Format",
			orderID: "invalid-uuid",
			inputOrder: model.OrderInsertRequest{
				StatusID: "123e4567-e89b-12d3-a456-426614174002",
			},
			mockError:      nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:    "Update Error",
			orderID: "123e4567-e89b-12d3-a456-426614174000",
			inputOrder: model.OrderInsertRequest{
				StatusID: "123e4567-e89b-12d3-a456-426614174002",
			},
			mockError:      errors.New("update error"),
			expectedStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Only call Update when it's a valid UUID
			if tt.expectedStatus != fiber.StatusBadRequest || (tt.expectedStatus == fiber.StatusBadRequest && tt.name != "Invalid UUID Format") {
				mockUsecase.On("Update", &tt.orderID, mock.AnythingOfType("*Entities.Order")).Return(tt.mockError).Once()
			}

			body, _ := json.Marshal(tt.inputOrder)
			req := httptest.NewRequest("PUT", "/orders/"+tt.orderID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestDelete(t *testing.T) {
	mockUsecase := new(MockOrderUsecase)
	handler := NewOrderHandler(mockUsecase)

	app := fiber.New()
	app.Delete("/orders/:id", handler.Delete)

	tests := []struct {
		name           string
		orderID        string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			orderID:        "123e4567-e89b-12d3-a456-426614174000",
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Delete Error",
			orderID:        "123e4567-e89b-12d3-a456-426614174000",
			mockError:      errors.New("delete error"),
			expectedStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.On("Delete", &tt.orderID).Return(tt.mockError).Once()

			req := httptest.NewRequest("DELETE", "/orders/"+tt.orderID, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockUsecase.AssertExpectations(t)
		})
	}
}

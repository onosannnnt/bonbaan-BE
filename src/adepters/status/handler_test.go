package statusAdapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStatusUsecase struct {
    mock.Mock
}

func (m *MockStatusUsecase) Insert(status *Entities.Status) error {
    args := m.Called(status)
    return args.Error(0)
}

func (m *MockStatusUsecase) GetStatusByID(id *string) (*Entities.Status, error) {
    args := m.Called(id)
    return args.Get(0).(*Entities.Status), args.Error(1)
}

func (m *MockStatusUsecase) GetStatusByName(name *string) (*Entities.Status, error) {
    args := m.Called(name)
    return args.Get(0).(*Entities.Status), args.Error(1)
}

func (m *MockStatusUsecase) GetAll() ([]*Entities.Status, error) {
    args := m.Called()
    return args.Get(0).([]*Entities.Status), args.Error(1)
}

func (m *MockStatusUsecase) Update(status *Entities.Status) error {
    args := m.Called(status)
    return args.Error(0)
}

func (m *MockStatusUsecase) Delete(id *string) error {
    args := m.Called(id)
    return args.Error(0)
}

func TestInsertStatus(t *testing.T) {
    mockUsecase := new(MockStatusUsecase)
    handler := NewStatusHandler(mockUsecase)

    app := fiber.New()
    app.Post("/status", handler.InsertStatus)

    tests := []struct {
        name           string
        inputStatus    Entities.Status
        mockError      error
        expectedStatus int
    }{
        {
            name:           "Success",
            inputStatus:    Entities.Status{Name: "Active"},
            mockError:      nil,
            expectedStatus: fiber.StatusCreated,
        },
        {
            name:           "Status Already Exists",
            inputStatus:    Entities.Status{Name: "Active"},
            mockError:      errors.New("this role already exists"),
            expectedStatus: fiber.StatusConflict,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUsecase.On("Insert", &tt.inputStatus).Return(tt.mockError).Once()

            body, _ := json.Marshal(tt.inputStatus)
            req := httptest.NewRequest("POST", "/status", bytes.NewReader(body))
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

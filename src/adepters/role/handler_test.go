package roleAdapter

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

type MockRoleUsecase struct {
    mock.Mock
}

func (m *MockRoleUsecase) InsertRole(role *Entities.Role) error {
    args := m.Called(role)
    return args.Error(0)
}

func (m *MockRoleUsecase) GetAll() (*[]Entities.Role, error) {
    args := m.Called()
    return args.Get(0).(*[]Entities.Role), args.Error(1)
}

func TestInsertRole(t *testing.T) {
    mockUsecase := new(MockRoleUsecase)
    handler := NewRoleHandler(mockUsecase)

    app := fiber.New()
    // Register the POST route for InsertRole
    app.Post("/roles", handler.InsertRole())

    tests := []struct {
        name           string
        inputRole      Entities.Role
        mockError      error
        expectedStatus int
    }{
        {
            name:           "Success",
            inputRole:      Entities.Role{Role: "admin"},
            mockError:      nil,
            expectedStatus: fiber.StatusCreated,
        },
        {
            name:           "Role Already Exists",
            inputRole:      Entities.Role{Role: "admin"},
            mockError:      errors.New("this role already exists"),
            expectedStatus: fiber.StatusConflict,
        },
        // {
        //     name:           "Missing Role",
        //     inputRole:      Entities.Role{Role: ""},
        //     mockError:      errors.New("Please fill all the require fields"),
        //     expectedStatus: fiber.StatusBadRequest,
        // },
		
			
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup the expected calls for this test case
            mockUsecase.On("InsertRole", &tt.inputRole).Return(tt.mockError).Once()

            body, _ := json.Marshal(tt.inputRole)
            req := httptest.NewRequest("POST", "/roles", bytes.NewReader(body))
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

func TestGetAll(t *testing.T) {
    mockUsecase := new(MockRoleUsecase)
    handler := NewRoleHandler(mockUsecase)

    app := fiber.New()
    // Register the GET route for GetAll
    app.Get("/roles", handler.GetAll())

    tests := []struct {
        name           string
        mockRoles      []Entities.Role
        mockError      error
        expectedStatus int
    }{
        {
            name:           "Success",
            mockRoles:      []Entities.Role{{Role: "admin"}, {Role: "user"}},
            mockError:      nil,
            expectedStatus: fiber.StatusOK,
        },
        {
            name:           "Internal Error",
            mockRoles:      nil,
            mockError:      errors.New("database error"),
            expectedStatus: fiber.StatusInternalServerError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup the expected calls for this test case
            mockUsecase.On("GetAll").Return(&tt.mockRoles, tt.mockError).Once()

            req := httptest.NewRequest("GET", "/roles", nil)
            resp, err := app.Test(req)
            if err != nil {
                t.Fatalf("Failed to test request: %v", err)
            }

            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            mockUsecase.AssertExpectations(t)
        })
    }
}



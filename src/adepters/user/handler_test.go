package userAdepter

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

// MockUserUsecase ...

type MockUserUsecase struct {
    mock.Mock
}

func (m *MockUserUsecase) Register(user *Entities.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserUsecase) Login(user *Entities.User) (*string, error) {
    args := m.Called(user)
    return args.Get(0).(*string), args.Error(1)
}

func (m *MockUserUsecase) Me(userID *string) (*Entities.User, error) {
    args := m.Called(userID)
    return args.Get(0).(*Entities.User), args.Error(1)
}

func (m *MockUserUsecase) ChangePassword(userID *string, req *model.ChangePasswordRequest) (*Entities.User, error) {
    args := m.Called(userID, req)
    return args.Get(0).(*Entities.User), args.Error(1)
}

func (m *MockUserUsecase) GetAll() (*[]Entities.User, error) {
    args := m.Called()
    return args.Get(0).(*[]Entities.User), args.Error(1)
}

func (m *MockUserUsecase) Update(user *model.UpdateRequest) (*Entities.User, error) {
    args := m.Called(user)
    return args.Get(0).(*Entities.User), args.Error(1)
}

func (m *MockUserUsecase) Delete(userID *string) error {
    args := m.Called(userID)
    return args.Error(0)
}

func TestRegister(t *testing.T) {
    mockUsecase := new(MockUserUsecase)
    handler := NewUserHandler(mockUsecase)

    app := fiber.New()
    app.Post("/register", handler.Register())

    tests := []struct {
        name           string
        inputUser      Entities.User
        mockError      error
        expectedStatus int
    }{
        {
            name:           "Success",
            inputUser:      Entities.User{Email: "test@example.com", Password: "password"},
            mockError:      nil,
            expectedStatus: fiber.StatusCreated,
        },
        {
            name:           "User Already Exists",
            inputUser:      Entities.User{Email: "test@example.com", Password: "password"},
            mockError:      errors.New("this account already exists"),
            expectedStatus: fiber.StatusConflict,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUsecase.On("Register", &tt.inputUser).Return(tt.mockError).Once()

            body, _ := json.Marshal(tt.inputUser)
            req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
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

func TestLogin(t *testing.T) {
    mockUsecase := new(MockUserUsecase)
    handler := NewUserHandler(mockUsecase)

    app := fiber.New()
    app.Post("/login", handler.Login)

    tests := []struct {
        name           string
        inputLogin     LoginRequest
        mockToken      *string
        mockError      error
        expectedStatus int
    }{
        {
            name:           "Success",
            inputLogin:     LoginRequest{EmailOrUsername: "test@example.com", Password: "password"},
            mockToken:      func() *string { s := "token"; return &s }(),
            mockError:      nil,
            expectedStatus: fiber.StatusOK,
        },
        {
            name:           "Invalid Credentials",
            inputLogin:     LoginRequest{EmailOrUsername: "test@example.com", Password: "wrongpassword"},
            mockToken:      nil,
            mockError:      errors.New("Invalid email, username or password"),
            expectedStatus: fiber.StatusUnauthorized,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUsecase.On("Login", mock.AnythingOfType("*Entities.User")).
                Return(tt.mockToken, tt.mockError).Once()

            body, _ := json.Marshal(tt.inputLogin)
            req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
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

func TestMe(t *testing.T) {
    mockUsecase := new(MockUserUsecase)
    handler := NewUserHandler(mockUsecase)

    app := fiber.New()
    app.Get("/me", handler.Me)

    // Use a valid UUID string
    validUUID := "123e4567-e89b-12d3-a456-426614174000"

    tests := []struct {
        name           string
        userID         string
        mockUser       *Entities.User
        mockError      error
        expectedStatus int
    }{
        {
            name:           "Success",
            userID:         validUUID,
            mockUser:       &Entities.User{ID: uuid.MustParse(validUUID), Email: "test@example.com"},
            mockError:      nil,
            expectedStatus: fiber.StatusOK,
        },
        {
            name:           "Internal Error",
            userID:         validUUID,
            mockUser:       nil,
            mockError:      errors.New("Internal Server Error"),
            expectedStatus: fiber.StatusInternalServerError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUsecase.On("Me", &tt.userID).Return(tt.mockUser, tt.mockError).Once()

            req := httptest.NewRequest("GET", "/me", nil)
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("UserID", tt.userID)

            resp, err := app.Test(req)
            if err != nil {
                t.Fatalf("Failed to test request: %v", err)
            }

            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            mockUsecase.AssertExpectations(t)
        })
    }
}

func TestChangePassword(t *testing.T) {
    mockUsecase := new(MockUserUsecase)
    handler := NewUserHandler(mockUsecase)

    app := fiber.New()
    app.Post("/change-password", handler.ChangePassword)

    validUUID := "123e4567-e89b-12d3-a456-426614174000"

    tests := []struct {
        name           string
        userID         string
        inputRequest   model.ChangePasswordRequest
        mockUser       *Entities.User
        mockError      error
        expectedStatus int
    }{
        {
            name:           "Success",
            userID:         validUUID,
            inputRequest:   model.ChangePasswordRequest{OldPassword: "oldpassword", NewPassword: "newpassword"},
            mockUser:       &Entities.User{ID: uuid.MustParse(validUUID), Email: "test@example.com"},
            mockError:      nil,
            expectedStatus: fiber.StatusOK,
        },
        {
            name:           "Internal Error",
            userID:         validUUID,
            inputRequest:   model.ChangePasswordRequest{OldPassword: "oldpassword", NewPassword: "newpassword"},
            mockUser:       nil,
            mockError:      errors.New("Internal Server Error"),
            expectedStatus: fiber.StatusInternalServerError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUsecase.On("ChangePassword", &tt.userID, &tt.inputRequest).Return(tt.mockUser, tt.mockError).Once()

            body, _ := json.Marshal(tt.inputRequest)
            req := httptest.NewRequest("POST", "/change-password", bytes.NewReader(body))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("UserID", tt.userID)

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
    mockUsecase := new(MockUserUsecase)
    handler := NewUserHandler(mockUsecase)

    app := fiber.New()
    app.Get("/users", handler.GetAll)

    // Use valid UUIDs
    userID1 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
    userID2 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

    tests := []struct {
        name           string
        mockUsers      []Entities.User
        mockError      error
        expectedStatus int
    }{
        {
            name: "Success",
            mockUsers: []Entities.User{
                {ID: userID1, Email: "test@example.com"},
                {ID: userID2, Email: "test2@example.com"},
            },
            mockError:      nil,
            expectedStatus: fiber.StatusOK,
        },
        {
            name:           "Internal Error",
            mockUsers:      nil,
            mockError:      errors.New("Internal Server Error"),
            expectedStatus: fiber.StatusInternalServerError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUsecase.On("GetAll").Return(&tt.mockUsers, tt.mockError).Once()

            req := httptest.NewRequest("GET", "/users", nil)
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
    mockUsecase := new(MockUserUsecase)
    handler := NewUserHandler(mockUsecase)

    app := fiber.New()
    app.Put("/update", handler.Update)

    validUUID := "123e4567-e89b-12d3-a456-426614174000"

    tests := []struct {
        name           string
        userID         string
        inputUser      model.UpdateRequest
        mockUser       *Entities.User
        mockError      error
        expectedStatus int
    }{
        {
            name:     "Success",
            userID:   validUUID,
            inputUser: model.UpdateRequest{ID:   validUUID,
                Email: "updated@example.com"},
            mockUser: &Entities.User{
                ID:    uuid.MustParse(validUUID),
                Email: "updated@example.com",
            },
            mockError:      nil,
            expectedStatus: fiber.StatusOK,
        },
        {
            name:     "Internal Error",
            userID:   validUUID,
            inputUser: model.UpdateRequest{ID:   validUUID,
                Email: "updated@example.com"},
            mockUser:  nil,
            mockError: errors.New("Internal Server Error"),
            expectedStatus: fiber.StatusInternalServerError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUsecase.On("Update", &tt.inputUser).Return(tt.mockUser, tt.mockError).Once()

            body, _ := json.Marshal(tt.inputUser)
            req := httptest.NewRequest("PUT", "/update", bytes.NewReader(body))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("UserID", tt.userID)
            
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
    mockUsecase := new(MockUserUsecase)
    handler := NewUserHandler(mockUsecase)

    app := fiber.New()
    app.Delete("/delete", handler.Delete)

    validUUID := "123e4567-e89b-12d3-a456-426614174000"

    tests := []struct {
        name           string
        userID         string
        mockError      error
        expectedStatus int
    }{
        {
            name:           "Success",
            userID:         validUUID,
            mockError:      nil,
            expectedStatus: fiber.StatusOK,
        },
        {
            name:           "Internal Error",
            userID:         validUUID,
            mockError:      errors.New("Internal Server Error"),
            expectedStatus: fiber.StatusInternalServerError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUsecase.On("Delete", &tt.userID).Return(tt.mockError).Once()

            req := httptest.NewRequest("DELETE", "/delete", nil)
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("UserID", tt.userID)

            resp, err := app.Test(req)
            if err != nil {
                t.Fatalf("Failed to test request: %v", err)
            }

            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            mockUsecase.AssertExpectations(t)
        })
    }
}








package serviceAdapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

    tests := []struct {
        name           string
        input          model.CreateServiceInput
        usecaseRet     error
        expectedStatus int
    }{
        {
            name: "Valid input with one package",
            input: model.CreateServiceInput{
                Name:        "Test Service",
                Description: "Test Description",
                Categories:  []string{"550e8400-e29b-41d4-a716-446655440000"},
                Packages: []model.PackageInput{
                    {
                        Name:        "Basic Package",
                        Item:        "Item 1",
                        Price:       100,
                        Description: "Basic package description",
                    },
                },
            },
            usecaseRet:     nil,
            expectedStatus: fiber.StatusCreated,
        },
        {
            name: "Valid input with two packages",
            input: model.CreateServiceInput{
                Name:        "Another Service",
                Description: "Another Description",
                Categories:  []string{"550e8400-e29b-41d4-a716-446655440000", "660e8400-e29b-41d4-a716-446655440111"},
                Packages: []model.PackageInput{
                    {
                        Name:        "Basic Package",
                        Item:        "Item 1",
                        Price:       120,
                        Description: "Basic package description",
                    },
                    {
                        Name:        "Premium Package",
                        Item:        "Item 2",
                        Price:       200,
                        Description: "Premium package description",
                    },
                },
            },
            usecaseRet:     nil,
            expectedStatus: fiber.StatusCreated,
        },
        {
            name: "Usecase failure",
            input: model.CreateServiceInput{
                Name:        "Failing Service",
                Description: "Should fail",
                Categories:  []string{"550e8400-e29b-41d4-a716-446655440000"},
                Packages: []model.PackageInput{
                    {
                        Name:        "Basic Package",
                        Item:        "Item 1",
                        Price:       80,
                        Description: "Basic package description",
                    },
                },
            },
            usecaseRet:     errors.New("creation error"),
            expectedStatus: fiber.StatusInternalServerError,
        },
		{
			name: "Invalid input with no packages",
			input: model.CreateServiceInput{
				Name:        "Invalid Service",
				Description: "No packages",
				Categories:  []string{"550e8400-e29b-41d4-a716-446655440000"},
				Packages:    []model.PackageInput{},
			},
			usecaseRet:     nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "Invalid input with Nagaive price no packages",
			input: model.CreateServiceInput{
				Name:        "Invalid Service",
				Description: "No packages",
				Categories:  []string{"550e8400-e29b-41d4-a716-446655440000"},
				Packages: []model.PackageInput{
					{
						Name:        "Basic Package",
						Item:        "Item 1",
						Price:       -999,
						Description: "Basic package description",
					},
				},
			},
			usecaseRet:     nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "Invalid input with no categories",
			input: model.CreateServiceInput{
				Name:        "Invalid Service",
				Description: "No categories",
				Categories:  []string{},
				Packages: []model.PackageInput{
					{
						Name:        "Basic Package",
						Item:        "Item 1",
						Price:       80,
						Description: "Basic package description",
					},
				},
			},
			usecaseRet:     nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "Invalid input with no name",
			input: model.CreateServiceInput{
				Name:        "",
				Description: "No name",
				Categories:  []string{"550e8400-e29b-41d4-a716-446655440000"},
				Packages: []model.PackageInput{
					{
						Name:        "Basic Package",
						Item:        "Item 1",
						Price:       80,
						Description: "Basic package description",
					},
				},
			},
			usecaseRet:     nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "Invalid input with no description",
			input: model.CreateServiceInput{
				Name:        "Invalid Service",
				Description: "",
				Categories:  []string{"550e8400-e29b-41d4-a716-446655440000"},
				Packages: []model.PackageInput{
					{
						Name:        "Basic Package",
						Item:        "Item 1",
						Price:       80,
						Description: "Basic package description",
					},
				},
			},
			usecaseRet:     nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "Invalid input with no items in package",
			input: model.CreateServiceInput{
				Name:        "Invalid Service",
				Description: "No items in package",
				Categories:  []string{"550e8400-e29b-41d4-a716-446655440000"},
				Packages: []model.PackageInput{
					{
						Name 	  : "Basic Package",
						Item 	  : "",
						Price 	  : 80,
						Description: "Basic package description",
					},
				},
			},
			usecaseRet:     nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			// 1. Partially invalid: one valid package, one invalid (negative price).
			//    If your system requires *all* packages to be valid, expect a bad request.
			name: "Mixed validity in packages",
			input: model.CreateServiceInput{
				Name:        "Mixed Package Service",
				Description: "Has one valid package and one invalid package",
				Rate:        3,
				Categories:  []string{"550e8400-e29b-41d4-a716-446655440000"},
				Packages: []model.PackageInput{
					{
						Name:        "Valid Package",
						Item:        "Item A",
						Price:       100,
						Description: "A valid package",
					},
					{
						Name:        "Invalid Package",
						Item:        "Item B",
						Price:       -50, // Negative -> invalid
						Description: "Should trigger a bad request",
					},
				},
			},
			usecaseRet:     nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			// 2. Large number of categories (e.g., testing upper limits of category handling).
			//    Adjust the count or response code based on your logic (could still be valid).
			name: "Service with large category list",
			input: model.CreateServiceInput{
				Name:        "Category Heavy Service",
				Description: "Contains a high volume of categories for stress-testing",
				Rate:        5,
				Categories: func() []string {
					// Example: Generate 50 category UUIDs (just repeated valid ones for illustration).
					categories := make([]string, 50)
					for i := 0; i < 50; i++ {
						categories[i] = "550e8400-e29b-41d4-a716-446655440000"
					}
					return categories
				}(),
				Packages: []model.PackageInput{
					{
						Name:        "Basic Package",
						Item:        "Item 1",
						Price:       150,
						Description: "Valid package",
					},
				},
			},
			usecaseRet:     nil,
			expectedStatus: fiber.StatusCreated, // or fiber.StatusBadRequest if your logic forbids too many categories
		},
		{
			// 3. Invalid UUID in categories. Should fail if your system strictly validates UUID format.
			name: "Invalid category UUID format",
			input: model.CreateServiceInput{
				Name:        "Invalid UUID Service",
				Description: "One category has an invalid UUID",
				Rate:        2,
				Categories: []string{
					"550e8400-e29b-41d4-a716-446655440000", // valid
					"not-a-valid-uuid",                    // invalid
				},
				Packages: []model.PackageInput{
					{
						Name:        "Basic Package",
						Item:        "Item 1",
						Price:       100,
						Description: "Basic package description",
					},
				},
			},
			usecaseRet:     nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			// 4. Repeated category IDs. Depending on your logic, you might treat duplicates as invalid or simply ignore duplicates.
			name: "Repeated categories",
			input: model.CreateServiceInput{
				Name:        "Repeated Categories Service",
				Description: "Contains the same category multiple times",
				Rate:        1,
				Categories: []string{
					"550e8400-e29b-41d4-a716-446655440000",
					"550e8400-e29b-41d4-a716-446655440000",
				},
				Packages: []model.PackageInput{
					{
						Name:        "Basic Package",
						Item:        "Item 1",
						Price:       100,
						Description: "Package with repeated categories",
					},
				},
			},
			usecaseRet:     nil,
			// Choose the status based on your system's handling of duplicates:
			expectedStatus: fiber.StatusCreated, // or fiber.StatusBadRequest if duplicates are not allowed
		},
		{
			// 5. Special/unicode characters in the name. Useful to test input sanitization or encoding issues.
			name: "Service name with special characters",
			input: model.CreateServiceInput{
				Name:        "Service ☺️ #$%!",
				Description: "Name includes emojis and special chars",
				Rate:        0,
				Categories: []string{
					"550e8400-e29b-41d4-a716-446655440000",
				},
				Packages: []model.PackageInput{
					{
						Name:        "Basic Package",
						Item:        "🎉",
						Price:       999,
						Description: "Package name with emoji",
					},
				},
			},
			usecaseRet:     nil,
			// If your system supports and properly encodes these characters, it should pass:
			expectedStatus: fiber.StatusCreated, // or fiber.StatusBadRequest if special chars are forbidden
		},
		{
			name: "Service with 50 packages",
			input: model.CreateServiceInput{
				Name:        "Test Service with 50 Packages",
				Description: "A large number of packages to stress test creation logic",
				Rate:        3,
				Categories:  []string{"550e8400-e29b-41d4-a716-446655440000"},
				Packages: func() []model.PackageInput {
					pkgs := make([]model.PackageInput, 50)
					for i := 0; i < 50; i++ {
						pkgs[i] = model.PackageInput{
							Name:        fmt.Sprintf("Package #%d", i+1),
							Item:        fmt.Sprintf("Item %d", i+1),
							Price:       100 + i,
							Description: fmt.Sprintf("This is package number %d.", i+1),
						}
					}
					return pkgs
				}(),
				Attachments: []string{"attachment_example.pdf"},
			},
			usecaseRet:     nil,
			expectedStatus: fiber.StatusCreated,
		},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUsecase := new(MockServiceUsecase)
            handler := NewServiceHandler(mockUsecase)
            app := fiber.New()
            app.Post("/services", handler.CreateService)

            // Build a multipart form request:
            var body bytes.Buffer
            writer := multipart.NewWriter(&body)

            // Write the JSON field value.
            jsonField, err := writer.CreateFormField("json")
            assert.NoError(t, err)
            jsonBytes, err := json.Marshal(tt.input)
            assert.NoError(t, err)
            _, err = jsonField.Write(jsonBytes)
            assert.NoError(t, err)

            // Add an attachment file field (simulate a file upload).
            fileField, err := writer.CreateFormFile("attachments", "test.txt")
            assert.NoError(t, err)
            _, err = fileField.Write([]byte("dummy file content"))
            assert.NoError(t, err)

            // Close the multipart writer to set the terminating boundary.
            err = writer.Close()
            assert.NoError(t, err)

            // Expect that the handler converts the input to an Entities.Service with at least one attachment.
            mockUsecase.
                On("CreateService", mock.MatchedBy(func(s *Entities.Service) bool {
                    return s.Name == tt.input.Name &&
                        s.Description == tt.input.Description &&
                        s.Rate == tt.input.Rate &&
                        len(s.Categories) == len(tt.input.Categories) &&
                        len(s.Packages) == len(tt.input.Packages) &&
                        len(s.Attachments) > 0
                })).
                Return(tt.usecaseRet).Once()

            req := httptest.NewRequest("POST", "/services", &body)
            req.Header.Set("Content-Type", writer.FormDataContentType())

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

			// Read and log the response body
			// bodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			// t.Logf("TestGetAllServices (%s) response: %s", tt.name, string(bodyBytes))

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
			// bodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			// t.Logf("TestGetByServiceID (%s) response: %s", tt.name, string(bodyBytes))

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
			// bodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			// t.Logf("TestGetPackagebyServiceID (%s) response: %s", tt.name, string(bodyBytes))

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
	// bodyResp, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	// t.Logf("TestUpdateService response: %s", string(bodyResp))

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
	// bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	// t.Logf("TestDeleteService response: %s", string(bodyBytes))

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

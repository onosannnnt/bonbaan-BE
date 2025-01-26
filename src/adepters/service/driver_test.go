package serviceAdapter

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// func setupTestDB(t *testing.T) *gorm.DB {
// 	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
// 	assert.NoError(t, err)
// 	err = db.AutoMigrate(&Entities.Service{})
// 	assert.NoError(t, err)
// 	return db
// }

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to open database: %v", err))
	}
	db.AutoMigrate(&Entities.Service{})
	return db
}

func TestServiceDriver_Create(t *testing.T) {
	db := setupTestDB()
	driver := NewServiceDriver(db)

	service := &Entities.Service{

		ID:          uuid.New(),
		Name:        "Test Service",
		Description: "Test Description",
		Rate:        5,
	}

	// Insert the service into the database
	err := driver.Insert(service)
	assert.NoError(t, err)

	var found Entities.Service
	fmt.Println(found.ID)
	// Query the service by ID
	err = db.First(&found, "id = ?", service.ID).Error
	assert.NoError(t, err)

	// Verify that the inserted service matches the expected values
	assert.Equal(t, service.Name, found.Name)
	assert.Equal(t, service.Description, found.Description)
}

func TestServiceDriver_GetAll(t *testing.T) {
	db := setupTestDB()
	driver := NewServiceDriver(db)
	services := []Entities.Service{
		{ID: uuid.New(), Name: "Service 1", Description: "Description 1", Rate: 5},
		{ID: uuid.New(), Name: "Service 2", Description: "Description 2", Rate: 4},
	}
	for _, s := range services {
		err := db.Create(&s).Error
		assert.NoError(t, err)
	}
	found, err := driver.GetAll()
	assert.NoError(t, err)
	assert.Len(t, *found, 3)
}

func TestServiceDriver_GetByID(t *testing.T) {
	db := setupTestDB()
	driver := NewServiceDriver(db)
	service := &Entities.Service{
		ID:          uuid.New(),
		Name:        "Test Service2",
		Description: "Test Description",
		Rate:        5,
	}
	err := db.Create(service).Error
	assert.NoError(t, err)
	id := service.ID.String()
	found, err := driver.GetByID(&id)
	assert.NoError(t, err)
	assert.Equal(t, service.Name, found.Name)
}

func TestServiceDriver_Update(t *testing.T) {
	db := setupTestDB()
	driver := NewServiceDriver(db)
	service := &Entities.Service{
		ID:          uuid.New(),
		Name:        "Test Service",
		Description: "Test Description",
		Rate:        5,
	}
	err := db.Create(service).Error
	assert.NoError(t, err)
	service.Name = "Updated Service"
	err = driver.Update(service)
	assert.NoError(t, err)
	var found Entities.Service
	err = db.First(&found, "id = ?", service.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Updated Service", found.Name)
}

func TestServiceDriver_Delete(t *testing.T) {
	db := setupTestDB()
	driver := NewServiceDriver(db)
	service := &Entities.Service{
		ID:          uuid.New(),
		Name:        "Test Service",
		Description: "Test Description",
		Rate:        5,
	}
	err := db.Create(service).Error
	assert.NoError(t, err)
	id := service.ID.String()
	err = driver.Delete(&id)
	assert.NoError(t, err)
	var found Entities.Service
	err = db.First(&found, "id = ?", service.ID).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

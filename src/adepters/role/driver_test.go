package roleAdapter

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
	db.AutoMigrate(&Entities.Role{})
	return db
  }

func TestRoleDriver_Insert(t *testing.T) {
	db := setupTestDB()
	driver := NewRoleDriver(db)

	role := &Entities.Role{
		ID: uuid.New(),
		Role: "Test Role",
	}

	err := driver.Insert(role)
	assert.NoError(t, err)

	var found Entities.Role
	err = db.First(&found, "id = ?", role.ID).Error
	assert.NoError(t, err)

	assert.Equal(t, role.Role, found.Role)

}

func TestRoleDriver_GetAll(t *testing.T) {
	db := setupTestDB()
	driver := NewRoleDriver(db)

	roles := []Entities.Role{
		{ID: uuid.New(), Role: "Role 1"},
		{ID: uuid.New(), Role: "Role 2"},
	}

	for _, r := range roles {
		err := db.Create(&r).Error
		assert.NoError(t, err)
	}

	found, err := driver.GetAll()
	assert.NoError(t, err)
	assert.Len(t, *found, 3)
}






















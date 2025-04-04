package statusAdapter

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        panic(fmt.Sprintf("Failed to open database: %v", err))
    }
    db.AutoMigrate(&Entities.Status{})
    return db
}

func TestStatusDriver_Insert(t *testing.T) {
    db := setupTestDB()
    driver := NewStatusDriver(db)

    status := &Entities.Status{
        ID:   uuid.New(),
        Name: "teststatus",
    }

    err := driver.Insert(status)
    assert.NoError(t, err)

    var found Entities.Status
    err = db.First(&found, "id = ?", status.ID).Error
    assert.NoError(t, err)
    assert.Equal(t, status.Name, found.Name)
}

func TestStatusDriver_FindStatusByID(t *testing.T) {
    db := setupTestDB()
    driver := NewStatusDriver(db)

    status := &Entities.Status{
        ID:   uuid.New(),
        Name: "teststatus",
    }

    err := db.Create(status).Error
    assert.NoError(t, err)

    idStr := status.ID.String()
    found, err := driver.FindStatusByID(&idStr)
    assert.NoError(t, err)
    assert.Equal(t, status.Name, found.Name)
}


func TestStatusDriver_FindStatusByName(t *testing.T) {
    db := setupTestDB()
    driver := NewStatusDriver(db)

    status := &Entities.Status{
        ID:   uuid.New(),
        Name: "teststatus",
    }

    err := db.Create(status).Error
    assert.NoError(t, err)

    found, err := driver.FindStatusByName(&status.Name)
    assert.NoError(t, err)
    assert.Equal(t, status.Name, found.Name)
}

func TestStatusDriver_FindAll(t *testing.T) {
    db := setupTestDB()
    driver := NewStatusDriver(db)

    statuses := []*Entities.Status{
        {ID: uuid.New(), Name: "status1"},
        {ID: uuid.New(), Name: "status2"},
    }

    for _, s := range statuses {
        err := db.Create(s).Error
        assert.NoError(t, err)
    }

    found, err := driver.FindAll()
    assert.NoError(t, err)
    assert.Len(t, found, 5)
}

func TestStatusDriver_Update(t *testing.T) {
    db := setupTestDB()
    driver := NewStatusDriver(db)

    status := &Entities.Status{
        ID:   uuid.New(),
        Name: "teststatus",
    }

    err := db.Create(status).Error
    assert.NoError(t, err)

    status.Name = "updatedstatus"
    err = driver.Update(status)
    assert.NoError(t, err)

    var found Entities.Status
    err = db.First(&found, "id = ?", status.ID).Error
    assert.NoError(t, err)
    assert.Equal(t, "updatedstatus", found.Name)
}

func TestStatusDriver_Delete(t *testing.T) {
    db := setupTestDB()
    driver := NewStatusDriver(db)

    status := &Entities.Status{
        ID:   uuid.New(),
        Name: "teststatus",
    }

    err := db.Create(status).Error
    assert.NoError(t, err)

    idStr := status.ID.String()
    err = driver.Delete(&idStr)
    assert.NoError(t, err)

    var found Entities.Status
    err = db.First(&found, "id = ?", status.ID).Error
    assert.Error(t, err)
    assert.Equal(t, gorm.ErrRecordNotFound, err)
}

    package orderAdepter

    import (
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
            panic("Failed to connect to test database")
        }
        db.AutoMigrate(&Entities.Order{}, &Entities.Status{}, &Entities.User{}, &Entities.Service{})
        return db
    }

    func TestOrderDriver_Insert(t *testing.T) {
        db := setupTestDB()
        driver := NewOrderDriver(db)

        status := &Entities.Status{ID: uuid.New(), Name: "pending"}
        user := &Entities.User{ID: uuid.New(), Username: "Test User"}
        service := &Entities.Service{ID: uuid.New(), Name: "Test Service"}

        db.Create(status)
        db.Create(user)
        db.Create(service)

        order := &Entities.Order{
            ID:        uuid.New(),
            UserID:    user.ID,
            ServiceID: service.ID,
            StatusID:  status.ID,
        }

        err := driver.Insert(order)
        assert.NoError(t, err)

        var found Entities.Order
        err = db.First(&found, "id = ?", order.ID).Error
        assert.NoError(t, err)
        assert.Equal(t, order.ID, found.ID)
    }

    func TestOrderDriver_GetDefaultStatus(t *testing.T) {
        db := setupTestDB()
        driver := NewOrderDriver(db)

        status := &Entities.Status{
            ID:   uuid.New(),
            Name: "pending",
        }
        db.Create(status)

        found, err := driver.GetDefaultStatus()
        assert.NoError(t, err)
        assert.Equal(t, "pending", found.Name)
    }

    func TestOrderDriver_GetAll(t *testing.T) {
        db := setupTestDB()
        driver := NewOrderDriver(db)

        status := &Entities.Status{ID: uuid.New(), Name: "pending"}
        user := &Entities.User{ID: uuid.New(), Username: "Test User"}
        service := &Entities.Service{ID: uuid.New(), Name: "Test Service"}

        db.Create(status)
        db.Create(user)
        db.Create(service)

        orders := []*Entities.Order{
            {ID: uuid.New(), UserID: user.ID, ServiceID: service.ID, StatusID: status.ID},
            {ID: uuid.New(), UserID: user.ID, ServiceID: service.ID, StatusID: status.ID},
        }
        for _, o := range orders {
            db.Create(o)
        }

        found, err := driver.GetAll()
        assert.NoError(t, err)
        assert.Len(t, found, 3)
    }

    func TestOrderDriver_GetByID(t *testing.T) {
        db := setupTestDB()
        driver := NewOrderDriver(db)

        status := &Entities.Status{ID: uuid.New(), Name: "pending"}
        user := &Entities.User{ID: uuid.New(), Username: "Test User"}
        service := &Entities.Service{ID: uuid.New(), Name: "Test Service"}

        db.Create(status)
        db.Create(user)
        db.Create(service)

        order := &Entities.Order{
            ID:        uuid.New(),
            UserID:    user.ID,
            ServiceID: service.ID,
            StatusID:  status.ID,
        }
        db.Create(order)

        idStr := order.ID.String()
        found, err := driver.GetByID(&idStr)
        assert.NoError(t, err)
        assert.Equal(t, order.ID, found.ID)
    }

    func TestOrderDriver_Update(t *testing.T) {
        db := setupTestDB()
        driver := NewOrderDriver(db)

        status := &Entities.Status{ID: uuid.New(), Name: "pending"}
        user := &Entities.User{ID: uuid.New(), Username: "Test User"}
        service := &Entities.Service{ID: uuid.New(), Name: "Test Service"}

        db.Create(status)
        db.Create(user)
        db.Create(service)

        order := &Entities.Order{
            ID:        uuid.New(),
            UserID:    user.ID,
            ServiceID: service.ID,
            StatusID:  status.ID,
        }
        db.Create(order)

        newStatus := &Entities.Status{ID: uuid.New(), Name: "completed"}
        db.Create(newStatus)

        order.StatusID = newStatus.ID
        idStr := order.ID.String()
        err := driver.Update(&idStr, order)
        assert.NoError(t, err)

        var found Entities.Order
        err = db.First(&found, "id = ?", order.ID).Error
        assert.NoError(t, err)
        assert.Equal(t, newStatus.ID, found.StatusID)
    }

    func TestOrderDriver_Delete(t *testing.T) {
        db := setupTestDB()
        driver := NewOrderDriver(db)

        status := &Entities.Status{ID: uuid.New(), Name: "pending"}
        user := &Entities.User{ID: uuid.New(), Username: "Test User"}
        service := &Entities.Service{ID: uuid.New(), Name: "Test Service"}

        db.Create(status)
        db.Create(user)
        db.Create(service)

        order := &Entities.Order{
            ID:        uuid.New(),
            UserID:    user.ID,
            ServiceID: service.ID,
            StatusID:  status.ID,
        }
        db.Create(order)

        idStr := order.ID.String()
        err := driver.Delete(&idStr)
        assert.NoError(t, err)

        var found Entities.Order
        err = db.First(&found, "id = ?", order.ID).Error
        assert.Error(t, err)
        assert.Equal(t, gorm.ErrRecordNotFound, err)
    }
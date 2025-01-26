package userAdepter

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// const projectDirName = "bonbaan-BE"

// func LoadEnv() {
//     // Find the current file path
//     _, filename, _, _ := runtime.Caller(0)
//     currentFileDir := filepath.Dir(filename)

//     // Walk up the directory tree to find the .env file
//     var configPath string
//     for {
//         configPath = filepath.Join(currentFileDir, ".env.test")
//         if _, err := os.Stat(configPath); !os.IsNotExist(err) {
//             break
//         }
//         parentDir := filepath.Dir(currentFileDir)
//         if parentDir == currentFileDir {
//             log.Fatalf(".env file not found")
//             os.Exit(-1)
//         }
//         currentFileDir = parentDir
//     }

//     // Load the .env file
//     err := godotenv.Load(configPath)
//     if err != nil {
//         log.Fatalf("Problem loading .env file: %v", err)
//         os.Exit(-1)
//     }
// }

func setupTestDB() *gorm.DB {

    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        panic(fmt.Sprintf("Failed to open database: %v", err))
    }
    db.AutoMigrate(&Entities.User{})

    return db
}

// func init() {
// 	LoadEnv()
// }


func TestUserDriver_Insert(t *testing.T) {
	
    db := setupTestDB()
    driver := NewUserDriver(db)

    user := &Entities.User{
		ID:       uuid.New(),
		Username: "testuser",
        Email:    "test1@gmail.com",
        Password: "password123",
    }

    err := driver.Insert(user)
    assert.NoError(t, err)

    var found Entities.User
    err = db.First(&found, "id = ?", user.ID).Error
    assert.NoError(t, err)
    assert.Equal(t, user.Username, found.Username)
    assert.Equal(t, user.Email, found.Email)
}


// Uncomment and update the following tests as needed

// func TestUserDriver_FindByEmailOrUsername(t *testing.T) {
// 	db := setupTestDB()
// 	driver := NewUserDriver(db)

// 	user := &Entities.User{
// 		ID: uuid.New(),
// 		Username: "testuser2",
// 		Email: "test2@test.com",
// 		Password: "password123",
// 	}

// 	err := db.Create(user).Error
// 	assert.NoError(t, err)

// 	found, err := driver.FindByEmailOrUsername(user)
// 	assert.NoError(t, err)
// 	assert.Equal(t, user.Username, found.Username)
// 	assert.Equal(t, user.Email, found.Email)
// }

// func TestUserDriver_FindByID(t *testing.T) {
// 	db := setupTestDB()
// 	driver := NewUserDriver(db)

// 	user := &Entities.User{
// 		ID: uuid.New(),
// 		Username: "testuser3",
// 		Email: "test3@test.com",
// 		Password: "password123",
// 	}

// 	err := db.Create(user).Error
// 	assert.NoError(t, err)

// 	id := user.ID.String()
// 	found, err := driver.FindByID(&id)
// 	assert.NoError(t, err)
// 	assert.Equal(t, user.Username, found.Username)
// }

// func TestUserDriver_Update(t *testing.T) {
// 	db := setupTestDB()
// 	driver := NewUserDriver(db)

// 	user := &Entities.User{
// 		ID: uuid.New(),
// 		Username: "testuser4",
// 		Email: "test4@test.com",
// 		Password: "password123",
// 	}

// 	err := db.Create(user).Error
// 	assert.NoError(t, err)

// 	user.Username = "updateduser"
// 	updated, err := driver.Update(user)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "updateduser", updated.Username)

// 	var found Entities.User
// 	err = db.First(&found, "id = ?", user.ID).Error
// 	assert.NoError(t, err)
// 	assert.Equal(t, "updateduser", found.Username)
// }

// func TestUserDriver_Delete(t *testing.T) {
// 	db := setupTestDB()
// 	driver := NewUserDriver(db)

// 	user := &Entities.User{
// 		ID: uuid.New(),
// 		Username: "testuser5",
// 		Email: "test5@test.com",
// 		Password: "password123",
// 	}

// 	err := db.Create(user).Error
// 	assert.NoError(t, err)

// 	id := user.ID.String()
// 	err = driver.Delete(&id)
// 	assert.NoError(t, err)

// 	var found Entities.User
// 	err = db.First(&found, "id = ?", user.ID).Error
// 	assert.Equal(t, gorm.ErrRecordNotFound, err)
// }

// func TestUserDriver_FindAll(t *testing.T) {
// 	db := setupTestDB()
// 	driver := NewUserDriver(db)

// 	users := []Entities.User{
// 		{ID: uuid.New(), Username: "user1", Email: "user1@test.com", Password: "pass1"},
// 		{ID: uuid.New(), Username: "user2", Email: "user2@test.com", Password: "pass2"},
// 	}

// 	for _, u := range users {
// 		err := db.Create(&u).Error
// 		assert.NoError(t, err)
// 	}

// 	found, err := driver.FindAll()
// 	assert.NoError(t, err)
// 	assert.Len(t, *found, len(users)+4)
// }

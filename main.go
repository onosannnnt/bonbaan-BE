package main

import (
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/onosannnnt/bonbaan-BE/src/Config"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	router "github.com/onosannnnt/bonbaan-BE/src/routers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Replace with your project and bucket details.
var (
	projectID  = "webpro-421315"
	bucketName = "webpro-421315.firebasestorage.app" // Typically your Firebase Storage bucket name
)


func main() {
	// image_upload()
	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", Config.DbHost, Config.DbPort, Config.DbUser, Config.DbPassword, Config.DbSchema)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	// Postgres install uuid-ossp extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		panic("failed to create uuid-ossp extension")
	}

	//Initialize Entities
	Entities.InitEntity(db)

	//check Entities in database

	if err != nil {
		panic("failed to connect database")
	}

	app := fiber.New(fiber.Config{
		BodyLimit: math.MaxInt64,
	})

	app.Use(cors.New())

	router.InitUserRouter(app, db)
	router.InitRoleRouter(app, db)
	router.InitStatusRouter(app, db)
	router.InitOrderRouter(app, db)
	router.ServiceRouter(app, db)
	router.InitPackageRouter(app, db)
	router.InitCategoryRouter(app, db)
	app.Listen(":" + Config.Port)
	
}

func Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello, World!"})
}



// // uploadImages handles the upload of multiple images and returns shareable download URLs.
// func uploadImages(w http.ResponseWriter, r *http.Request) {
	

// 	// Parse the multipart form with a max memory of 10 MB.
// 	if err := r.ParseMultipartForm(10 << 20); err != nil {
// 		http.Error(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Retrieve all files under the form field "images".
// 	files := r.MultipartForm.File["images"]
// 	if len(files) == 0 {
// 		http.Error(w, "No images provided", http.StatusBadRequest)
// 		return
// 	}
// 	ctx := context.Background()
// 	// Initialize the Cloud Storage client using your service account credentials.
// 	client, err := storage.NewClient(ctx, option.WithCredentialsFile("C:\\Users\\Mayoi\\Downloads\\webpro-421315-firebase-adminsdk-fbsvc-4467ef61f7.json"))
// 	if err != nil {
// 		http.Error(w, "Failed to create storage client: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer client.Close()

// 	// Slice to hold shareable URLs for all images.
// 	shareableURLs := make([]string, 0, len(files))

// 	// Iterate over each uploaded file.
// 	for _, fileHeader := range files {
// 		// Open the file.
// 		file, err := fileHeader.Open()
// 		if err != nil {
// 			http.Error(w, "Error opening file: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		// Generate a unique object name.
// 		// Using UnixNano ensures a high-resolution timestamp.
// 		objectName := fmt.Sprintf("images/%d_%s", time.Now().UnixNano(), fileHeader.Filename)

// 		// Generate a random download token.
// 		token := uuid.New().String()

// 		// Create a writer to upload the file to the storage bucket.
// 		wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
// 		// Set the metadata with the download token.
// 		wc.Metadata = map[string]string{
// 			"firebaseStorageDownloadTokens": token,
// 		}

// 		// Copy the file's content to Cloud Storage.
// 		if _, err = io.Copy(wc, file); err != nil {
// 			file.Close() // Close the file before returning.
// 			http.Error(w, "Failed to write file to bucket: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		// Close the file and the writer.
// 		file.Close()
// 		if err := wc.Close(); err != nil {
// 			http.Error(w, "Failed to close writer: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		// Construct the shareable URL.
// 		// The object name must be URL-encoded.
// 		shareableURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
// 			bucketName, url.QueryEscape(objectName), token)
// 		shareableURLs = append(shareableURLs, shareableURL)
// 	}

// 	// Return the shareable URLs as JSON.
// 	w.Header().Set("Content-Type", "application/json")
// 	response, err := json.Marshal(map[string][]string{"urls": shareableURLs})
// 	if err != nil {
// 		http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(response)
// }

// func image_upload() {
// 	http.HandleFunc("/upload", uploadImages)
// 	log.Println("Server is running on :8080...")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
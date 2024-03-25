package middleware

import (
	"fmt"
	"time"
)

func GenerateUniqueFilename(distribId string, ext string) (string, error) {
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("%s_%s.%s", distribId, timestamp, ext)
	return filename, nil
}

// func uploadHandler(objectFiles dto.ContactQueryFileMime) error {

// 	// Set size limit (optional)
// 	err := c.ParseMultipartForm(32 << 20) // Limit to 32 MB (adjust as needed)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	// Access uploaded file
// 	file, err := c.FormFile("file") // Replace "file" with your actual form field name
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
// 	}
// 	defer file.Close()

// 	// Generate unique filename
// 	filename, err := generateUniqueFilename("uploaded_file", filepath.Ext(file.Filename))
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate filename"})
// 	}

// 	// Define destination directory (ensure proper permissions)
// 	destinationDir := "./uploads" // Replace with your desired directory

// 	// Create destination path
// 	destinationPath := filepath.Join(destinationDir, filename)

// 	// Open destination file for writing
// 	destinationFile, err := os.Create(destinationPath)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create file"})
// 	}
// 	defer destinationFile.Close()

// 	// Write file content
// 	_, err = io.Copy(destinationFile, file)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
// 	}

// 	// Success response
// 	return c.JSON(fiber.StatusOK, fiber.Map{"message": "File uploaded successfully"})
// }

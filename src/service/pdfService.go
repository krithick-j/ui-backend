package service

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"ui-back-end/src/repositories"
	"ui-back-end/src/tmplts"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
)

func GenerateConsentForm(distrib_id string) (fiber.Map, int) {

	// Retrieve the user data from the database
	user, err := repositories.GetUserByID(distrib_id)
	if err.Error != nil {
		return fiber.Map{"error": "Error in getting distrib ID", "err": err.Error.Error()}, fiber.StatusInternalServerError
	}

	// Parse the HTML template from the string
	tmpl, res := template.New("DistribApplicationForm").Parse(tmplts.DistribApplicationFormTemplate)
	if res != nil {
		log.Println("Failed to parse template:", err)
		return fiber.Map{"error": "Failed to parse template"}, fiber.StatusInternalServerError
	}

	// Populate the template with user data
	var htmlContent bytes.Buffer
	if err := tmpl.Execute(&htmlContent, user); err != nil {
		log.Println("Failed to execute template:", err)
		return fiber.Map{"error": "Failed to execute template"}, fiber.StatusInternalServerError
	}

	// Generate PDF from HTML content
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	// pdf.MultiCell(0, 10, fmt.Sprintf("%v", htmlContent), "", "", false)
	pdf.Write(0, htmlContent.String())

	// Save the PDF file to the specified path
	filename := fmt.Sprintf("./assets/%s-idcard.pdf", distrib_id)
	res = pdf.OutputFileAndClose(filename)
	if res != nil {
		log.Println("Failed to save PDF:", res.Error())
		return fiber.Map{"error": "Failed to save PDF"}, fiber.StatusInternalServerError
	}

	// Return success response with file location
	return fiber.Map{"data": "DistribApplicationForm successfully generated", "file": filename}, fiber.StatusOK
}

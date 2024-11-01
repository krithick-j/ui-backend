package service

import (
	"bytes"
	"fmt"
	"html/template"
	"ui-back-end/src/repositories"
	"ui-back-end/src/tmplts"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

func GenerateConsentForm(distrib_id string, tx *gorm.DB) (fiber.Map, int) {

	// Retrieve the user data from the database
	user, err := repositories.GetUserByID(distrib_id, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetUserByID", "GenerateConsentForm", fiber.StatusInternalServerError, tx)
	}

	// Parse the HTML template from the string
	tmpl, err := template.New("DistribApplicationForm").Parse(tmplts.DistribApplicationFormTemplate)
	if err != nil {
		return utils.CommonErrorMessage(err, "Failed to parse template", fiber.StatusInternalServerError, tx)
	}

	// Populate the template with user data
	var htmlContent bytes.Buffer
	if err := tmpl.Execute(&htmlContent, user); err != nil {
		return utils.CommonErrorMessage(err, "Failed to execute template", fiber.StatusInternalServerError, tx)
	}

	// Generate PDF from HTML content
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	// pdf.MultiCell(0, 10, fmt.Sprintf("%v", htmlContent), "", "", false)
	pdf.Write(0, htmlContent.String())

	// Save the PDF file to the specified path
	savePath := fmt.Sprintf("./assets/%s-.pdf", distrib_id)
	err = pdf.OutputFileAndClose(savePath)
	if err != nil {
		return utils.CommonErrorMessage(err, "Failed to save PDF", fiber.StatusInternalServerError, tx)
	}
	fileName := fmt.Sprintf("media/%s-ack-letter.pdf", distrib_id)
	// Return success response with file location
	return fiber.Map{"data": "DistribApplicationForm successfully generated", "file": fileName}, fiber.StatusOK
}

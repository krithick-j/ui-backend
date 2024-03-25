package middleware

import (
	"fmt"
	"mime/multipart"
	"time"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
)

func GenerateUniqueFilename(distribId string, fileName string, ext string) (string, error) {
	println("ext-->", ext)
	if ext != "" {
		timestamp := time.Now().Format("2006-01-02T15-04-05")
		filename := fmt.Sprintf("%s_%s_%s%s", distribId, timestamp, fileName, ext)
		return filename, nil
	}
	return "", nil
}

func UploadFileToServer(fileName dto.ContactQueryFileForm, file models.ContactQueryFile) error {

	//Upload to contactQueryUploads
	// Success response
	return nil
}

func ParseForm(form map[string][]*multipart.FileHeader) (*dto.ContactQueryFileForm, error) {
	obj := new(dto.ContactQueryFileForm)
	if form["aadhaar_front"] != nil {
		obj.AadhaarFront = *form["aadhaar_front"][0]
	}
	if form["aadhaar_back"] != nil {
		obj.AadhaarBack = *form["aadhaar_back"][0]
	}
	if form["pan_card"] != nil {
		obj.PanCard = *form["pan_card"][0]
	}
	if form["passport_size"] != nil {
		obj.PassportSize = *form["passport_size"][0]
	}

	return obj, nil
}

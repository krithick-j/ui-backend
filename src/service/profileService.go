package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"ui-back-end/configs"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/go-pdf/fpdf"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
)

func IDCardFactory() *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	// Outer Rect
	pdf.Rect(5, 5, 180, 60, "D")
	// Inner Left Rect
	pdf.Rect(10, 10, 80, 50, "D")
	// Inner Right Rect
	pdf.Rect(100, 10, 80, 50, "D")
	//Disclaimer Box
	pdf.SetFillColor(234, 234, 245)
	pdf.RoundedRect(48, 38, 40, 15, 1, "1234", "DF")
	//Address Footer
	pdf.Rect(100, 45, 80, 15, "DF")
	pdf.Image("./assets/images/uilogo.png", 70, 12, 6, 0, false, "png", 0, "")
	pdf.Image("./assets/images/uilogo.png", 104, 49, 6, 0, false, "png", 0, "")
	// Inner Right Rect
	pdf.SetFillColor(0, 255, 0)
	pdf.RoundedRect(11, 42, 32, 6, 2, "1234", "DF")
	return pdf
}

func IDCardAddContent(pdf *fpdf.Fpdf, distrib_id string, userdata models.User) *fpdf.Fpdf {

	userphoto, _ := strings.CutPrefix(userdata.KYCPhoto, "media/")
	extension := filepath.Ext(userphoto)
	extension, _ = strings.CutPrefix(extension, ".")

	userphoto = filepath.Join("./assets", userphoto)
	if _, err := os.Stat(userphoto); errors.Is(err, os.ErrNotExist) {
		configs.Log.Errorf("File Does Not exit %s", userphoto)
	}
	if _, err := os.Stat("./assets/images/uilogo.png"); errors.Is(err, os.ErrNotExist) {
		configs.Log.Error("File Does Not exit", zap.String("image", "./assets/images/uilogo.png"))
	}
	pdf.Image(userphoto, 15, 10, 25, 0, false, extension, 0, "")
	var (
		currX float64 = 12
		currY float64 = 45
	)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "UI Distributor")
	currX += 8
	currY += 7
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Of")
	currX += -8
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Universe International - India")
	pdf.SetFont("Arial", "B", 10)
	currX += 36
	currY -= 34
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, userdata.Name)
	pdf.SetFont("Arial", "", 8)
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Dist ID: "+userdata.DistribID)
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Email: "+userdata.EmailAddress)
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Phone: +91-"+userdata.MobilePhoneNo)
	pdf.SetFont("Arial", "I", 8)
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• No Registration Fees"))
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• No Deposits"))
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• No Investments"))
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• No Job Offerings"))
	currX += 55
	currY -= 34
	pdf.SetFont("Arial", "", 6)
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Phone No: 9080296128")
	currY += 4
	//pdf.SetXY(currX, currY)
	//pdf.Cell(0, 0, "CIN: XXXXXXX-XXXXX")
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "GST: Applied")
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• Check Government ID Card for Proof"))
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• This card is not transferable"))
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• This card should be shown before presentation"))
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.Cell(0, 0, tr("In Case of Any Enquiry, Write to:"))
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "support@ui-network.com")
	//Footer
	currX += 10
	currY += 6
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 0, "Universe International - India")
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 6)
	pdf.Cell(0, 0, "KASI ARCADE FIRST FLOOR, VOC STREET,KAIKANKUPPAM")
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "RAMAPURAM, Chennai - PIN 600087 - TAMIL NADU")

	return pdf
}

func GenerateIDCard(distrib_id string) error {
	configs.Log.Info("Started Generating ID Card")
	userdata, _ := repositories.GetUserByID(distrib_id)
	// Create ID Card Layout
	pdf := IDCardFactory()
	pdf = IDCardAddContent(pdf, distrib_id, userdata)
	filename := fmt.Sprintf("./assets/%s-idcard.pdf", distrib_id)
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		configs.Log.Error("Error is ", zap.String("outfile", err.Error()))
		return err
	}
	return nil
}

func SendEmailCode(toMail string) error {
	otp := GenOPT()
	err := repositories.SaveOTP("email", toMail, otp)
	if err != nil {
		configs.Log.Errorln("Error on calling SaveOTP repositories fn from SendEmailCode service fn", err.Error())
		return err
	}
	msg := fmt.Sprintf("Your OTP for Email Verification is %s", otp)
	tomail := fmt.Sprintf("<%s>", toMail)
	SendPlainMail(tomail, "Your Email Verification OTP", msg)
	return nil
}

func VerifyEmailCode(otp string, email string) error {
	return repositories.CheckAndUpdateOTP("email", email, otp)
}

func SendPhoneCode(c *fiber.Ctx, phone string) error {
	otp := GenOPT()
	err := repositories.SaveOTP("phone", phone, otp)
	if err != nil {
		configs.Log.Errorln("Error on calling SaveOTP repositories fn from SendPhoneCode service fn", err.Error())
		return err
	}
	urlStr := "https://www.textguru.in/api/v22.0/?"

	// Data to be sent in the POST request
	data := url.Values{}
	data.Set("username", "jega.in")
	data.Set("password", "60423479")
	data.Set("source", "GSENTS")
	data.Set("dmobile", "91"+phone)
	data.Set("dlttempid", "1707171500974884924")
	data.Set("message", fmt.Sprintf("Dear Customer,\r\nThis is your OTP for Login %s for your mobile number verification On https://ui-network.com.\r\nGSENTS", otp))

	// Create a new POST request
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create an HTTP client and send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print the response status and body
	configs.Log.Infof("Response status: %s\n", resp.Status)
	configs.Log.Infof("Response body: %s\n", string(body))
	return nil
}

func VerifyPhoneCode(otp string, phone string) error {
	return repositories.CheckAndUpdateOTP("phone", phone, otp)
}

func KycUpload(c *fiber.Ctx, form *multipart.Form) error {
	user := models.User{}
	user.DistribID = form.Value["distrib_id"][0]
	for fs, fhs := range form.File {
		for _, fh := range fhs {
			extension := filepath.Ext(fh.Filename)
			fullPath := "./assets/" + user.DistribID + "-" + fs + extension
			mediapath := "media/" + user.DistribID + "-" + fs + extension
			err := c.SaveFile(fh, fullPath)
			if err != nil {
				return err
			}
			switch fs {
			case "aadhar":
				user.KYCAdhaar = mediapath
			case "consent":
				user.KYCConsentDoc = mediapath
			case "pan":
				user.KYCPAN = mediapath
			case "user-image":
				user.KYCPhoto = mediapath
			}
		}
	}
	user.KYCStatus = "pending"
	repositories.UpdateKyc(&user)
	return nil
}

func ApproveKYC(distrib_id string) error {
	user := models.User{}
	user.DistribID = distrib_id
	user.KYCStatus = "verified"

	repositories.UpdateKyc(&user)
	return nil
}

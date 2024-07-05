package service

import (
	"fmt"
	"ui-back-end/src/middleware"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/go-pdf/fpdf"
	"github.com/gofiber/fiber/v2"
)

func DistributorFormFactory(distribInformation models.User, referrerDistribInformation models.User) *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 1)
	pdf.AddPage()
	var (
		currX float64 = 75
		currY float64 = 8
	)

	// // Function to check if there is enough space for the table
	// checkSpaceForTable := func(pdf *fpdf.Fpdf, rowCount int, rowHeight float64, lineHeight float64, headerHeight float64) {
	// 	_, pageHeight := pdf.GetPageSize()
	// 	_, _, _, bottomMargin := pdf.GetMargins()
	// 	availableHeight := pageHeight - pdf.GetY() - bottomMargin - 10
	// 	requiredHeight := lineHeight + headerHeight + float64(int(rowHeight)*rowCount)
	// 	if requiredHeight > availableHeight {
	// 		pdf.AddPage()
	// 		currY = 10.0
	// 	}
	// }

	pdf.Image("assets/images/uilogo.png", currX, currY, 6, 0, false, "png", 0, "")
	//Universe International and Address
	pdf.SetFont("Arial", "B", 12)
	currX += 8
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Universe International")
	pdf.SetFont("Arial", "", 8)
	currX += 5
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "(Unit of GS ENTERPRISES)")
	currX += 4
	currY += 4
	currX += 12
	currY += 4

	//Line
	currX = 0
	currY += 5
	pdf.Line(currX, currY, currX+600, currY)

	//INVOICE
	pdf.SetFont("Arial", "BU", 13)
	currX += 66
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "DISTRIBUTOR APPLICATION FORM")

	//Distributor ID No:
	currX = 10
	currY += 9
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Distributor ID No:")
	currX += 25
	currY -= 2
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, distribInformation.DistribID, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	//Date
	currX += 55
	currY += 2

	pdf.SetXY(currX, currY)
	pdf.Cell(5, 0, "Date:")
	currX += 20
	currY -= 2
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(55, 4, middleware.FormatTimeByLocation(distribInformation.CreatedAt, "Asia/Kolkata", "02-01-2006 15:04:05"), "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currY += 10
	currX = 10
	pdf.SetXY(currX, currY)

	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 0, "REFERRER INFORMATION")
	pdf.SetTextColor(0, 0, 0)

	//Distributor ID No:
	currX = 10
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Referrer ID:")
	currX += 25
	currY -= 2
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, referrerDistribInformation.DistribID, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	//Date
	currX += 55
	currY += 2

	pdf.SetXY(currX, currY)
	pdf.Cell(5, 0, "Referrer Name:")
	currX += 25
	currY -= 2
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(55, 4, referrerDistribInformation.Name, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currY += 10
	currX = 10
	pdf.SetXY(currX, currY)

	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 0, "APPLICATION INFORMATION")
	pdf.SetTextColor(0, 0, 0)

	//Title
	currX = 10
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Title (If Individual) Mr./Mrs./Ms:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, referrerDistribInformation.DistribID, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currX += 78
	currY -= 5

	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Given Name:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, distribInformation.Name, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currY += 10
	//Title
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Mailing Address:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	mailingAddress := distribInformation.Address1 + distribInformation.Address2 + distribInformation.TownOrCity + distribInformation.StateOrProvince + distribInformation.Country + distribInformation.District + distribInformation.PinOrZipCode
	pdf.MultiCell(50, 4, mailingAddress, "1", "L", false)
	pdf.SetFont("Arial", "", 8)

	currX += 78
	currY -= 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Cheque Name:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, distribInformation.ChequeName, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	fmt.Println("xy-----------_>", pdf.GetY())
	currY = pdf.GetY() + 20
	//Title
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Shipping Address:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	shippingAddress := distribInformation.Address1 + distribInformation.Address2 + distribInformation.TownOrCity + distribInformation.StateOrProvince + distribInformation.Country + distribInformation.District + distribInformation.PinOrZipCode
	pdf.MultiCell(50, 4, shippingAddress, "1", "L", false)
	pdf.SetFont("Arial", "", 8)

	currX += 78
	currY -= 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Cheque Name:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, distribInformation.ChequeName, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	// checkSpaceForTable(pdf, len(bvDistributionTable), rowHeight, lineHeight, headerHeight)
	pdf.AddPage()
	currY = 10.0
	currX = 10.0
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "TERMS AND CONDITIONS")
	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(128, 128, 128)

	txt := "1. I have read, understood and agreed to be bound by all the terms and conditions set forth by Universe International Direct Selling (India) Pvt Ltd regarding this transaction."
	currY += 10.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    I have also read and agreed to comply with the Policies and Procedures as stated."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "2. Refund Policy : "
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")
	pdf.SetFont("Arial", "", 7)

	txt = "Distributors are hereby notified that Products are subject to the Company's Buy-Back Policy."
	currX += 25
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "3. The Company shall be obliged to buy-back any marketable product sold to a Customer/Distributor within fifteen (15) days from the date of invoice of the product after "
	currY += 6.0
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    withholding Tax Deducted at Source (TDS), Sales Incentive utilised, and other taxes if applicable, in accordance with its policies."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "4. The Customer/Distributor should raise a written request to the Company for the product refund within 15 days from the date of invoice. No refund requests will be entertained after"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "     15 days."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "5. Upon receipt and examination of the physical products, the final decision for a product refund rests with the Company."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "6. The Buy-Back Policy is only applicable for the cancellation of the full purchase order and upon the return of physical products to the company."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    In case of Combo products purchase or purchase order with multiple products, the distributor/customer must apply for refund conforming to all products of the said Combo set or "
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    purchase order."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "7. The company will not entertain a partial refund of selective products thereof. Subject to such products being in an unused state, accordingly the Company will process the refund "
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    of the payment made by the distributor/customer."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "Please send an email to "
	mail := `admin@ui-network.com`
	link := "https://mail.google.com/mail/?view=cm&fs=1&to=" + mail
	rest := " in case of further queries."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")
	currX += 27.8
	pdf.SetXY(currX, currY)
	pdf.SetTextColor(0, 0, 255) // Set text color to blue
	pdf.SetFont("Arial", "U", 7)
	pdf.WriteLinkString(0, mail, link)
	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(128, 128, 128)
	currX += 25.8
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, rest, "", 0, "L", true, 0, "")

	txt = "Please PRINT this receipt for your future reference. For questions and comments, please eMail: "
	currX = 10
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")
	mailtxt := "admin@ui-network.com"
	currX += 105
	pdf.SetXY(currX, currY)
	pdf.SetTextColor(0, 0, 255) // Set text color to blue
	pdf.SetFont("Arial", "U", 7)
	pdf.WriteLinkString(0, mailtxt, link)
	pdf.SetFont("Arial", "", 7)

	return pdf
}

func GenerateDistributorForm(distribId string) (string, int, error) {

	distributorInformation, res := repositories.GetUserByID(distribId)
	if res.Error != nil {
		return res.Error.Error(), fiber.StatusInternalServerError, res.Error
	}

	referrerDistribId, res := repositories.GetRefDistribIdByDistribId(distribId)
	if res.Error != nil {
		return res.Error.Error(), fiber.StatusInternalServerError, res.Error
	}

	referrerDistribInformation, res := repositories.GetUserByID(referrerDistribId)
	if res.Error != nil {
		return res.Error.Error(), fiber.StatusInternalServerError, res.Error
	}

	pdf := DistributorFormFactory(distributorInformation, referrerDistribInformation)

	filename := fmt.Sprintf("./distribApplicationForm/distribApplicationForm-%s.pdf", distribId)
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return err.Error(), fiber.StatusInternalServerError, err
	}

	return filename, fiber.StatusOK, err
}

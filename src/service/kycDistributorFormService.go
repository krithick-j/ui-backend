package service

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/middleware"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/go-pdf/fpdf"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DistributorFormFactory(distribInformation models.User, referrerDistribInformation models.User) *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 1)
	pdf.AddPage()
	var (
		currX float64 = 75
		currY float64 = 8
	)

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
	pdf.CellFormat(50, 4, referrerDistribInformation.Title, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currX += 95
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

	currY += 8
	//Title
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Mailing Address:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	mailingAddress := distribInformation.Address1 + ", " + distribInformation.Address2 + ", " + distribInformation.TownOrCity + ", " + distribInformation.StateOrProvince + ", " + distribInformation.Country + ", " + distribInformation.District + ", " + distribInformation.PinOrZipCode
	fmt.Printf("len(mailingAddress): %v\n", len(mailingAddress))
	if len(mailingAddress) > 120 {
		mailingAddress = distribInformation.Address1 + ", " + distribInformation.District + ", " + distribInformation.PinOrZipCode
	}
	if len(mailingAddress) > 120 {
		mailingAddress = distribInformation.Address1 + ", " + distribInformation.District + ", " + distribInformation.PinOrZipCode
	}
	pdf.MultiCell(90, 4, mailingAddress, "1", "L", false)
	pdf.SetFont("Arial", "", 8)

	currX += 95
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

	currY = pdf.GetY() + 10

	currY += 5

	//Shipping Address
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Shipping Address:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	shippingAddress := distribInformation.Address1 + ", " + distribInformation.Address2 + ", " + distribInformation.TownOrCity + ", " + distribInformation.StateOrProvince + ", " + distribInformation.Country + ", " + distribInformation.District + ", " + distribInformation.PinOrZipCode
	if len(shippingAddress) > 100 {
		shippingAddress = distribInformation.Address1 + ", " + distribInformation.District + ", " + distribInformation.PinOrZipCode
	}
	if len(shippingAddress) > 100 {
		shippingAddress = distribInformation.Address1 + ", " + distribInformation.District + ", " + distribInformation.PinOrZipCode
	}
	pdf.MultiCell(90, 4, shippingAddress, "1", "L", false)
	pdf.SetFont("Arial", "", 8)

	currX += 95
	currY -= 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Home Phone No & Mobile No:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	content := distribInformation.HomePhoneNo + " & " + distribInformation.MobilePhoneNo
	pdf.CellFormat(50, 4, content, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currY = pdf.GetY() + 8

	currX = 10

	//Valid ID No
	currX = 10
	currY += 8
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Valid ID No:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, distribInformation.ValidIdNo, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currX += 55
	currY -= 5

	//Email Address
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Email Address:")
	currX += 1
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, distribInformation.EmailAddress, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currX += 55
	currY -= 5
	//Nationality & DoB
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Nationality & Date Of Birth:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	content = distribInformation.Country + ", " + distribInformation.DateOfBirth
	pdf.CellFormat(50, 4, content, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currY += 5

	// Valid ID No
	currX = 10
	currY += 8
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Relationship to the Benificiary/Nominee:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, distribInformation.BeneficiaryRelationship, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currX += 95
	currY -= 5

	// Email Address
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Benificiary Name:")
	currX += 1
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, distribInformation.BenificiaryName, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currY += 8
	currX = 10
	//Pan Card
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "PAN Card:")
	currX += 20
	currY -= 2
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(100, 4, distribInformation.PanCard, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currY += 8
	currX = 10
	//Bank Name
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Bank Name:")
	currX += 20
	currY -= 2
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(100, 4, distribInformation.BankName, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currY += 8
	currX = 10
	//Bank Acct No
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Bank Acct No:")
	currX += 20
	currY -= 2
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(100, 4, distribInformation.BankAccNo, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)
	pdf.SetXY(currX, currY)

	currY += 8
	currX = 10
	//IFS Code
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "IFS Code:")
	currX += 20
	currY -= 2
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(100, 4, distribInformation.IFSCCode, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currY += 10
	currX = 10
	pdf.SetXY(currX, currY)

	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 0, "PREFERRED PLACEMENT INFORMATION")
	pdf.SetTextColor(0, 0, 0)

	currX = 10

	//Preferred Distributor ID
	currX = 10
	currY += 8
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Preferred Distributor ID:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)

	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, distribInformation.PreferredDistribId, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)
	currX += 55
	currY -= 5

	//Full Name
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Full Name:")
	currX += 1
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 4, distribInformation.PreferredDistribName, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currX += 55
	currY -= 5
	//Preferred Placement & Side
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(5, 0, "Preferred Placement & Side:")
	currX += 2
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	content = distribInformation.PreferredPlace + ", " + distribInformation.PreferredSide
	pdf.CellFormat(50, 4, content, "1", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 8)

	currY += 10
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 0, "TERMS AND CONDITIONS")
	pdf.SetTextColor(0, 0, 0)

	termsAndConditions := "1. For business entities, an authorised signatory of the company must sign this Distributor Application Form."
	currY += 10.0
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 0, termsAndConditions, "", 0, "L", true, 0, "")

	termsAndConditions = "2. You must be 21 years old and above to become a Distributor."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 0, termsAndConditions, "", 0, "L", true, 0, "")

	termsAndConditions = "3. By signing below. you certify and acknowledge that you have read and agreed to be bound by the Policies"
	// " and Procedures."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 0, termsAndConditions, "", 0, "L", true, 0, "")

	termsAndConditions = "    and Procedures."
	currY += 2.5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 0, termsAndConditions, "", 0, "L", true, 0, "")

	termsAndConditions = "4. I agree to adhere to the Know Your Customer ( KYC ) requirements as requested by Universe International Direct Selling"

	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 0, termsAndConditions, "", 0, "L", true, 0, "")

	termsAndConditions = "    (India) Pvt Ltd."
	currY += 2.5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 0, termsAndConditions, "", 0, "L", true, 0, "")

	currY += 10.0
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 0, "Received By: ", "", 0, "L", true, 0, "")
	currX += 20
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 7)
	pdf.CellFormat(0, 0, "                                     ", "", 0, "L", true, 0, "")

	currY += 10.0
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 0, "Received Date: ", "", 0, "L", true, 0, "")
	currX += 20
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 7)
	pdf.CellFormat(0, 0, "                                     ", "", 0, "L", true, 0, "")

	currY += 10.0
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 0, "Processed By: ", "", 0, "L", true, 0, "")
	currX += 20
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 7)
	pdf.CellFormat(0, 0, "                                     ", "", 0, "L", true, 0, "")

	currY += 10.0
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 0, "Processed Date: ", "", 0, "L", true, 0, "")
	currX += 20
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 7)
	pdf.CellFormat(0, 0, "                                     ", "", 0, "L", true, 0, "")

	currY += 10.0
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 7)
	pdf.CellFormat(0, 0, "                                                                                                                    ", "", 0, "L", true, 0, "")

	currX += 90
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 7)
	pdf.CellFormat(0, 0, "                                                                                                                    ", "", 0, "L", true, 0, "")

	currY += 5
	currX = 15
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(0, 0, "APPLICANT'S SIGNATURE (Required for processing)", "", 0, "L", true, 0, "")

	currX += 120
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(0, 0, "DATE", "", 0, "L", true, 0, "")

	// checkSpaceForTable(pdf, len(bvDistributionTable), rowHeight, lineHeight, headerHeight)
	pdf.AddPage()
	currY = 10.0
	currX = 10.0
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "DISTRIBUTOR APPLICATION FORM")
	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(128, 128, 128)

	txt := "1. I am at least 21 years old and eligible to enter a business contract in INDIA."
	currY += 10.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "2. I have undergone an orientation session (online or classroom) about the Company and Business Opportunity, and I do not have any doubts about either."
	currY += 6.0
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "3. By signing below, I certify and acknowledge that I have read, understood, and agreed to be bound by the Policies & Procedures of the company."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "4. I agree to submit self-attested copies of the PAN card, Address Proof, Signature Proof, and Identity proof within 30 days of registration under the Know Your Customer"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    (KYC) requirements, failing which the company has right to block access to the back office and commissions/incentives/rewards until such time the documents are "
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    submitted and accepted by the company."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "5. I understand that the company has a 30-Day Buy Back Policy (calculated from the date of invoice of the product). I hereby agree that I will always inform customers "
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    about this policy, with relevant conditions for the same."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "6. I will be solely responsible for complying with the Income tax, Service Tax, Professional Tax, GST, Sales Tax/VAT, Octroi, LBT and/or any other local taxes and levies for "
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    sales of products and services as may be applicable from time to time. Obtaining of all such licenses and registrations to run the business as a direct seller from time "
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    to time will be my responsibility. The company will not be responsible for any deviation and violation of any such legal/statutory requirement from my end."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "7. I hereby agree that TDS (Tax Deducted at Source) will be deducted by the company from the commissions payable to me. "
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "8. I hereby confirm that my referrer explained the Business Opportunity to me clearly. I also confirm that I have registered willingly, without any coercion or undue"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    pressure from anyone."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "9. I am aware that there is no fee for enrolling as a Distributor, and that there is no compulsion to purchase any products to enroll or otherwise in this business"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    oppurtunity."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "9. I am aware that there is no fee for enrolling as a Distributor, and that there is no compulsion to purchase any products to enroll or otherwise in this business"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "    official website, which may be revised according to the statutory requirements."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "11. I shall always commit and communicate as trained and authorised by the company and never mislead or give any unauthorised information nor any competitive"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      comparison/advantage of the products and services to the prospective customers/prospective Distributors."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "12. I shall not unfairly denigrate any other organisation, brand, or product, directly or indirectly."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "13. I shall not make any product claims on my own and would share the information only from the official publications authorised by the company. Also, I shall not"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      misrepresent the earning potential of this Business Opportunity and shall always transparently share documented facts published on the official website of the"
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      company."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "14. I acknowledge and confirm that I will be, always, an independent, self-employed entrepreneur in relation to this Agreement as a Distributor, and will not, at any time,"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      claim or represent to be an employee of the company, nor will I claim any benefits accruing to the employees of the company."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "15. I understand and accept that in case of my resignation or the termination of this Agreement due to non-compliance of any condition/s mentioned in the Distributor"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      Application form and Policies & Procedures, I will lose the right to claim commissions/incentives/rewards, access to the back office, and right to purchase"
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      products/sales tools, etc., and I will not be eligible to get any remuneration from the company."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "16. I hereby declare that I am not engaged and will not engage in any anti-national activities, directly or indirectly."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "17. I accept the right of the company to reject this application for registration at its absolute sole discretion, without specifying any reasons thereof."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "18. I am authorised to enter in an agreement with the company and sign the Application Form."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "19. I also confirm and agree to abide by such terms and conditions as modified or amended by the company from time to time."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "20. I understand that upgrade of systems and processes is an ongoing exercise and the company will continue to work on the same, and I hereby confirm that I have read"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      and understood the latest available information on the company's official website before signing the application form. I also understand that in case of any clarification"
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      on the same, I can write to admin@ui-network.com ."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "21. I fully understand that if I do not agree with the Policies & Procedures, I am free to resign by downloading, filling up, and sending the resignation form to the company."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "22. I undertake to confirm the accuracy of the information by referring to the company's official website before propagating the same."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "23. I hereby declare that the details furnished above are true and correct to the best of my knowledge and belief, and I undertake to inform you of any changes therein,"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      immediately."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "24. In case any of the above information is found to be false or untrue, I am aware that I may be held liable for it and the company may terminate this Agreement as soon"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      as the company confirms the same and take legal action against me to recover the commissions, bonuses, and rewards, etc. paid to me."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "25. I hereby acknowledge that I have reviewed and understood this Distributor Application and Agreement, along with all the relevant documents, including the Code of"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      Conduct, Policies & Procedures, and Compensation Plan defined herein as \"Materials\" which are incorporated herewith, and that I agree to be bound by all of them."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "26. I understand that this agreement will be terminated in case of no business done for a consecutive period of 24 months."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "27. I hereby agree to send in any complaints to the company on the official email address admin@ui-network.com and give the company 45 days to provide an"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      unbiased resolution according to the policies of the company before escalating it to the Department of Consumer Affairs or any other government body."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "28. I agree to receive updates from the company on my registered email address and registered mobile number."
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "29. I understand that the personal details I have shared with the company, along with data regarding the business, may be shared with the governmental agencies, if"
	currY += 6.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	txt = "      required."
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")

	currY += 10.0
	currX = 10
	pdf.SetXY(currX, currY)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "BU", 7)
	pdf.CellFormat(0, 0, "                                                                                                                    ", "", 0, "L", true, 0, "")

	currX += 90
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, "                                                                                                                    ", "", 0, "L", true, 0, "")

	currY += 5
	currX = 15
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(0, 0, "APPLICANT'S SIGNATURE (Required for processing)", "", 0, "L", true, 0, "")

	currX += 120
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(0, 0, "DATE", "", 0, "L", true, 0, "")

	currY += 5
	currX = 10

	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(128, 128, 128)

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

func GenerateDistributorForm(distribId string, tx *gorm.DB) (string, int, error) {

	distributorInformation, err := repositories.GetUserByID(distribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetUserByID repositories fn from GenerateDistributorForm service fn", err.Error())
		return err.Error(), fiber.StatusInternalServerError, err
	}

	referrerDistribId, err := repositories.GetRefDistribIdByDistribId(distribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetRefDistribIdByDistribId repositories fn from GenerateDistributorForm service fn", err.Error())
		return err.Error(), fiber.StatusInternalServerError, err
	}

	referrerDistribInformation, err := repositories.GetUserByID(referrerDistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetUserByID repositories fn from GenerateDistributorForm service fn", err.Error())
		return err.Error(), fiber.StatusInternalServerError, err
	}

	pdf := DistributorFormFactory(distributorInformation, referrerDistribInformation)

	filename := fmt.Sprintf("./assets/%s-distrib-form.pdf", distribId)
	err = pdf.OutputFileAndClose(filename)
	if err != nil {
		return err.Error(), fiber.StatusInternalServerError, err
	}

	return filename, fiber.StatusOK, err
}

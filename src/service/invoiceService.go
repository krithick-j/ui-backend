package service

import (
	"fmt"
	"ui-back-end/src/dto"
	"ui-back-end/src/middleware"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"math"

	"github.com/go-pdf/fpdf"
	"github.com/gofiber/fiber/v2"
)

// how many digit's groups to process
const groupsNumber int = 4

var _smallNumbers = []string{
	"Zero", "One", "Two", "Three", "Four",
	"Five", "Six", "Seven", "Eight", "Nine",
	"Ten", "Eleven", "Twelve", "Thirteen", "Fourteen",
	"Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen",
}
var _tens = []string{
	"", "", "Twenty", "Thirty", "Forty", "Fifty",
	"Sixty", "Seventy", "Eighty", "Ninety",
}
var _scaleNumbers = []string{
	"", "Thousand", "million", "billion",
}

type digitGroup int

// Convert converts number into the words representation.
func Convert(number int) string {
	return convert(number, false)
}

// ConvertAnd converts number into the words representation
// with " and " added between number groups.
func ConvertAnd(number int) string {
	return convert(number, true)
}

func intMod(x, y int) int {
	return int(math.Mod(float64(x), float64(y)))
}

func digitGroup2Text(group digitGroup, useAnd bool) (ret string) {
	hundreds := group / 100
	tensUnits := intMod(int(group), 100)

	if hundreds != 0 {
		ret += _smallNumbers[hundreds] + " Hundred"

		if tensUnits != 0 {
			ret += separator(useAnd)
		}
	}

	tens := tensUnits / 10
	units := intMod(tensUnits, 10)

	if tens >= 2 {
		ret += _tens[tens]

		if units != 0 {
			ret += "-" + _smallNumbers[units]
		}
	} else if tensUnits != 0 {
		ret += _smallNumbers[tensUnits]
	}

	return
}

// separator returns proper separator string between
// number groups.
func separator(useAnd bool) string {
	if useAnd {
		return " and "
	}
	return " "
}

func convert(number int, useAnd bool) string {
	// Zero rule
	if number == 0 {
		return _smallNumbers[0]
	}

	// Divide into three-digits group
	var groups [groupsNumber]digitGroup
	positive := math.Abs(float64(number))

	// Form three-digit groups
	for i := 0; i < groupsNumber; i++ {
		groups[i] = digitGroup(math.Mod(positive, 1000))
		positive /= 1000
	}

	var textGroup [groupsNumber]string
	for i := 0; i < groupsNumber; i++ {
		textGroup[i] = digitGroup2Text(groups[i], useAnd)
	}
	combined := textGroup[0]
	and := useAnd && (groups[0] > 0 && groups[0] < 100)

	for i := 1; i < groupsNumber; i++ {
		if groups[i] != 0 {
			prefix := textGroup[i] + " " + _scaleNumbers[i]

			if len(combined) != 0 {
				prefix += separator(and)
			}

			and = false

			combined = prefix + combined
		}
	}

	if number < 0 {
		combined = "minus " + combined
	}

	return combined
}

func InvoiceFactory(orderDetails models.OrdersHeader, iCouponsArr []dto.OrderedICouponOut, bvDistributionTable []dto.PlaceBv, refDistribId string) *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 1)
	pdf.AddPage()
	var (
		currX float64 = 60
		currY float64 = 8
	)

	// Function to check if there is enough space for the table
	checkSpaceForTable := func(pdf *fpdf.Fpdf, rowCount int, rowHeight float64, lineHeight float64, headerHeight float64) {
		_, pageHeight := pdf.GetPageSize()
		_, _, _, bottomMargin := pdf.GetMargins()
		availableHeight := pageHeight - pdf.GetY() - bottomMargin - 10
		requiredHeight := lineHeight + headerHeight + float64(int(rowHeight)*rowCount)
		fmt.Println("Required height--->", requiredHeight, "Available Height----->", availableHeight)
		if requiredHeight > availableHeight {
			pdf.AddPage()
			currY = 10.0
		}
	}

	pdf.Image("./assets/images/uilogo.png", currX, currY, 6, 0, false, "png", 0, "")
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
	currX += 80
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "INVOICE")

	//Icoupon Independent Representative Receipt
	currX = 10
	currY += 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 0, "ICoupon INDEPENDENT REPRESENTATIVE RECEIPT")
	pdf.SetTextColor(0, 0, 0)

	//Distributor ID No:
	currX = 10
	currY += 7
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Distributor ID No:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 0, orderDetails.DistribId)
	pdf.SetFont("Arial", "", 8)

	//Reference Distrib ID No
	currX += 60
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Reference Distributor ID No:")
	currX += 40
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, refDistribId)
	pdf.SetFont("Arial", "", 8)

	//Invoice No:
	currX = 10
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Invoice No:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.OrderId)
	pdf.SetFont("Arial", "", 8)

	//Order Date
	currX += 60
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Order Date:")
	currX += 40
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	formattedItcTime := middleware.FormatTimeByLocation(orderDetails.CreatedAt, "Asia/Kolkata", "02-01-2006")
	pdf.Cell(0, 0, formattedItcTime)
	pdf.SetFont("Arial", "", 8)
	//ORDER DETAILS END

	//Line
	currX = 0
	currY += 5
	pdf.Line(currX, currY, currX+600, currY)

	//ORDER DETAILS START

	currX = 10
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "CUSTOMER DETAILS")
	pdf.SetFont("Arial", "", 10)

	//Set Customer Address
	//To
	currX = 10
	currY += 7
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "To:")

	//Contact Name
	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Customer Name:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.ContactName)
	pdf.SetFont("Arial", "", 8)

	//Address
	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Address:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.Address)
	pdf.SetFont("Arial", "", 8)

	//City
	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "City:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.City)
	pdf.SetFont("Arial", "", 8)

	//City
	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "District:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.District)

	pdf.SetFont("Arial", "", 8)

	//ZipCode
	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Zip code:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.ZipCode)
	pdf.SetFont("Arial", "", 8)

	//State
	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "State:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.State)
	pdf.SetFont("Arial", "", 8)

	//Country
	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Country:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.Country)
	pdf.SetFont("Arial", "", 8)

	//Mobile Phone No
	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "City:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.MobilePhoneNo)
	pdf.SetFont("Arial", "", 8)

	//Home Phone No
	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "State:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.HomePhoneNo)
	pdf.SetFont("Arial", "", 8)

	rowHeight := 8.0
	lineHeight := 10.0
	headerHeight := 8.0

	checkSpaceForTable(pdf, len(orderDetails.OrdersLiner), rowHeight, lineHeight, headerHeight)
	//Line
	currX = 0
	currY += 5
	pdf.Line(currX, currY, currX+600, currY)

	//ORDER DETAILS START
	currX += 10
	currY += lineHeight
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "ORDER DETAILS")
	pdf.SetFont("Arial", "", 7)

	//Set Invoice Header
	pdf.SetFillColor(106, 163, 101)

	currY += 5
	pdf.SetXY(currX, currY)
	pdf.CellFormat(8, headerHeight, "SL No.", "1", 0, "C", true, 0, "")
	currX += 8
	pdf.SetXY(currX, currY)
	pdf.CellFormat(70, headerHeight, "Item", "1", 0, "C", true, 0, "")
	currX += 70
	pdf.SetXY(currX, currY)
	pdf.CellFormat(10, headerHeight, "Qty", "1", 0, "C", true, 0, "")
	currX += 10
	pdf.SetXY(currX, currY)
	pdf.MultiCell(20, 4, "Unit Price \nBefore GST", "1", "CT", true)
	currX += 20
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 7.5)
	pdf.CellFormat(12, headerHeight, "GST %", "1", 0, "R", true, 0, "")
	currX += 12
	pdf.SetXY(currX, currY)
	pdf.CellFormat(20, rowHeight, "GST Amount", "1", 0, "R", true, 0, "")
	currX += 20
	pdf.SetXY(currX, currY)
	pdf.MultiCell(18, 4, "Unit Price \nAfter GST", "1", "CT", true)
	currX += 18
	pdf.SetXY(currX, currY)
	pdf.CellFormat(14, headerHeight, "S & H", "1", 0, "C", true, 0, "")
	currX += 14
	pdf.SetXY(currX, currY)
	pdf.CellFormat(24, headerHeight, "Total (INR)", "1", 0, "C", true, 0, "")

	//Set Invoice Items
	pdf.SetFillColor(236, 253, 235)

	fmt.Println("order Details", orderDetails)
	for i, product := range orderDetails.OrdersLiner {
		fmt.Print("index:", i, "product details: ", product)
		//Set Invoice Header
		currX = 10
		currY += rowHeight
		pdf.SetXY(currX, currY)
		pdf.CellFormat(8, rowHeight, fmt.Sprintf("%v", i+1), "1", 0, "C", true, 0, "") //sl no
		currX += 8
		pdf.SetXY(currX, currY)
		pdf.CellFormat(70, rowHeight, product.Name, "1", 0, "C", true, 0, "")
		currX += 70
		pdf.SetXY(currX, currY)
		pdf.CellFormat(10, rowHeight, fmt.Sprintf("%v", product.Quantity), "1", 0, "R", true, 0, "")
		currX += 10
		pdf.SetXY(currX, currY)
		UnitPriceBeforeGST := product.UnitPrice / (1 + product.GstPercentage/100)
		pdf.CellFormat(20, rowHeight, fmt.Sprintf("%.2f", UnitPriceBeforeGST), "1", 0, "R", true, 0, "")
		currX += 20
		GSTAmount := product.UnitPrice - UnitPriceBeforeGST
		pdf.CellFormat(12, rowHeight, fmt.Sprintf("%v", product.GstPercentage), "1", 0, "R", true, 0, "")
		currX += 12
		pdf.SetXY(currX, currY)
		pdf.CellFormat(20, rowHeight, fmt.Sprintf("%.2f", GSTAmount), "1", 0, "R", true, 0, "")
		currX += 20
		pdf.SetXY(currX, currY)
		pdf.CellFormat(18, rowHeight, fmt.Sprintf("%v", product.UnitPrice), "1", 0, "R", true, 0, "")
		currX += 18
		pdf.SetXY(currX, currY)
		pdf.CellFormat(14, rowHeight, fmt.Sprintf("%v", product.SandH), "1", 0, "R", true, 0, "")
		currX += 14
		pdf.SetXY(currX, currY)
		pdf.CellFormat(24, rowHeight, fmt.Sprintf("%v", product.SubTotal), "1", 0, "R", true, 0, "")
	}
	//Total Invoice Details
	currX = 10
	currY += rowHeight
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 7.5)
	pdf.CellFormat(158, 8, "Total", "1", 0, "C", true, 0, "")
	pdf.SetFont("Arial", "", 7.5)
	currX += 158
	pdf.SetXY(currX, currY)
	pdf.CellFormat(14, rowHeight, fmt.Sprintf("%v", orderDetails.TotalSandH), "1", 0, "R", true, 0, "")
	currX += 14
	pdf.SetXY(currX, currY)
	pdf.CellFormat(24, rowHeight, fmt.Sprintf("%v", orderDetails.SubTotal), "1", 0, "R", true, 0, "")
	currX = 10
	currY += 10
	pdf.SetXY(currX, currY)
	pdf.SetFillColor(200, 200, 200)
	pdf.SetFont("Arial", "B", 7.5)
	pdf.CellFormat(172, 8, "Grand Total(Total Unit Price After GST + Total S and H)", "1", 0, "C", true, 0, "")
	pdf.SetFont("Arial", "", 7.5)
	currX += 172
	pdf.SetXY(currX, currY)
	pdf.CellFormat(24, 8, fmt.Sprintf("%v", orderDetails.TotalAmount), "1", 0, "R", true, 0, "")
	inwords := convert(int(orderDetails.TotalAmount), true)
	currX += 24
	currY += 15
	pdf.SetXY(currX, currY)

	//Amount in words
	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Amount In Words: ")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, fmt.Sprintf("%s Only", inwords))
	pdf.SetFont("Arial", "", 8)

	lineHeight = 5.0
	headerHeight = 8.0
	rowHeight = 8.0
	checkSpaceForTable(pdf, len(bvDistributionTable), rowHeight, lineHeight, headerHeight)
	//Line
	currX = 0
	currY += lineHeight
	pdf.Line(currX, currY, currX+600, currY)

	//BV DISTRIBUTION DETAILS START
	currX += 10
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "BV DISTRIBUTION DETAILS")
	pdf.SetFont("Arial", "", 7)

	//Set Tc Header
	pdf.SetFillColor(106, 163, 101)

	currY += 5
	pdf.SetXY(currX, currY)
	pdf.CellFormat(8, rowHeight, "TC", "1", 0, "C", true, 0, "")
	currX += 8
	pdf.SetXY(currX, currY)
	pdf.CellFormat(90, rowHeight, "Added BV", "1", 0, "C", true, 0, "")
	pdf.SetFillColor(236, 253, 235)

	for _, placeBvs := range bvDistributionTable {
		//Set TC items
		currX = 10
		currY += rowHeight
		pdf.SetXY(currX, currY)
		pdf.CellFormat(8, rowHeight, fmt.Sprintf("%v", placeBvs.Place), "1", 0, "C", true, 0, "") //sl no
		currX += rowHeight
		pdf.SetXY(currX, currY)
		pdf.CellFormat(90, rowHeight, fmt.Sprintf("%v", placeBvs.AddBv), "1", 0, "C", true, 0, "")
		currX += 30
		currX += 20
		pdf.SetXY(currX, currY)
		currX += 20
	}

	rowHeight = 8.0
	lineHeight = 10.0
	headerHeight = 5.0
	//Line
	currX = 0
	currY += 10

	// Check space for the table
	// checkSpaceForTable(pdf, len(iCouponsArr), rowHeight, lineHeight, headerHeight)
	pdf.AddPage()
	currY = 10.0

	// ICOUPON DETAILS START
	currX += 10
	currY += 6
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "ICOUPONS DETAILS")
	pdf.SetFont("Arial", "", 7)

	// Set Invoice Header
	pdf.SetFillColor(106, 163, 101)

	currY += headerHeight //5
	pdf.SetXY(currX, currY)
	pdf.CellFormat(8, rowHeight, "Sl No.", "1", 0, "C", true, 0, "")
	currX += 8
	pdf.SetXY(currX, currY)
	pdf.CellFormat(30, rowHeight, "ICoupon No", "1", 0, "C", true, 0, "")
	currX += 30
	pdf.SetXY(currX, currY)
	pdf.CellFormat(20, rowHeight, "V ID", "1", 0, "C", true, 0, "")
	currX += 20
	pdf.SetXY(currX, currY)
	pdf.CellFormat(20, rowHeight, "Total Value", "1", 0, "C", true, 0, "")
	currX += 20
	pdf.SetXY(currX, currY)
	pdf.CellFormat(20, rowHeight, "Used Value", "1", 0, "C", true, 0, "")

	// Set Invoice Items
	pdf.SetFillColor(236, 253, 235)

	for i, icoupon := range iCouponsArr {

		// Set Invoice Header
		currX = 10
		currY += rowHeight //8
		pdf.SetXY(currX, currY)
		pdf.CellFormat(8, rowHeight, fmt.Sprintf("%v", i+1), "1", 0, "C", true, 0, "") // sl no
		currX += 8
		pdf.SetXY(currX, currY)
		pdf.CellFormat(30, rowHeight, icoupon.VID, "1", 0, "C", true, 0, "")
		currX += 30
		pdf.SetXY(currX, currY)
		pdf.CellFormat(20, rowHeight, icoupon.VID, "1", 0, "C", true, 0, "")
		currX += 20
		pdf.SetXY(currX, currY)
		pdf.CellFormat(20, rowHeight, fmt.Sprintf("%.2f", icoupon.Value), "1", 0, "R", true, 0, "")
		currX += 20
		pdf.SetXY(currX, currY)
		pdf.CellFormat(20, rowHeight, fmt.Sprintf("%.2f", icoupon.UsedValue), "1", 0, "R", true, 0, "")
		currX += 20
	}

	//Line
	currX = 0
	currY += lineHeight //10

	currX += 10
	currY += 15
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "For GS Enterprises")
	currY += 20
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Godwin Selvan")
	currY += 8
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Digitally Generated Invoice. No Sign Required")
	return pdf
}

func GenerateInvoice(orderId string) (string, int, error) {

	orderDetails, err := repositories.GetOrderDetailsByOrderId(orderId)
	if orderDetails.DistribId == "" {
		return "Record not Found", fiber.StatusNotFound, err
	}

	iCouponsArr, err := GetICouponArrayByOrderId(orderId, orderDetails.DistribId)
	if err != nil {
		return err.Error(), fiber.StatusInternalServerError, err
	}

	// bvDistributionTable := []dto.PlaceBv{
	// 	{
	// 		Place: "001",
	// 		AddBv: 250,
	// 	},
	// 	{
	// 		Place: "002",
	// 		AddBv: 250,
	// 	},
	// 	{
	// 		Place: "003",
	// 		AddBv: 300,
	// 	},
	// }
	bvDistributionTable, err := GetBvDistributionTableByOrdeId(orderId, orderDetails.DistribId)
	if err != nil {
		return err.Error(), fiber.StatusInternalServerError, err
	}

	referrerDistribId, res := repositories.GetRefDistribIdByDistribId(orderDetails.DistribId)
	if res.Error != nil {
		return res.Error.Error(), fiber.StatusInternalServerError, res.Error
	}
	pdf := InvoiceFactory(orderDetails, iCouponsArr, bvDistributionTable, referrerDistribId)
	filename := fmt.Sprintf("./tmp/invoice-%s.pdf", orderId)
	err = pdf.OutputFileAndClose(filename)
	if err != nil {
		return err.Error(), fiber.StatusInternalServerError, err
	}
	return filename, fiber.StatusOK, err
}

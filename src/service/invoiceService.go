package service

import (
	"fmt"
	"ui-back-end/src/middleware"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"math"

	"github.com/go-pdf/fpdf"
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

func InvoiceFactory(orderDetails models.OrdersHeader) *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	var (
		currX float64 = 60
		currY float64 = 8
	)
	pdf.Image("./assets/images/uilogo.png", currX, currY, 6, 0, false, "png", 0, "")

	pdf.SetFont("Arial", "B", 12)
	currX += 8
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Universe International")
	pdf.SetFont("Arial", "", 8)
	currX -= 30
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "GS ENTERPRISES FLAT NO D, KASI ARCADE FIRST FLOOR, VOC STREET")
	currX += 4
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "KAIKANKUPPAM, RAMAPURAM, CHENNAI. TAMIL NADU PIN 600087.")
	currX += 12
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Ph: 9080296128 Email: selvangodson@gmail.com")
	currX = 0
	currY += 5
	pdf.Line(currX, currY, currX+600, currY)
	pdf.SetFont("Arial", "BU", 12)
	currX += 80
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "INVOICE")
	pdf.SetFont("Arial", "", 8)

	//Set Customer Address
	currX = 10
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "To:")
	currX += 3
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, orderDetails.ContactName)
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, orderDetails.Address)
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, orderDetails.City)
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, orderDetails.State)
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, orderDetails.ZipCode)
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, orderDetails.Country)
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, orderDetails.MobilePhoneNo)

	//Set Invoice Details
	currX += 120
	currY -= 16
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, fmt.Sprintf("Invoice No: %s", orderDetails.OrderId))
	currY += 4
	pdf.SetXY(currX, currY)
	formattedItcTime := fmt.Sprintf("Invoice Date: %s", middleware.FormatTimeByLocation(orderDetails.CreatedAt, "Asia/Kolkata", "02-01-2006"))
	pdf.Cell(0, 0, formattedItcTime)

	//Set Invoice Header
	pdf.SetFillColor(200, 200, 200)
	currX = 10
	currY += 20
	pdf.SetXY(currX, currY)
	pdf.CellFormat(10, 8, "Sl No.", "1", 0, "C", true, 0, "")
	currX += 10
	pdf.SetXY(currX, currY)
	pdf.CellFormat(100, 8, "Item", "1", 0, "C", true, 0, "")
	currX += 100
	pdf.SetXY(currX, currY)
	pdf.CellFormat(10, 8, "Qty", "1", 0, "C", true, 0, "")
	currX += 10
	pdf.SetXY(currX, currY)
	pdf.CellFormat(20, 8, "Unit Price", "1", 0, "C", true, 0, "")
	currX += 20
	pdf.SetXY(currX, currY)
	pdf.CellFormat(16, 8, "GST %", "1", 0, "C", true, 0, "")
	currX += 16
	pdf.SetXY(currX, currY)
	pdf.CellFormat(16, 8, "S & H", "1", 0, "C", true, 0, "")
	currX += 16
	pdf.SetXY(currX, currY)
	pdf.CellFormat(20, 8, "Total (INR)", "1", 0, "C", true, 0, "")

	//Set Invoice Items
	pdf.SetFillColor(255, 255, 255)
	fmt.Println("order Details", orderDetails)
	for i, product := range orderDetails.OrdersLiner {
		fmt.Print("index:", i, "product details: ", product)
		//Set Invoice Header
		currX = 10
		currY += 8
		pdf.SetXY(currX, currY)
		pdf.CellFormat(10, 8, fmt.Sprintf("%v", i+1), "1", 0, "C", true, 0, "") //sl no
		currX += 10
		pdf.SetXY(currX, currY)
		pdf.CellFormat(100, 8, product.Name, "1", 0, "C", true, 0, "")
		currX += 100
		pdf.SetXY(currX, currY)
		pdf.CellFormat(10, 8, fmt.Sprintf("%v", product.Quantity), "1", 0, "C", true, 0, "")
		currX += 10
		pdf.SetXY(currX, currY)
		pdf.CellFormat(20, 8, fmt.Sprintf("%v", product.UnitPrice), "1", 0, "R", true, 0, "")
		currX += 20
		pdf.SetXY(currX, currY)
		pdf.CellFormat(16, 8, fmt.Sprintf("%v", product.GstPercentage), "1", 0, "C", true, 0, "")
		currX += 16
		pdf.SetXY(currX, currY)
		pdf.CellFormat(16, 8, fmt.Sprintf("%v", product.SandH), "1", 0, "R", true, 0, "")
		currX += 16
		pdf.SetXY(currX, currY)
		pdf.CellFormat(20, 8, fmt.Sprintf("%v", product.SubTotal), "1", 0, "R", true, 0, "")
	}
	//Total Invoice Details
	currX = 10
	currY += 8
	pdf.SetXY(currX, currY)
	pdf.CellFormat(156, 8, "Total", "1", 0, "C", true, 0, "")
	currX += 156
	pdf.SetXY(currX, currY)
	pdf.CellFormat(16, 8, fmt.Sprintf("%v", orderDetails.TotalSandH), "1", 0, "R", true, 0, "")
	currX += 16
	pdf.SetXY(currX, currY)
	pdf.CellFormat(20, 8, fmt.Sprintf("%v", orderDetails.SubTotal), "1", 0, "R", true, 0, "")
	currX = 10
	currY += 8
	pdf.SetXY(currX, currY)
	pdf.CellFormat(172, 8, "Grand Total", "1", 0, "C", true, 0, "")
	currX += 172
	pdf.SetXY(currX, currY)
	pdf.CellFormat(20, 8, fmt.Sprintf("%v", orderDetails.TotalAmount), "1", 0, "R", true, 0, "")
	inwords := convert(int(orderDetails.TotalAmount), true)
	currX = 10
	currY += 15
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, fmt.Sprintf("Amount In Words: %v Only", inwords))
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

func GenerateInvoice(orderId string) (string, error) {

	orderDetails, err := repositories.GetOrderDetailsByOrderId(orderId)
	if err != nil {
		return "", err
	}
	// Create ID Card Layout
	pdf := InvoiceFactory(orderDetails)
	//pdf = IDCardAddContent(pdf)
	filename := fmt.Sprintf("./tmp/invoice-%s.pdf", orderId)
	err = pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}
	return filename, err
}

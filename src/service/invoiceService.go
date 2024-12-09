package service

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"math"

	"github.com/go-pdf/fpdf"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

func GenerateInvoice(orderId string, productType string, tx *gorm.DB) (string, int, error) {

	orderDetails, err := repositories.GetOrderDetailsByOrderId(orderId, tx)
	if err == gorm.ErrRecordNotFound {
		configs.Log.
			Infoln("Record not Found on calling the GetOrderDetailsByOrderId fn from GenerateInvoice service fn", err.Error())
		return "Record not Found", fiber.StatusNotFound, err
	}

	res, status := GetICouponArrayTotalValueByOrderId(orderId, orderDetails.DistribId, tx)
	if status != fiber.StatusOK {
		tx.Rollback()
		configs.Log.Errorln("Error on calling the GetICouponArrayTotalValueByOrderId fn from GenerateInvoice service fn", res["error"])
		return res["error"].(string), status, res["err"].(error)
	}
	iCouponsArr := res["data"].([]dto.OrderedICouponOut)

	bvDistributionTable, err := GetBvDistributionTableByOrdeId(orderId, orderDetails.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling the GetBvDistributionTableByOrdeId fn from GenerateInvoice service fn", res["error"])
		return err.Error(), fiber.StatusInternalServerError, err
	}

	referrerDistribId, err := repositories.GetRefDistribIdByDistribId(orderDetails.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.
			Errorln("Error on calling the GetRefDistribIdByDistribId fn from GenerateInvoice service fn", res["error"])
		return err.Error(), fiber.StatusInternalServerError, err
	}

	var pdf *fpdf.Fpdf
	if productType == "bv" {
		bvDistributionTable = append(bvDistributionTable, dto.PlaceBv{
			Place: "",
			AddBv: 0,
		})
		bvPdf := BvInvoiceFactory(orderDetails, iCouponsArr, bvDistributionTable, referrerDistribId)
		pdf = bvPdf
	} else if productType == "rsp" {
		rspPdf := RspInvoiceFactory(orderDetails, iCouponsArr, referrerDistribId)
		pdf = rspPdf
	} else if productType == "ep" {
		EpInvoiceFactory(orderDetails, referrerDistribId)
	} else {
		return "", fiber.StatusNotFound, fmt.Errorf("INVALID PRODUCT TYPE")
	}

	filename := fmt.Sprintf("./tmp/invoice-%s.pdf", orderId)
	err = pdf.OutputFileAndClose(filename)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln(fmt.Sprintf("Error on calling %s from %s service fn: %s", "OutputFileAndClose", "GenerateInvoice", err.Error()))
		return err.Error(), fiber.StatusInternalServerError, err
	}
	return filename, fiber.StatusOK, err
}

func GeneralDetails(currY float64, pdf *fpdf.Fpdf, currX float64, orderDetails models.OrdersHeader, refDistribId string) {
	currY += 25
	pdf.SetXY(currX, currY)
	pdf.Line(currX, currY, currX+600, currY)

	pdf.SetFont("Arial", "BU", 13)
	currX += 80
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "INVOICE")

	currX = 10
	currY += 10
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 0, "ICoupon INDEPENDENT REPRESENTATIVE RECEIPT")
	pdf.SetTextColor(0, 0, 0)

	currX = 10
	currY += 7
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Distributor ID No:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 0, orderDetails.DistribId)
	pdf.SetFont("Arial", "", 8)

	currX += 60
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Reference Distributor ID No:")
	currX += 40
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, refDistribId)
	pdf.SetFont("Arial", "", 8)

	currX = 10
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Invoice No:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.OrderId)
	pdf.SetFont("Arial", "", 8)

	currX += 60
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Order Date:")
	currX += 40
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	formattedItcTime := utils.FormatTimeByLocation(orderDetails.CreatedAt, "Asia/Kolkata", "02-01-2006 15:04:05")
	pdf.Cell(0, 0, formattedItcTime)
	pdf.SetFont("Arial", "", 8)
}

func CustomerDetails(currY float64, pdf *fpdf.Fpdf, currX float64, orderDetails models.OrdersHeader) {
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "CUSTOMER DETAILS")
	pdf.SetFont("Arial", "", 10)

	currX = 10
	currY += 7
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "To:")

	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Customer Name:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.ContactName)
	pdf.SetFont("Arial", "", 8)

	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Address:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.Address)
	pdf.SetFont("Arial", "", 8)

	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "City:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.City)
	pdf.SetFont("Arial", "", 8)

	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "District:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.District)

	pdf.SetFont("Arial", "", 8)

	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Zip code:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.ZipCode)
	pdf.SetFont("Arial", "", 8)

	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "State:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.State)
	pdf.SetFont("Arial", "", 8)

	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Country:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.Country)
	pdf.SetFont("Arial", "", 8)

	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Mobile Phone No:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.MobilePhoneNo)
	pdf.SetFont("Arial", "", 8)

	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Home Phone No:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.HomePhoneNo)
	pdf.SetFont("Arial", "", 8)

	currX = 10
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Email Address:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.ContactEmail)
	pdf.SetFont("Arial", "", 8)
}

func OrderDetailsTable(currX float64, currY float64, lineHeight float64, pdf *fpdf.Fpdf, headerHeight float64, rowHeight float64, orderDetails models.OrdersHeader, checkSpaceForTable func(pdf *fpdf.Fpdf, rowHeight float64, lineHeight float64, headerHeight float64)) string {
	currX += 10
	currY += lineHeight
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "ORDER DETAILS")
	pdf.SetFont("Arial", "", 7)

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

	pdf.SetFillColor(236, 253, 235)

	fmt.Println("order Details", orderDetails)
	for i, product := range orderDetails.OrdersLiners {
		checkSpaceForTable(pdf, rowHeight, lineHeight, headerHeight)

		fmt.Print("index:", i, "product details: ", product)

		currX = 10
		currY += rowHeight
		pdf.SetXY(currX, currY)
		pdf.CellFormat(8, rowHeight, fmt.Sprintf("%v", i+1), "1", 0, "C", true, 0, "")
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
	return inwords
}

func BvDistributionTable(currX float64, currY float64, pdf *fpdf.Fpdf, rowHeight float64, bvDistributionTable []dto.PlaceBv, checkSpaceForTable func(pdf *fpdf.Fpdf, rowHeight float64, lineHeight float64, headerHeight float64), lineHeight float64, headerHeight float64) {
	currX += 10
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "BV DISTRIBUTION DETAILS")
	pdf.SetFont("Arial", "", 7)

	pdf.SetFillColor(106, 163, 101)

	currY += 5
	pdf.SetXY(currX, currY)
	pdf.CellFormat(8, rowHeight, "TC", "1", 0, "C", true, 0, "")
	currX += 8
	pdf.SetXY(currX, currY)
	pdf.CellFormat(90, rowHeight, "Added BV", "1", 0, "C", true, 0, "")
	pdf.SetFillColor(236, 253, 235)

	totalBv := 0.0
	for _, placeBvs := range bvDistributionTable {

		currX = 10
		currY += rowHeight
		pdf.SetXY(currX, currY)
		pdf.CellFormat(8, rowHeight, fmt.Sprintf("%v", placeBvs.Place), "1", 0, "C", true, 0, "")
		currX += rowHeight
		pdf.SetXY(currX, currY)
		pdf.CellFormat(90, rowHeight, fmt.Sprintf("%.2f", placeBvs.AddBv), "1", 0, "R", true, 0, "")
		currX += 30
		currX += 20
		pdf.SetXY(currX, currY)
		currX += 20
		totalBv += placeBvs.AddBv
		checkSpaceForTable(pdf, rowHeight, lineHeight, headerHeight)
	}

	currX = 10
	currY += 10
	pdf.SetXY(currX, currY)
	pdf.SetFillColor(200, 200, 200)
	pdf.SetFont("Arial", "B", 7.5)
	pdf.CellFormat(54, rowHeight, "Total Distribution Points", "1", 0, "C", true, 0, "")
	pdf.SetFont("Arial", "", 7.5)
	currX += 54
	pdf.SetXY(currX, currY)
	pdf.CellFormat(44, rowHeight, fmt.Sprintf("%.2f", totalBv), "1", 0, "R", true, 0, "")
}

func ICouponTable(currX float64, currY float64, pdf *fpdf.Fpdf, headerHeight float64, rowHeight float64, iCouponsArr []dto.OrderedICouponOut, checkSpaceForTable func(pdf *fpdf.Fpdf, rowHeight float64, lineHeight float64, headerHeight float64), lineHeight float64) {
	currX += 10
	currY += 6
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "ICOUPON PAYMENT DETAILS")
	pdf.SetFont("Arial", "", 7)

	pdf.SetFillColor(106, 163, 101)

	currY += headerHeight
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

	pdf.SetFillColor(236, 253, 235)

	totalUsedValue := 0.0
	for i, icoupon := range iCouponsArr {
		checkSpaceForTable(pdf, rowHeight, lineHeight, headerHeight)

		currX = 10
		currY += rowHeight
		pdf.SetXY(currX, currY)
		pdf.CellFormat(8, rowHeight, fmt.Sprintf("%v", i+1), "1", 0, "C", true, 0, "")
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
		totalUsedValue += icoupon.UsedValue
	}

	currX = 10
	currY += 10
	pdf.SetXY(currX, currY)
	pdf.SetFillColor(200, 200, 200)
	pdf.SetFont("Arial", "B", 7.5)
	pdf.CellFormat(54, rowHeight, "Total ICoupon Used Value", "1", 0, "C", true, 0, "")
	pdf.SetFont("Arial", "", 7.5)
	currX += 54
	pdf.SetXY(currX, currY)
	pdf.CellFormat(44, rowHeight, fmt.Sprintf("%.2f", totalUsedValue), "1", 0, "R", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
}

func TermsAndConditionsPage(pdf *fpdf.Fpdf) {
	pdf.AddPage()
	currY := 10.0
	currX := 10.0
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "TERMS AND CONDITIONS")
	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(128, 128, 128)

	txt := "1. I have read, understood and agreed to be bound by all the terms and conditions set forth by Ubiquitous Infinity Network Pvt Ltd regarding this transaction."
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
	currY -= 3.5
	pdf.SetXY(currX, currY)
	pdf.SetTextColor(0, 0, 255)
	pdf.SetFont("Arial", "U", 7)

	pdf.CellFormat(0, 7, mail, "", 0, "L", false, 0, link)
	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(128, 128, 128)
	currX += 25.8
	currY += 3.5
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, rest, "", 0, "L", true, 0, "")

	txt = "Please PRINT this receipt for your future reference. For questions and comments, please eMail: "
	currX = 10
	currY += 3.0
	pdf.SetXY(currX, currY)
	pdf.CellFormat(0, 0, txt, "", 0, "L", true, 0, "")
	mailtxt := "admin@ui-network.com"
	currX += 105
	currY -= 3.5
	pdf.SetXY(currX, currY)
	pdf.SetTextColor(0, 0, 255)
	pdf.SetFont("Arial", "U", 7)
	pdf.CellFormat(0, 7, mailtxt, "", 0, "L", false, 0, link)
	pdf.SetFont("Arial", "", 7)
}

func Header(currX float64, currY float64, pdf *fpdf.Fpdf) float64 {
	currX += 12
	currY += 2
	pdf.SetXY(currX, currY)
	pdf.Image("assets/images/uilogo.jpeg", currX, currY, 47, 18, false, "jpeg", 0, "")
	return currY
}
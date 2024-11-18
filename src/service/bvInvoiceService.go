package service

import (
	"fmt"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/utils"

	"github.com/go-pdf/fpdf"
)

func BvInvoiceFactory(orderDetails models.OrdersHeader, iCouponsArr []dto.OrderedICouponOut, bvDistributionTable []dto.PlaceBv, refDistribId string) *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 1)
	pdf.AddPage()
	var (
		currX float64 = 60
		currY float64 = 8
	)
	// Function to check if there is enough space for the table
	checkSpaceForTable := func(pdf *fpdf.Fpdf, rowHeight float64, lineHeight float64, headerHeight float64) {
		_, pageHeight := pdf.GetPageSize()
		_, _, _, bottomMargin := pdf.GetMargins()
		availableHeight := pageHeight - pdf.GetY() - bottomMargin - 10
		requiredHeight := lineHeight + headerHeight + float64(int(rowHeight))
		if requiredHeight > availableHeight {
			pdf.AddPage()
			currY = 10.0
		}
	}
	currX += 12
	currY += 2
	pdf.SetXY(currX, currY)
	pdf.Image("assets/images/uilogo.jpeg", currX, currY, 47, 18, false, "jpeg", 0, "")

	//Line
	currX = 0
	currY += 25
	pdf.SetXY(currX, currY)
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
	formattedItcTime := utils.FormatTimeByLocation(orderDetails.CreatedAt, "Asia/Kolkata", "02-01-2006 15:04:05")
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
	pdf.Cell(0, 0, "Mobile Phone No:")
	currX += 35
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 0, orderDetails.MobilePhoneNo)
	pdf.SetFont("Arial", "", 8)

	//Home Phone No
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

	rowHeight := 8.0
	lineHeight := 10.0
	headerHeight := 8.0

	checkSpaceForTable(pdf, rowHeight, lineHeight, headerHeight)

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
	for i, product := range orderDetails.OrdersLiners {
		checkSpaceForTable(pdf, rowHeight, lineHeight, headerHeight)

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

	//Line
	currX = 0
	currY += lineHeight

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

	totalBv := 0.0
	for _, placeBvs := range bvDistributionTable {
		//Set TC items
		currX = 10
		currY += rowHeight
		pdf.SetXY(currX, currY)
		pdf.CellFormat(8, rowHeight, fmt.Sprintf("%v", placeBvs.Place), "1", 0, "C", true, 0, "") //sl no
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

	rowHeight = 8.0
	lineHeight = 10.0
	headerHeight = 5.0
	//Line
	currX = 0
	currY += 10

	// Check space for the table
	checkSpaceForTable(pdf, rowHeight, lineHeight, headerHeight)
	// pdf.AddPage()
	currY = 10.0

	// ICOUPON DETAILS START
	currX += 10
	currY += 6
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(0, 0, "ICOUPON PAYMENT DETAILS")
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

	totalUsedValue := 0.0
	for i, icoupon := range iCouponsArr {
		checkSpaceForTable(pdf, rowHeight, lineHeight, headerHeight)

		// Set Invoice Items
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

	TermsAndConditionsPage(pdf)

	return pdf
}

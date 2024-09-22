package service

import (
	"fmt"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/utils"

	"github.com/go-pdf/fpdf"
)

func RspInvoiceFactory(orderDetails models.OrdersHeader, iCouponsArr []dto.OrderedICouponOut, refDistribId string) *fpdf.Fpdf {
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

	checkSpaceForTable(pdf, len(orderDetails.OrdersLiners), rowHeight, lineHeight, headerHeight)
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
		checkSpaceForTable(pdf, len(orderDetails.OrdersLiners), rowHeight, lineHeight, headerHeight)

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

	pdf.AddPage()
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
	currX = 10
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
		checkSpaceForTable(pdf, len(iCouponsArr), rowHeight, lineHeight, headerHeight)

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

	pdf.AddPage()
	currY = 10.0
	currX = 10.0
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

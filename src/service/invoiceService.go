package service

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/repositories"

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

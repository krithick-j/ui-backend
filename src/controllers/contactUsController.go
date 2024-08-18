package controllers

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
)

func AddEnquiryType(c *fiber.Ctx) error {
	var request dto.EnquiryTypeIn
	if err := c.BodyParser(&request); err != nil {
		configs.Log.Errorln("Error on parsing request from AddEnquiryType controllers function", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}
	tx := configs.DB.Begin()
	res, status := service.AddEnquiryType(request, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func GetAllEnquiryType(c *fiber.Ctx) error {

	tx := configs.DB.Begin()
	res, status := service.GetAllEnquiryType(tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func SubmitContactUsQuery(c *fiber.Ctx) error {

	var request dto.ContactUsIn
	//for parsing textfield in form data
	if err := c.BodyParser(&request); err != nil {
		configs.Log.Errorln("Error on parsing request from SubmitContactUsQuery controllers fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid text format", "err": err.Error()})
	}

	if form, err := c.MultipartForm(); err != nil {
		configs.Log.Errorln("Error on calling MultipartForm controller fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form format", "err": err.Error()})
	} else {
		if fileFormObj, err := utils.ParseForm(form.File); err != nil {
			configs.Log.Errorln("Error on calling utils.ParseForm controller function", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot Parse input Form", "err": err.Error()})
		} else {
			tx := configs.DB.Begin()
			res, formName, status := service.SubmitContactUsQuery(*fileFormObj, request, tx)
			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
			}
			if status == 200 {

				if formName.AadhaarFront != "" {
					path := fmt.Sprintf("contactQueryUploads/%s", formName.AadhaarFront)
					c.SaveFile(&fileFormObj.AadhaarFront, path)
				}
				if formName.AadhaarBack != "" {
					path := fmt.Sprintf("contactQueryUploads/%s", formName.AadhaarBack)
					c.SaveFile(&fileFormObj.AadhaarBack, path)
				}
				if formName.PanCard != "" {
					path := fmt.Sprintf("contactQueryUploads/%s", formName.PanCard)
					c.SaveFile(&fileFormObj.PanCard, path)
				}
				if formName.PassportSize != "" {
					path := fmt.Sprintf("contactQueryUploads/%s", formName.PassportSize)
					c.SaveFile(&fileFormObj.PassportSize, path)
				}
			}
			return c.Status(status).JSON(res["success"])
		}
	}
}

func GetAllContactQueries(c *fiber.Ctx) error {

	tx := configs.DB.Begin()
	res, status := service.GetAllContactQueries(tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func SwitchContactQueryStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	tx := configs.DB.Begin()
	res, status := service.SwitchContactQueryStatus(id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

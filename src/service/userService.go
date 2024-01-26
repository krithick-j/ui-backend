package service

import (
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RecursiveUser struct {
	Name  string         `json:"name"`
	Link  string         `json:"link"`
	Left  *RecursiveUser `json:"left"`
	Right *RecursiveUser `json:"right"`
}

func GetUserByDistId(dist_id string) fiber.Map {

	var user models.User
	var result *gorm.DB

	user, result = repositories.GetUserByID(dist_id, user)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": user}
}

func GetUsers() fiber.Map {
	var users []models.User
	var result *gorm.DB
	users, result = repositories.GetAllUsers(users)
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": users}
}

func FindRecursiveFind(ruser *RecursiveUser, dist_id string, side string) {
	var user models.User
	var result *gorm.DB
	user, result = repositories.GetUserByID(dist_id, user)
	var nuser *RecursiveUser = new(RecursiveUser)
	nuser.Name = user.Name
	nuser.Link = "Link"
	if result.Error == gorm.ErrRecordNotFound {
		panic("Not Found")
	}
	if result.Error != nil {
		panic("Some Error")
	}
	if user.Lside != "" {
		FindRecursiveFind(nuser, user.Lside, "left")
	}
	if user.Rside != "" {
		FindRecursiveFind(nuser, user.Rside, "right")
	}
	if side == "left" {
		ruser.Left = nuser
	} else {
		ruser.Right = nuser
	}

}

func GetTreeUserByDistId(dist_id string) fiber.Map {

	var user models.User
	var result *gorm.DB
	data := new(RecursiveUser)
	user, result = repositories.GetUserByID(dist_id, user)
	data.Name = user.Name
	data.Link = "Link"
	if user.Lside != "" {
		FindRecursiveFind(data, user.Lside, "left")
	}
	if user.Rside != "" {
		FindRecursiveFind(data, user.Rside, "right")
	}
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": data}
}

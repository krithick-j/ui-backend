package controllers

// func TestController(c *fiber.Ctx) error {

// 	var payload dto.AddTc

// 	if err := c.BodyParser(&payload); err != nil {
// 		configs.Log.Errorln("Error on parsing payload from GetBvCounterByStartDate controllers fn", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
// 	}

// 	status, _ := service.AddTc(payload)
// 	return c.Status(status).JSON(res)
// }

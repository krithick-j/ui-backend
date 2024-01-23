package controllers


// func CreateBuyer(c *fiber.Ctx) error {
// 	var payload *models.User

// 	// Incoming Request Body
// 	if err := c.BodyParser(&payload); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
// 	}
// 	// hashedPassword, err1 := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
// 	// if err1 != nil {
// 	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err1.Error()})
// 	// }

// 	buyer := models.Buyer{
// 		Username: payload.Username,
// 		PhoneNo:  payload.PhoneNo,
// 		Email:    strings.ToLower(payload.Email),
// 		Password: string(hashedPassword),
// 		Image:    payload.Image,
// 		Age:      payload.Age,
// 	}

// 	// Check if user with the same email already exists
// 	if result := configs.DB.Where("email = ?", buyer.Email).First(&buyer); result.RowsAffected > 0 {
// 		// print("----------------_>%v", result.Error)/
// 		log.Info("User with email already exists.")
// 		return c.Status(http.StatusConflict).JSON(fiber.Map{"message": "User with email already exists..."})
// 	}

// 	// Check if user with the same phone number already exists
// 	if result := configs.DB.Where("phone_no = ?", buyer.PhoneNo).Take(&buyer); result.RowsAffected > 0 {
// 		log.Info("User with phone number already exists.")
// 		return c.Status(http.StatusConflict).JSON(fiber.Map{"message": "User with Phone number already exists"})
// 	}
// 	// Save the user to the database
// 	if err := repositories.BuyerSave(&buyer); err != nil {
// 		log.Info("Error saving user to the database:", err)
// 		return c.Status(http.StatusBadGateway).JSON(err)
// 	}

// 	log.Info("User created successfully.")

// 	return c.Status(http.StatusCreated).JSON(buyer)
// }
package controllers

import (
    "github.com/gofiber/fiber/v2"
    //"gorm.io/gorm"
    "hostel/database"
    "hostel/models"
)

func StudentSignIn(c *fiber.Ctx) error {
	var input models.Student

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Check for existing student by email or roll number
	var existing models.Student
	if err := database.DB.
		Where("email = ? OR roll_number = ?", input.Email, input.RollNumber).
		First(&existing).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Student with this email or roll number already exists"})
	}

	// Find a room with less than 3 students
	var rooms []models.Room
	if err := database.DB.Find(&rooms).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch rooms"})
	}

	var assignedRoom *models.Room
	for _, room := range rooms {
		if room.Occupied < 3 {
			assignedRoom = &room
			break
		}
	}

	if assignedRoom == nil {
		return c.Status(400).JSON(fiber.Map{"error": "No available rooms with free beds"})
	}

	// Assign room and create student
	student := models.Student{
		Name:       input.Name,
		Email:      input.Email,
		Password:   input.Password,
		RollNumber: input.RollNumber,
		Phone:      input.Phone,
		RoomID:     assignedRoom.ID,
	}

	if err := database.DB.Create(&student).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create student"})
	}

	// Increment the room's Occupied count
	assignedRoom.Occupied++
	database.DB.Save(&assignedRoom)
	if err := database.DB.Save(assignedRoom).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update room occupancy"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":      "Student signed in successfully",
		"assignedRoom": assignedRoom.RoomNumber,
	})
}



func StudentLogin(c *fiber.Ctx) error {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    var student models.Student
    if err := database.DB.Where("email = ?", input.Email).First(&student).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Student not found"})
    }

    if student.Password != input.Password {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
    }

    return c.JSON(student)
}

func AdminLogin(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var admin models.Admin
	if err := database.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Admin not found"})
	}

	if admin.Password != input.Password {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
	}

	// Return limited info if needed or full object
	return c.Status(200).JSON(fiber.Map{
		"message":  "Login successful",
		"admin_id": admin.ID,
		"username": admin.Username,
	})
}


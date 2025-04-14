package controllers

import (
	"hostel/database"
	"hostel/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var student models.Student
	if err := database.DB.Preload("Room").First(&student, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Student not found"})
	}

	var payments []models.Payment
	database.DB.Where("student_id = ?", id).Find(&payments)

	var checkInOuts []models.CheckInOut
	database.DB.Where("student_id = ?", id).Find(&checkInOuts)

	// Format check-in/out records
	var checkInOutFormatted []fiber.Map
for _, record := range checkInOuts {
	var checkOutStr string
	if record.CheckOut != nil {
		checkOutStr = record.CheckOut.Format("2006-01-02 15:04:05")
	} else {
		checkOutStr = "Not checked out"
	}

	checkInOutFormatted = append(checkInOutFormatted, fiber.Map{
		"check_in":  record.CheckIn.Format("2006-01-02 15:04:05"),
		"check_out": checkOutStr,
	})
}


	return c.JSON(fiber.Map{
		"student":           student,
		"payments":          payments,
		"check_in_out_logs": checkInOutFormatted,
	})
}

func CheckInStudent(c *fiber.Ctx) error {
	type CheckInInput struct {
		StudentID uint      `json:"student_id"`
		CheckIn   time.Time `json:"check_in"`
	}

	var input CheckInInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var existing models.CheckInOut
	database.DB.Where("student_id = ? AND check_out IS NULL", input.StudentID).First(&existing)
	if existing.ID != 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Already checked in"})
	}

	checkIn := models.CheckInOut{
        StudentID: input.StudentID,
        CheckIn:   input.CheckIn,
        CheckOut:  nil,
    }
    
    if checkIn.CheckIn.IsZero() {
        checkIn.CheckIn = time.Now()
    }
    

	if checkIn.CheckIn.IsZero() {
		checkIn.CheckIn = time.Now()
	}

	if err := database.DB.Create(&checkIn).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check in"})
	}

	return c.JSON(checkIn)
}

func CheckOutStudent(c *fiber.Ctx) error {
	type CheckOutInput struct {
		StudentID uint      `json:"student_id"`
		CheckOut  time.Time `json:"check_out"`
	}

	var input CheckOutInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var record models.CheckInOut
	if err := database.DB.Where("student_id = ? AND check_out IS NULL", input.StudentID).Order("id desc").First(&record).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "No active check-in found"})
	}

	if input.CheckOut.IsZero() {
        now := time.Now()
        record.CheckOut = &now
    } else {
        record.CheckOut = &input.CheckOut
    }
    

	if err := database.DB.Save(&record).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check out"})
	}

	return c.JSON(record)
}

package controllers

import (
	"hostel/models"
	"github.com/gofiber/fiber/v2"
	//"gorm.io/gorm"
	 "hostel/database"
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
        checkInOutFormatted = append(checkInOutFormatted, fiber.Map{
            "check_in":  record.CheckIn.Format("2006-01-02 15:04:05"),
            "check_out": record.CheckOut.Format("2006-01-02 15:04:05"),
        })
    }

    return c.JSON(fiber.Map{
        "student":           student,
        "payments":          payments,
        "check_in_out_logs": checkInOutFormatted,
    })
}

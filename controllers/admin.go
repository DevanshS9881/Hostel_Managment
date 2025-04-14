package controllers

import (
    "github.com/gofiber/fiber/v2"
    "hostel/database"
    "hostel/models"
)

func GetAllStudents(c *fiber.Ctx) error {
    var students []models.Student
    if err := database.DB.Preload("Room").Find(&students).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch students"})
    }

    // Optionally fetch payment status for each student
    var payments []models.Payment
    database.DB.Find(&payments)

    return c.JSON(students)
}

type RoomWithStudents struct {
    RoomNumber   string   `json:"room_number"`
    MaxBeds      uint     `json:"max_beds"`
    OccupiedBeds uint     `json:"occupied_beds"`
    Students     []string `json:"students"` // list of student names
}

func GetAllRooms(c *fiber.Ctx) error {
    var rooms []models.Room
    if err := database.DB.Find(&rooms).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch rooms"})
    }

    var result []RoomWithStudents
    for _, room := range rooms {
        var students []models.Student
        database.DB.Where("room_id = ?", room.ID).Find(&students)

        studentNames := []string{}
        for _, s := range students {
            studentNames = append(studentNames, s.Name)
        }

        result = append(result, RoomWithStudents{
            RoomNumber:   room.RoomNumber,
            MaxBeds:      3,
            OccupiedBeds: uint(len(students)),
            Students:     studentNames,
        })
    }

    return c.JSON(result)
}


func GetAllCheckInOutRecords(c *fiber.Ctx) error {
	var records []models.CheckInOut
	if err := database.DB.Preload("Student").Find(&records).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch records"})
	}

	var response []fiber.Map
	for _, r := range records {
		checkOutStr := "Not checked out"
		if r.CheckOut != nil {
			checkOutStr = r.CheckOut.Format("2006-01-02 15:04:05")
		}

		response = append(response, fiber.Map{
			"student_name": r.Student.Name,
			"roll_number":  r.Student.RollNumber,
			"check_in":     r.CheckIn.Format("2006-01-02 15:04:05"),
			"check_out":    checkOutStr,
		})
	}

	return c.JSON(response)
}
func GetPayments(c *fiber.Ctx) error {
	studentID := c.Query("student_id") // Optional filter

	var payments []models.Payment
	query := database.DB

	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}

	if err := query.Order("paid_at desc").Find(&payments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve payments"})
	}

	return c.JSON(fiber.Map{
		"payments": payments,
	})
}



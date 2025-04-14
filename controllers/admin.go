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


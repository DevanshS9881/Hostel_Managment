package models

import "time"

// Admin table
type Admin struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
}

// Room table
type Room struct {
    ID          uint   `gorm:"primaryKey"`
    RoomNumber  string `gorm:"unique;not null"`
    Capacity    int    // Number of beds
    Occupied    int    // Currently occupied beds
    Students    []Student
}

// Student table
type Student struct {
    ID         uint   `gorm:"primaryKey"`
    Name       string `gorm:"not null"`
	Password   string `gorm:"not null"`
    RollNumber string `gorm:"unique;not null"`
    Email      string `gorm:"unique;not null"`
    Phone      string
    RoomID     uint   // Foreign key to Room
    Room       Room   `gorm:"foreignKey:RoomID"`
}

// Payment table
type Payment struct {
    ID         uint      `gorm:"primaryKey"`
    StudentID  uint      `gorm:"not null"`
    Student    Student   `gorm:"foreignKey:StudentID"`
    Amount     float64   `gorm:"not null"`
    PaidAt     time.Time `gorm:"autoCreateTime"`
    TransactionID string `gorm:"type:varchar(100);uniqueIndex"`
    Description string
}

// CheckInOut table
type CheckInOut struct {
    ID         uint       `gorm:"primaryKey"`
    StudentID  uint       `gorm:"not null"`
    Student    Student    `gorm:"foreignKey:StudentID"`
    CheckIn    time.Time  `gorm:"not null"`
    CheckOut   *time.Time // Nullable
}

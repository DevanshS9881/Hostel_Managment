package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"hostel/models" // Replace with your actual module name
)

var DB *gorm.DB

// Convert string to uint
func Convert(port string) uint {
	u64, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		log.Fatal("Error:", err)
	}
	return uint(u64)
}

// Load a key from .env
func Load(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	return os.Getenv(key)
}

func InitDB() error {
	// Build DSN using values from .env
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Load("user"),
		Load("password"),
		Load("host"),
		Load("port"),
		Load("dbname"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
		return err
	}

	// Auto-migrate all tables
	err = DB.AutoMigrate(
		&models.Admin{},
		&models.Room{},
		&models.Student{},
		&models.Payment{},
		&models.CheckInOut{},
	)
	if err != nil {
		log.Fatal("AutoMigrate failed: ", err)
		return err
	}

	log.Println("✅ Database connected & tables migrated successfully")
	return nil
}

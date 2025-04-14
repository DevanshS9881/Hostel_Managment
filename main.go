package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // CORS middleware
	"github.com/gofiber/fiber/v2/middleware/logger"

	"hostel/database"
	"hostel/routes"
)

func main() {
	// Step 1: Initialize DB
	err := database.InitDB()
	if err != nil {
		log.Fatal("Database initialization failed: ", err)
	}

	// Step 2: Create a new Fiber app
	app := fiber.New()

	// Step 3: Enable CORS (allow all origins for dev)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Optional: Log every request (for debugging)
	app.Use(logger.New())

	// Step 4: Serve frontend static files
	app.Static("/", "./frontend")

	// Step 5: Setup home and API routes
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("🏠 Hostel Management System API is up!")
	})
	routes.SetRoutes(app)

	// Step 6: Start server
	log.Println("🚀 Server is running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

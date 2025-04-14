package routes

import (
	"hostel/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App){
	app.Post("/register_s",controllers.StudentSignIn)
	app.Post("/login_a",controllers.AdminLogin)
	app.Post("/login_s",controllers.StudentLogin)
	app.Get("/getStudents",controllers.GetAllStudents)
	app.Get("/getRooms",controllers.GetAllRooms)
	app.Get("/getStudent/:id", controllers.GetStudentByID)
	app.Post("/checkin", controllers.CheckInStudent)
    app.Post("/checkout", controllers.CheckOutStudent)
	app.Get("/checkinoutRecords", controllers.GetAllCheckInOutRecords)
	app.Post("/makePayment", controllers.CreatePayment)


}
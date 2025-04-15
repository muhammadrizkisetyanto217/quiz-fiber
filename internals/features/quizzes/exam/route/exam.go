package route

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	examController "quiz-fiber/internals/features/quizzes/exam/controller"
	userController "quiz-fiber/internals/features/user/auth/controller"
)

func ExamRoute(app *fiber.App, db *gorm.DB) {
	// 🔒 Grup API dengan Auth Middleware
	api := app.Group("/api", userController.AuthMiddleware(db))

	// 📝 Exam Routes
	examCtrl := examController.NewExamController(db)
	examRoutes := api.Group("/exams")

	examRoutes.Get("/", examCtrl.GetExams)                     // GET /api/exams
	examRoutes.Get("/:id", examCtrl.GetExam)                   // GET /api/exams/:id
	examRoutes.Get("/unit/:unitId", examCtrl.GetExamsByUnitID) // GET /api/exams/unit/:unitId
	examRoutes.Post("/", examCtrl.CreateExam)                  // POST /api/exams
	examRoutes.Put("/:id", examCtrl.UpdateExam)                // PUT /api/exams/:id
	examRoutes.Delete("/:id", examCtrl.DeleteExam)             // DELETE /api/exams/:id

	// ✅ User Exam Routes
	userExamCtrl := examController.NewUserExamController(db)
	userExamRoutes := api.Group("/user-exams")

	userExamRoutes.Post("/", userExamCtrl.Create)                  // POST /api/user-exams
	userExamRoutes.Get("/", userExamCtrl.GetAll)                   // GET /api/user-exams
	userExamRoutes.Get("/user/:user_id", userExamCtrl.GetByUserID) // GET /api/user-exams/user/:user_id
	userExamRoutes.Get("/:id", userExamCtrl.GetByID)               // GET /api/user-exams/:id
	userExamRoutes.Delete("/:id", userExamCtrl.Delete)             // DELETE /api/user-exams/:id
}

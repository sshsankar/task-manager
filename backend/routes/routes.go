package routes

import (
    "github.com/gofiber/fiber/v2"
)

// SetupRoutes - Registers API routes
func SetupRoutes(app *fiber.App) {
    app.Post("/tasks", CreateTask)
    app.Get("/tasks", GetTasks)
    app.Get("/ai-tasks", GetTaskSuggestions) // ✅ AI API Route
    app.Delete("/tasks/:id", DeleteTask)
}

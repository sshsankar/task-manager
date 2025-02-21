package routes

import (
    "backend/database"
    "backend/models"
    "encoding/json"
    "net/http"
    "bytes"
    "os"
    "log"
    "github.com/gofiber/fiber/v2"
)

// CreateTask - Adds a new task
func CreateTask(c *fiber.Ctx) error {
    task := new(models.Task)
    if err := c.BodyParser(task); err != nil {
        return err
    }

    _, err := database.DB.Exec("INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)", task.Title, task.Description, "pending")
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{"message": "Task created"})
}

// GetTasks - Fetch all tasks
func GetTasks(c *fiber.Ctx) error {
    rows, err := database.DB.Query("SELECT id, title, description, status FROM tasks")
    if err != nil {
        return err
    }
    defer rows.Close()

    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
        tasks = append(tasks, task)
    }

    return c.JSON(tasks)
}

// GetTaskSuggestions - Fetch AI-generated task suggestions using OpenAI API
func GetTaskSuggestions(c *fiber.Ctx) error {
    openaiAPIKey := os.Getenv("OPENAI_API_KEY")

    prompt := "Suggest 3 new tasks for a project management system."

    requestData := map[string]interface{}{
        "model": "gpt-3.5-turbo",
        "messages": []map[string]string{
            {"role": "system", "content": "You are a helpful assistant."},
            {"role": "user", "content": prompt},
        },
    }

    requestBody, _ := json.Marshal(requestData)

    req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return c.Status(500).SendString(err.Error())
    }
    defer resp.Body.Close()

    var response map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&response)

    log.Println("AI Response:", response) // Debugging

    return c.JSON(response)
}

// DeleteTask - Deletes a task by ID
func DeleteTask(c *fiber.Ctx) error {
    taskID := c.Params("id")

    // Check if the task exists
    var task models.Task
    err := database.DB.QueryRow("SELECT id, title, description, status FROM tasks WHERE id = $1", taskID).Scan(&task.ID, &task.Title, &task.Description, &task.Status)
    if err != nil {
        if err.Error() == "sql: no rows in result set" {
            return c.Status(404).JSON(fiber.Map{"message": "Task not found"})
        }
        return c.Status(500).SendString(err.Error())
    }

    // Proceed with deleting the task
    _, err = database.DB.Exec("DELETE FROM tasks WHERE id = $1", taskID)
    if err != nil {
        return c.Status(500).SendString(err.Error())
    }

    return c.JSON(fiber.Map{"message": "Task deleted"})
}

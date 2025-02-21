package main

import (
    "backend/database"
    "backend/routes"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "github.com/joho/godotenv"
    "log"
    "os"
)

// WebSocket clients
var clients = make(map[*websocket.Conn]bool)

func main() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    openaiKey := os.Getenv("OPENAI_API_KEY")
    log.Println("Loaded OpenAI API Key:", openaiKey) // Debugging

    // Connect to database
    database.ConnectDB()

    // Initialize Fiber app
    app := fiber.New()

    // Enable CORS
    
    app.Use(func(c *fiber.Ctx) error {
    	c.Set("Access-Control-Allow-Origin", "*") // Allow all origins
    	c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
    	c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    	// Handle preflight requests (OPTIONS method)
    	if c.Method() == "OPTIONS" {
        	return c.SendStatus(fiber.StatusNoContent)
    	}

    	return c.Next()
    })


    // WebSocket route for real-time updates
    app.Get("/ws", websocket.New(func(c *websocket.Conn) {
        clients[c] = true
        defer func() {
            delete(clients, c)
            c.Close()
        }()

        for {
            _, msg, err := c.ReadMessage()
            if err != nil {
                log.Println("WebSocket Error:", err)
                break
            }
            broadcast(msg)
        }
    }))

    // Register API routes
    routes.SetupRoutes(app)

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default to 8080 if not set
    }
    log.Println("Server running on port:", port)
    log.Fatal(app.Listen(":" + port))
}

// Broadcast function to send updates to all WebSocket clients
func broadcast(message []byte) {
    for client := range clients {
        err := client.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Println("WebSocket Broadcast Error:", err)
            client.Close()
            delete(clients, client)
        }
    }
}

package routes

import (
    "backend/database"
    "backend/models"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "os"
    "time"
)

// Signup Route
func Register(c *fiber.Ctx) error {
    user := new(models.User)
    if err := c.BodyParser(user); err != nil {
        return err
    }

    // Hash password before saving
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)

    // Insert user into database
    _, err = database.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{"message": "User registered successfully"})
}

// Login Route
func Login(c *fiber.Ctx) error {
    user := new(models.User)
    if err := c.BodyParser(user); err != nil {
        return err
    }

    var storedPassword string
    err := database.DB.QueryRow("SELECT password FROM users WHERE email = $1", user.Email).Scan(&storedPassword)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
    }

    // Compare hashed password
    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
    }

    // Generate JWT token
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["email"] = user.Email
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

    t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{"token": t})
}

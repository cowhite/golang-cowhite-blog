package users

import (
    "github.com/gofiber/fiber/v2"
)

func SetupUrls(app fiber.Router) {
    app.Post("/api/v1/signup", SignupView)
    app.Post("/api/v1/login", LoginView)
}
package main

import (
	"github.com/gofiber/fiber/v2"
    "src/blog"
    "src/users"
)

func main() {
  app := fiber.New()

  blog.SetupUrls(app)
  users.SetupUrls(app)

  app.Listen(":3000")
}

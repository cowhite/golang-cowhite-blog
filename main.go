package main

import (
	"github.com/gofiber/fiber/v2"
    "src/blog"
)

func main() {
  app := fiber.New()

  blog.SetupUrls(app)

  app.Listen(":3000")
}

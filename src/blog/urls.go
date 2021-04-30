package blog

import (
	"github.com/gofiber/fiber/v2"
)

func SetupUrls(app fiber.Router) {
	app.Get("/", HomeView)
	app.Post("/categories/new", AddCategoryView)
	app.Post("/blogposts/new", AddBlogPostView)
}
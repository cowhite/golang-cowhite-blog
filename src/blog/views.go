package blog

import (
	"github.com/gofiber/fiber/v2"
    // "strconv"
    "fmt"
)

func HomeView(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"Hello": "World",
	})
}

func AddCategoryView(c *fiber.Ctx) error {
	category := new(Category)
	c.BodyParser(category)

	errors := ValidateNewCategory(c, category)
	if errors.Error {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	category = NewCategory(c, category)
	fmt.Println("res,,")
	fmt.Println(category)
	return c.Status(fiber.StatusCreated).JSON(category)
}

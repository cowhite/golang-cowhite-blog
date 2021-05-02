package blog

import (
	"github.com/gofiber/fiber/v2"
    // "strconv"
    "fmt"
    "src/utils"
)

func HomeView(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"Hello": "World",
	})
}

func AddCategoryView(c *fiber.Ctx) error {
	fmt.Println("AddCategoryView...")

	category := new(Category)
	c.BodyParser(category)

	errors := ValidateNewCategory(c, category)
	if errors != nil {
		var response utils.Response = utils.GetErrorResponse(errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	category = NewCategory(c, category)
	return c.Status(fiber.StatusCreated).JSON(category)
}

func AddBlogPostView(c *fiber.Ctx) error {
	fmt.Println("AddBlogPostView...")
	blogpost := new(BlogPost)
	err := c.BodyParser(blogpost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	// errors := ValidateNewBlogPost(c, blogpost)
	// if errors.Error {
	// 	return c.Status(fiber.StatusBadRequest).JSON(errors)
	// }

	blogpost = NewBlogPost(c, blogpost)
	fmt.Println("res,,")
	fmt.Println(blogpost)
	return c.Status(fiber.StatusCreated).JSON(blogpost)
}

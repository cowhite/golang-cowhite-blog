package utils

// import (
// 	"github.com/gofiber/fiber/v2"
// )

type ErrorResponse struct {
    Message     string
    Error       bool
    Errors      []string
}

// func GetResponse(c *fiber.Ctx, res fiber.Map) error {
// 	var status int
// 	if res.error {
// 		status = fiber.StatusBadRequest
// 	} else if res.created == true {
// 		status = fiber.StatusCreated
// 	} else {
// 		status = fiber.StatusOK
// 	}
// 	return c.Status(status).JSON(res)
// }


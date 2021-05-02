package users

import (
    "github.com/gofiber/fiber/v2"
    // "strconv"
    "log"
    "fmt"
    "src/utils"
)

func SignupView(c *fiber.Ctx) error {
    fmt.Println("SignupView...")

    user := new(User)
    err := c.BodyParser(user)
    if err != nil {
    	log.Fatal(err)
    }

    var errors []utils.ErrorResponse = user.ValidateNewUser(c)
    if errors != nil {
        var response utils.Response = utils.GetErrorResponse(errors)
        return c.Status(fiber.StatusBadRequest).JSON(response)
    }

    user.Create(c)
    return c.Status(fiber.StatusCreated).JSON(user)
}

func LoginView(c *fiber.Ctx) error {
	user_login := new(UserLogin)
	err := c.BodyParser(user_login)
    if err != nil {
    	log.Fatal(err)
    }
    var response UserLoginResponse = user_login.ValidateLogin(c)
    if response.Error {
		return c.Status(fiber.StatusBadRequest).JSON(response)
    } else {
		return c.Status(fiber.StatusOK).JSON(response)
	}

}
package utils

import (
	"src/config"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gosimple/slug"
)

type ErrorResponse struct {
	Key string
	Value string
	Tag string
}

type Response struct {
    Message     string
    Error       bool
    // Errors      []string
    Errors []ErrorResponse
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

func GetErrorResponse(errors []ErrorResponse) Response {
	return Response{
		Message: "Please check errors",
		Error: true,
		Errors: errors,
	}
}

func GenerateUniqueSlug(c *fiber.Ctx, title string, collection_name string, slug_field string) string {
	if title == "" {
		return ""
	}
	if slug_field == "" {
		slug_field = "slug"
	}

	var temp_title string = title
	var title_slug string
	var err error
	var i int = 0

	JSONData := &bson.D{}
	db := config.Connect().Db
	collection := db.Collection(collection_name)

	for true {
		title_slug = slug.Make(temp_title)
		err = collection.FindOne(c.Context(), bson.M{"slug": title_slug}).Decode(JSONData)
		if err != nil {
			break
		}
		i++
		temp_title = title + "-" + strconv.Itoa(i)
	}
	return title_slug
}


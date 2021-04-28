package blog

import (
	"time"
	"fmt"
	"log"
	"src/config"
	"strconv"

	"github.com/gofiber/fiber/v2"
	// "github.com/fatih/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gosimple/slug"

    "src/utils"

)

// Todo - todo model
type BlogPost struct {
	Id        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string   `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

}

type Category struct {
	Id 				primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title           string `json:"title"`
	Slug 			string `json:"slug"`
	description		string `json:"description"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

func Add(c *fiber.Ctx, x int, y int) int {
	fmt.Println(time.Now())
	NewBlogPost(c, "ddd")
	return x + y
}

func NewBlogPost(c *fiber.Ctx, title string) *BlogPost {
	blog := new(BlogPost)
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	blog.Title = title

	db := config.Connect().Db

	blog_collection := db.Collection("blogs")
	result, err := blog_collection.InsertOne(c.Context(), blog)

	if err != nil {
		log.Fatal(err)
	}

	err = blog_collection.FindOne(c.Context(), bson.M{"_id": result.InsertedID}).Decode(&blog)

	if err != nil {
		log.Fatal(err)
	}

	return blog
}

func ValidateNewCategory(c *fiber.Ctx, category *Category) *utils.ErrorResponse {
	db := config.Connect().Db

	// category := new(Category)

	// Validate
	category_collection := db.Collection("categories")
	err := category_collection.FindOne(c.Context(), bson.M{"title": category.Title}).Decode(&category)

	// var res Response
	if err == nil {
		return &utils.ErrorResponse{
			Error: true,
			Message: "Title already exists",
			Errors: []string{"Title already exists"},
		}
	}
	return &utils.ErrorResponse{
		Error: false,
		Message: "",
	}

}

func NewCategory(c *fiber.Ctx, category *Category) *Category {
	db := config.Connect().Db

	category_collection := db.Collection("categories")

	// title not found, so create it
	title_slug := GenerateUniqueSlug(c, category.Title, "categories", "slug")
	fmt.Println("title_slug")
	fmt.Println(title_slug)
	category.Slug = title_slug
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	result, err := category_collection.InsertOne(c.Context(), category)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result")
	fmt.Println(result)
	fmt.Println("err2")
	fmt.Println(err)

	err = category_collection.FindOne(c.Context(), bson.M{"_id": result.InsertedID}).Decode(&category)

	return category

	// return fiber.Map{
	// 	"error": false,
	// 	"results": fiber.Map{
	// 		"id": result.InsertedID,
	// 		"title": category.Title,
	// 		"slug": category.Slug,
	// 		"updatedAt": category.CreatedAt,
	// 		"createdAt": category.UpdatedAt,
	// 	},
	// 	"message": "success",
	// }

}

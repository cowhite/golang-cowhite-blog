package blog

import (
	"time"
	"fmt"
	"log"
	"src/config"
	// "strconv"

	"github.com/gofiber/fiber/v2"
	// "github.com/fatih/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "github.com/gosimple/slug"

    "src/utils"

)

// Todo - todo model
type BlogPost struct {
	Id        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string   `json:"title"`
	Slug 			string `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CategoryIds []string `json:"category_ids"`

}

type Category struct {
	Id 				primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title           string `json:"title"`
	Slug 			string `json:"slug"`
	description		string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BlogPostCategory struct {
	Id 				primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BlogPostId      string `json:"blogpost_id"`
	CategoryId      string `json:"category_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// func Add(c *fiber.Ctx, x int, y int) int {
// 	fmt.Println(time.Now())
// 	NewBlogPost(c, "ddd")
// 	return x + y
// }

// func NewBlogPost(c *fiber.Ctx, title string) *BlogPost {
// 	blog := new(BlogPost)
// 	blog.CreatedAt = time.Now()
// 	blog.UpdatedAt = time.Now()
// 	blog.Title = title

// 	db := config.Connect().Db

// 	blog_collection := db.Collection("blogs")
// 	result, err := blog_collection.InsertOne(c.Context(), blog)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = blog_collection.FindOne(c.Context(), bson.M{"_id": result.InsertedID}).Decode(&blog)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return blog
// }

func ValidateNewCategory(c *fiber.Ctx, category *Category) []utils.ErrorResponse {
	db := config.Connect().Db

	// category := new(Category)

	// Validate
	category_collection := db.Collection("categories")
	err := category_collection.FindOne(c.Context(), bson.M{"title": category.Title}).Decode(&category)

	var errors []utils.ErrorResponse

	// var res Response
	if err == nil {
		errors = append(errors, utils.ErrorResponse{
			Key: "title",
			Value: "Title already exists",
			// Error: true,
			// Message: "Title already exists",
			// Errors: []string{"Title already exists"},
		})
		return errors
	}
	return []utils.ErrorResponse{}
	// return &utils.ErrorResponse{
	// 	Error: false,
	// 	Message: "",
	// }

}

func NewCategory(c *fiber.Ctx, category *Category) *Category {
	db := config.Connect().Db

	category_collection := db.Collection("categories")

	// title not found, so create it
	title_slug := utils.GenerateUniqueSlug(c, category.Title, "categories", "slug")
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

func NewBlogPost(c *fiber.Ctx, blogpost *BlogPost) *BlogPost {
	db := config.Connect().Db

	blogpost_collection := db.Collection("blogposts")

	// title not found, so create it
	title_slug := utils.GenerateUniqueSlug(c, blogpost.Title, "blogposts", "slug")
	fmt.Println("title_slug")
	fmt.Println(title_slug)
	blogpost.Slug = title_slug
	blogpost.CreatedAt = time.Now()
	blogpost.UpdatedAt = time.Now()
	result, err := blogpost_collection.InsertOne(c.Context(), blogpost)
	if err != nil {
		log.Fatal(err)
	}

	categories := GetCategories(c, blogpost.CategoryIds)
	NewBlogPostCategoryBatch(c, result.InsertedID.(primitive.ObjectID).Hex(), categories)

	err = blogpost_collection.FindOne(c.Context(), bson.M{"_id": result.InsertedID}).Decode(&blogpost)

	return blogpost
}

func NewBlogPostCategoryBatch(c *fiber.Ctx, blogpostId string, categories []Category) []BlogPostCategory {
	db := config.Connect().Db

	blogpost_category_collection := db.Collection("blogpost_categories")

	var blogpost_categories []BlogPostCategory
	var blogpost_category BlogPostCategory

	// TODO: Check for duplicates

	for _, category := range categories {
		blogpost_category = BlogPostCategory{
			BlogPostId: blogpostId, CategoryId: category.Id.Hex(),
			CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		blogpost_categories = append(blogpost_categories, blogpost_category)
	}

	var ui []interface{}
	for _, t := range blogpost_categories{
	    ui = append(ui, t)
	}

	_, err := blogpost_category_collection.InsertMany(c.Context(), ui)
	if err != nil {
		log.Fatal(err)
	}

	return blogpost_categories
}

func GetCategories(c *fiber.Ctx, categoryIds []string) []Category {
	db := config.Connect().Db
	var categoryIdObjects []primitive.ObjectID
	var categoryIdObject primitive.ObjectID
	var err error;

	for _, categoryId := range categoryIds {
		categoryIdObject, err = primitive.ObjectIDFromHex(categoryId)
		if err != nil {
			log.Fatal(err)
		}
		categoryIdObjects = append(categoryIdObjects, categoryIdObject)
	}
	fmt.Println("categoryIdObjects")
	fmt.Println(categoryIdObjects)

	category_collection := db.Collection("categories")

	var categories []Category

	results, err2 := category_collection.Find(c.Context(), bson.M{"_id": bson.M{"$in": categoryIdObjects}})
	results.All(c.Context(), &categories)
	fmt.Println("categories")
	fmt.Println(categories)

	if err2 != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return categories
}


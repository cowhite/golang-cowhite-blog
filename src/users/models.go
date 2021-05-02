package users

import (
	"time"
	"fmt"
	"log"
	"strings"

	"src/config"
	// "strconv"

	"github.com/gofiber/fiber/v2"
	fiberUtils "github.com/gofiber/fiber/v2/utils"
	// // "github.com/fatih/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/go-playground/validator"

	// "github.com/gosimple/slug"


    "src/utils"

)

type UserInterface interface {}

type User struct {
	Id        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName     string   `json:"first_name"`
	LastName 	  string `json:"last_name"`
	Email 		  string `json:"email" validate:"required,email"`
	Username 	  string `json:"username"`
	Password 	  string `json:"password"`
	IsActive 	  bool `json:"is_active"`
	IsAdmin 	  bool `json:"is_admin"`
	IsStaff 	  bool `json:"is_staff"`
	Token 		  string `json:"token"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

}

type UserDetail struct {
	Id        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName     string   `json:"first_name"`
	LastName 	  string `json:"last_name"`
	Email 		  string `json:"email" validate:"required,email"`
	Username 	  string `json:"username"`
	IsActive 	  bool `json:"is_active"`
	IsAdmin 	  bool `json:"is_admin"`
	IsStaff 	  bool `json:"is_staff"`
	Token 		  string `json:"token"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

}

type UserLogin struct {
	Password string `json:"password" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UserLoginResponse struct {
	User UserDetail `json:user`
	Errors []utils.ErrorResponse `json:errors`
	Error bool `json:error`
}

func (user_login UserLogin) ValidateLogin(c *fiber.Ctx) UserLoginResponse {
	db := config.Connect().Db
	users_collection := db.Collection("users")

	var user User
	var errors []utils.ErrorResponse

	// Basic validations
	validate := validator.New()
	err := validate.Struct(user_login)
	fmt.Println("err")
	fmt.Println(err)
	fmt.Println(user_login)
	fmt.Println(user_login.Email)
	fmt.Println(user_login.Password)

	if err != nil {
        for _, err2 := range err.(validator.ValidationErrors) {
        	fmt.Println(err2)
            var element utils.ErrorResponse
            element.Key = err2.Field()
            element.Tag = err2.Tag()
            errors = append(errors, element)
        }
	} else {

		// Get password from db
		err = users_collection.FindOne(c.Context(), bson.M{"email": strings.ToLower(user_login.Email)}).Decode(&user)

		if err != nil {
			// Check if email exists
			errors = append(errors, utils.ErrorResponse{
				Key: "base",
				Value: "Invalid credentials",
			})
		} else {
			// Check if password is correct
			err = CheckPassword(user_login.Password, user.Password)
			if err != nil {
				errors = append(errors, utils.ErrorResponse{
					Key: "base",
					Value: "Invalid credentials",
				})
			}
		}
	}

	var is_error bool = false
	if errors != nil {
		is_error = true
	}

	// Prepare data to send to user
	var user_detail UserDetail
	err = users_collection.FindOne(c.Context(), bson.M{"email": strings.ToLower(user_login.Email)}).Decode(&user_detail)

	var response UserLoginResponse = UserLoginResponse{
		Errors: errors,
		Error: is_error,
	}
	if is_error == false {
		response.User = user_detail
	}
	return response
}



func (user User) ValidateNewUser(c *fiber.Ctx) []utils.ErrorResponse {

	var errors []utils.ErrorResponse;

	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		fmt.Println(err.(validator.ValidationErrors))
        for _, err2 := range err.(validator.ValidationErrors) {
            var element utils.ErrorResponse
            element.Key = err2.Field()
            element.Tag = err2.Tag()
            errors = append(errors, element)
        }
	}

	db := config.Connect().Db
	users_collection := db.Collection("users")

	// Unique email validation
	err = users_collection.FindOne(c.Context(), bson.M{"email": strings.ToLower(user.Email)}).Decode(&user)

	if err == nil {
		errors = append(errors, utils.ErrorResponse{
			Key: "Email",
			Value: "Email already exists",
		})
	}
	return errors
}

func (user *User) Create(c *fiber.Ctx) {
	// Hash password, dont save plain password in db
	if user.Password != "" {
		user.Password = user.SetPassword(user.Password)
	}

	// Lowercase email
	user.Email = strings.ToLower(user.Email)

	// For now, username and email are same
	user.Username = user.Email
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Token = fiberUtils.UUIDv4()

	// Save to db
	db := config.Connect().Db
	users_collection := db.Collection("users")

	result, err := users_collection.InsertOne(c.Context(), user)

	if err != nil {
		log.Fatal(err)
	}

	// Fetch user from db and return it
	err = users_collection.FindOne(c.Context(), bson.M{"_id": result.InsertedID}).Decode(&user)

}

func (user User) GetFullName() string {
	return user.FirstName + " " + user.LastName
}

func (user User) SetPassword(password string) string {
	// Set hashed password
	password, _ = HashPassword(password)
	return password
}

func (user User) Get(c *fiber.Ctx) User {
	db := config.Connect().Db
	users_collection := db.Collection("users")

	err := users_collection.FindOne(c.Context(), bson.M{"_id": user.Id.Hex()}).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	return user

}




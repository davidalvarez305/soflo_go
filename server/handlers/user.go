package handlers

import (
	"fmt"
	"os"

	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/davidalvarez305/soflo_go/server/database"
	"github.com/davidalvarez305/soflo_go/server/models"
	"github.com/davidalvarez305/soflo_go/server/sessions"
	"github.com/davidalvarez305/soflo_go/server/types"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	var u types.User
	var user types.UserRequestBody
	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	if user.Password == "" || user.Username == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "Missing Fields.",
		})
	}

	hashedPassword, err3 := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err3 != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error hashing password",
		})
	}

	u.Username = user.Username
	u.Password = hashedPassword

	fmt.Printf("User %v", u)

	data, err2 := actions.CreateUser(u)

	if err2 != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Unable to Create User.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": data,
	})
}

func GetUser(c *fiber.Ctx) error {
	sess, err := sessions.Sessions.Get(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Unable to retrieve cookie.",
		})
	}

	k := sess.Get(os.Getenv("COOKIE_NAME"))

	if k == nil {
		return c.Status(404).JSON(fiber.Map{
			"data": "Not found.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": k,
	})
}

func Logout(c *fiber.Ctx) error {

	sess, err := sessions.Sessions.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Unable to get cookie.",
		})
	}

	k := sess.Get(os.Getenv("COOKIE_NAME"))

	if k == nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Not found.",
		})
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Unable to destroy session.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": "Logged out!",
	})
}

func Login(c *fiber.Ctx) error {
	var u types.UserRequestBody
	var user models.User
	err := c.BodyParser(&u)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Bad Input.",
		})
	}

	result := database.DB.Where("username = ?", &u.Username).First(&user)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Incorrect username.",
		})
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(u.Password))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Incorrect password.",
		})
	}

	id := sessions.Sessions.KeyGenerator()

	sess, err := sessions.Sessions.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Unable to get cookie.",
		})
	}

	sess.Set(os.Getenv("COOKIE_NAME"), id)

	if err := sess.Save(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Unable to save session.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": id,
	})
}

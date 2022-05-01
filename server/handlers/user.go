package handlers

import (
	"os"

	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/davidalvarez305/soflo_go/server/database"
	"github.com/davidalvarez305/soflo_go/server/models"
	"github.com/davidalvarez305/soflo_go/server/sessions"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	var user models.User
	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	data, err2 := actions.CreateUser(user)

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
		panic(err)
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
	var user models.User
	type body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var reqBody body
	err := c.BodyParser(&reqBody)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Bad Input.",
		})
	}

	result := database.DB.Where("email = ?", &reqBody.Email).First(&user)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Incorrect e-mail.",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(user.Password))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Incorrect password.",
		})
	}

	id := sessions.Sessions.KeyGenerator()

	sess, err := sessions.Sessions.Get(c)
	if err != nil {
		panic(err)
	}

	sess.Set(os.Getenv("COOKIE_NAME"), id)

	if err := sess.Save(); err != nil {
		panic(err)
	}

	return c.Status(200).JSON(fiber.Map{
		"data": id,
	})
}

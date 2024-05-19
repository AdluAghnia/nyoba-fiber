package handler

import (
	"time"

	"github.com/AdluAghnia/nyoba-fiber/config"
	"github.com/AdluAghnia/nyoba-fiber/models"
	"github.com/AdluAghnia/nyoba-fiber/repository"
	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
)

func LoginHandler(c *fiber.Ctx) error {
	// get user information from FORM
	username := c.FormValue("name")
	password := c.FormValue("passwd")

	// find usr credentials

	user, err := repository.FindByCredentials(username, password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// create JWT claims, include the user id and expiry time
	day := time.Hour * 24

	claims := jtoken.MapClaims{
		"ID":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(day * 1).Unix(),
	}

	// create token

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	// Generate encode token and send it as response
	t, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(models.LoginResponse{
		Token: t,
	})
}

func FrontpageHandler(c *fiber.Ctx) error {
	// Render HTML Template
	return c.Render("index", fiber.Map{
		"Title": "This is front page",
	}, "layouts/main")
}

func Protected(c *fiber.Ctx) error {
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	name := claims["username"].(string)

	return c.SendString("Hello" + name)
}

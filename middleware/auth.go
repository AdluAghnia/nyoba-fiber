package middlewares

import (
	"strings"

	"github.com/AdluAghnia/nyoba-fiber/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Middleware JWT function

func JwtMiddleware(c *fiber.Ctx) error {
	// Get the user form context and return it
	tokenString := c.Get("user")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing JWT")
	}

	tokenString = strings.Replace(tokenString, "Bearer :", "", 1)

	// Parse JWT Token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return config.Secret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid JWT Token")
	}

	return c.Next()

}

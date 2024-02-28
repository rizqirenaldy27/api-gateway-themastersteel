package middleware

import (
	"context"
	"service-user/config"
	"service-user/helpers"
	"service-user/model"

	"github.com/gofiber/fiber/v2"
)

func Authentication(c *fiber.Ctx) error {
	accessToken := c.Get("access_token")

	if len(accessToken) == 0 {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token: Access token missing")
	}

	claims, err := helpers.VerifyToken(accessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token: Failed to verify token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token: Missing email claim")
	}

	var user model.User
	db := config.GetPgxPool()

	query := "SELECT id, email, password FROM users WHERE email = $1"
	if err := db.QueryRow(context.Background(), query, email).Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token: User not found")
	}

	c.Locals("user", user)

	return c.Next()
}

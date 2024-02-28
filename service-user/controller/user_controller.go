package controller

import (
	"service-user/helpers"
	"service-user/model"

	"service-user/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WebResponse struct {
	Code   int
	Status string
	Data   interface{}
}

func Register(c *fiber.Ctx) error {
	var requestBody model.User
	pgPool := config.GetPgxPool()

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userId := uuid.New()

	ctx, cancel := config.NewPostgresContext()
	defer cancel()

	_, err := pgPool.Exec(ctx, `
		INSERT INTO users (id, email, password) 
		VALUES ($1, $2, $3)`,
		userId, requestBody.Email, helpers.HashPassword([]byte(requestBody.Password)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":   201,
		"status": "OK",
		"data":   requestBody.Email,
	})
}

func Login(c *fiber.Ctx) error {
	pgPool := config.GetPgxPool()

	var requestBody model.User
	var result model.User

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := config.NewPostgresContext()
	defer cancel()

	err := pgPool.QueryRow(ctx, `
		SELECT id, email, password 
		FROM users 
		WHERE email = $1`, requestBody.Email).
		Scan(&result.Id, &result.Email, &result.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if !helpers.ComparePassword([]byte(result.Password), []byte(requestBody.Password)) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	accessToken := helpers.SignToken(requestBody.Email)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":        fiber.StatusOK,
		"status":      "OK",
		"accessToken": accessToken,
		"data":        result,
	})
}

func Auth(c *fiber.Ctx) error {
	return c.JSON("OK")
}

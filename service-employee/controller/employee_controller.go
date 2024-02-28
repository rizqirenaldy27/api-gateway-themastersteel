package controller

import (
	"fmt"
	"net/http"
	"service-employee/config"
	"service-employee/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var user_uri string = "http://service-user:3001/user"

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func CreateEmployee(c *fiber.Ctx) error {
	pgPool := config.GetPgxPool()
	var requestBody model.Employee

	// Parse body request
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "Invalid request body",
		})
	}

	requestBody.Id = uuid.New().String()

	access_token := c.Get("access_token")
	if len(access_token) == 0 {
		return c.Status(401).SendString("Invalid token: Access token missing")
	}

	req, err := http.NewRequest("GET", user_uri+"/auth", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic(err)
	}

	// Set headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("access_token", access_token)

	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		panic(err)
	}
	defer resp.Body.Close()

	// Print the response
	// fmt.Println("Response Status:", resp.Status)
	// fmt.Println("Response Headers:", resp.Header)

	if resp.Status != "200 OK" {
		c.Status(401).SendString("invalid token")
	}

	ctx, cancel := config.NewPostgresContext()
	defer cancel()

	if _, err := pgPool.Exec(ctx, `
		INSERT INTO employees (id, name) 
		VALUES ($1, $2)`,
		requestBody.Id, requestBody.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create employee")
	}

	return c.Status(fiber.StatusCreated).JSON(WebResponse{
		Code:   fiber.StatusCreated,
		Status: "Created",
		Data:   requestBody,
	})
}

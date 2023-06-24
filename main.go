package main

import (
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	listenPort := flag.String("port", "8000", "Port to run the API server")
	flag.Parse()

	fmt.Println(*listenPort)

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	app.Get("/foo", handleFoo)
	apiv1.Get("/user", handleUser)
	app.Listen(":" + *listenPort)
}

func handleUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"user": "Me",
	})
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"msg": "Working!",
	})
}

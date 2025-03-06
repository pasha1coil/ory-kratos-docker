package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/ory/client-go"
	"log"
	"ory-kratos-docker/middleware"
	"os/signal"
	"syscall"
)

func handler(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := fiber.New()

	c := client.NewConfiguration()
	c.Servers = client.ServerConfigurations{{URL: "http://kratos:4433"}}
	ory := client.NewAPIClient(c)

	app.Use(middleware.KratosMiddleware(ory))

	app.Get("/", handler)
	app.Get("/public", func(c *fiber.Ctx) error {
		return c.SendString("Public content")
	})
	app.Get("/private", func(c *fiber.Ctx) error {
		return c.SendString("private content")
	})

	app.Get("/check", func(c *fiber.Ctx) error {
		sessionToken := c.Get("Cookie")
		return c.Status(fiber.StatusOK).SendString(sessionToken)
	})

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()

	fmt.Println("Shutting down server...")
}

package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
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

	//c := client.NewConfiguration()
	//c.Servers = client.ServerConfigurations{{URL: "http://kratos:4433"}}
	//ory := client.NewAPIClient(c)

	app.Use(middleware.OauthKeeperMiddleware())

	app.Get("/", handler)
	app.Get("/public", func(c *fiber.Ctx) error {
		return c.SendString("Public content")
	})
	app.Get("/private", func(c *fiber.Ctx) error {
		token, ok := c.Context().UserValue(middleware.CtxTokenKey).(string)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.SendString("private content, token: " + token)
	})

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()

	fmt.Println("Shutting down server...")
}

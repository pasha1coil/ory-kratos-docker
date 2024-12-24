package test

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/ory/client-go"
	"ory-kratos-docker/middleware"
	"os/signal"
	"syscall"
	"testing"
)

// for test
// ory proxy http://localhost:3000 --project 'pj-id'

func handler(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Test_Srv(t *testing.T) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	c := client.NewConfiguration()
	c.Servers = client.ServerConfigurations{{URL: "http://localhost:4000/.ory"}}
	ory := client.NewAPIClient(c)

	app := fiber.New()

	app.Use(middleware.KratosMiddleware(ory))

	app.Get("/", handler)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			t.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()

	fmt.Println("Shutting down server...")
}

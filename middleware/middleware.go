package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/ory/client-go"
)

func KratosMiddleware(ory *client.APIClient) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ory == nil {
			return fiber.NewError(fiber.StatusInternalServerError, "empty ory client")
		}

		sessionToken := ctx.Get("Cookie")
		if sessionToken == "" {
			fmt.Println("sessionToken nil")
			return ctx.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		session, err := getSessionFromToken(ctx, ory, sessionToken)
		if err != nil || !*session.Active {
			fmt.Println("!*session.Active", err.Error())
			return ctx.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		fmt.Println("Active session:", session.Id)
		fmt.Println("Active identity:", session.Identity.Id)

		return ctx.Next()
	}
}

func getSessionFromToken(ctx *fiber.Ctx, ory *client.APIClient, token string) (*client.Session, error) {
	session, _, err := ory.FrontendAPI.ToSession(ctx.Context()).Cookie(token).Execute()
	if err != nil {
		return nil, fmt.Errorf("error retrieving session: %w", err)
	}

	return session, nil
}

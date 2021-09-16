package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WebSocketUpgradeMiddleware(context *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client
	// requested upgrade to the WebSocket protocol.
	if websocket.IsWebSocketUpgrade(context) {
		context.Locals("allowed", true)
		return context.Next()
	}

	return fiber.ErrUpgradeRequired
}
func XApiKeyMiddleware(context *fiber.Ctx) error {
	// todo investigar por que no llega X_API_KEY
	/*
		requestApiKey := context.Get("X_API_KEY")
		serverApiKey := utils.GetEnvVariable("X_API_KEY")
		fmt.Println("requestApiKey")
		fmt.Println(requestApiKey)
		fmt.Println("serverApiKey")
		fmt.Println(serverApiKey)
		// context.h

		if requestApiKey != serverApiKey {
			return context.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"error": "unauthorized",
			})
		}
	*/

	return context.Next()
}

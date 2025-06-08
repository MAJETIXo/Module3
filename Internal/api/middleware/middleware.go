package middleware

import "github.com/gofiber/fiber/v2"

func Authorization(token string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}

}

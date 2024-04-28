package main

import (
	"github.com/gofiber/fiber/v2"
)

func Protected(c *fiber.Ctx) error {
	sess, err := globalSession.Get(c)
	if err != nil {
		return err
	}

	if sess.Fresh() || len(sess.Keys()) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	// Continue to the next middleware or route handling function
	return c.Next()
}

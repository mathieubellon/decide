package main

import (
	"github.com/gofiber/fiber/v2"
)

func Protect(c *fiber.Ctx) error {
	// Perform tasks before the route handling function
	println("Middleware: Request received")

	// Continue to the next middleware or route handling function
	return c.Next()
}

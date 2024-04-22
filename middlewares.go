package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func MakeSession(c *fiber.Ctx) error {
	// Perform tasks before the route handling function
	println("Middleware: Request received")

	// Continue to the next middleware or route handling function
	return c.Next()
}

func Onboarding(c *fiber.Ctx) error {
	// Perform tasks before the route handling function
	println("Middleware: Onboarding")
	sess, err := globalSession.Get(c) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}

	if sess.Get("userEmail") != nil {
		println("User is authenticated")
		var user User

		err := db.Find(&user, sess.Get("userID").(uint)).Error
		if err != nil {
			return err
		}

		log.Println(user.ID)
		log.Println(user.UUID)
		log.Println(user.WorkspaceID)

		return c.Next()
	}

	// Continue to the next middleware or route handling function
	return c.Next()
}

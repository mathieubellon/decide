package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"

	"github.com/shareed2k/goth_fiber"
)

func main() {
	godotenv.Load()
	app := fiber.New()
	app.Use(logger.New())

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://127.0.0.1:8000/auth/callback/google"),
	)

	app.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	app.Get("/auth/callback/:provider", func(ctx *fiber.Ctx) error {
		CompleteUserAuthOptions := goth_fiber.CompleteUserAuthOptions{
			ShouldLogout: false,
		}
		user, err := goth_fiber.CompleteUserAuth(ctx, CompleteUserAuthOptions)
		if err != nil {
			log.Fatal(err)
		}

		return ctx.SendString(user.Email)
	})
	app.Get("/logout", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}

		return ctx.SendString("logout")
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./index.html")
	})

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}

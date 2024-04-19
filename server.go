package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"github.com/shareed2k/goth_fiber"
)

func main() {

	godotenv.Load()
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:8000, http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Static("/static", "./static")

	store := session.New(session.Config{
		Storage: postgres.New(postgres.Config{
			ConnectionURI: os.Getenv("POSTGRES_DATABASE_URL"),
			Database:      "godecide",
			Table:         "sessions",
			Reset:         true,
		}),
	})

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://127.0.0.1:8000/auth/callback/google"),
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), "http://127.0.0.1:8000/auth/callback/github"),
	)

	app.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	app.Get("/auth/callback/:provider", func(ctx *fiber.Ctx) error {
		CompleteUserAuthOptions := goth_fiber.CompleteUserAuthOptions{
			ShouldLogout: false,
		}
		_, err := goth_fiber.CompleteUserAuth(ctx, CompleteUserAuthOptions)
		if err != nil {
			log.Fatal(err)
		}
		sess, err := store.Get(ctx)
		if err != nil {
			return err
		}
		sess.Set("user", "bill")
		sess.Set("provider", "github")
		sess.Save()

		return ctx.Redirect("/")
	})
	app.Get("/logout", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}
		// Destroy session
		sess, err := store.Get(ctx)
		if err != nil {
			panic(err)
		}
		if err := sess.Destroy(); err != nil {
			panic(err)
		}

		return ctx.Redirect("/")
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("./index.html", nil)
	})

	app.Get("/ideas", func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		if err != nil {
			return err
		}
		return ctx.JSON(&fiber.Map{
			"page":    "ideas",
			"session": sess.ID(),
		})
	})
	app.Get("/me", func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		if err != nil {
			return err
		}
		return ctx.JSON(&fiber.Map{
			"page":     "me",
			"session":  sess.ID(),
			"provider": sess.Get("provider"),
			"keys":     sess.Keys(),
		})
	})

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}

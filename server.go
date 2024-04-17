package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres"
	"github.com/gofiber/template/django/v3"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"

	"github.com/shareed2k/goth_fiber"
)

func main() {

	godotenv.Load()
	app := fiber.New(fiber.Config{
		Views:        django.New("./templates", ".django"),
		ViewsLayout:  "main",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:8000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Static("/static", "./static")

	storage := postgres.New(postgres.Config{
		ConnectionURI: os.Getenv("POSTGRES_DATABASE_URL"),
		Database:      "postgres",
		Table:         "sessions",
		Reset:         true,
	})

	store := session.New(session.Config{
		Storage: storage,
	})

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
		log.Println(user)

		return ctx.Redirect("/")
	})
	app.Get("/logout", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}

		return ctx.SendString("logout")
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		sess.Set("name", "asa")
		sess.Set("provider", "google")
		keys := sess.Keys()
		sess.Save()
		if err != nil {
			return err
		}
		fmt.Println(time.Now().Format("15:04:05.000000"), sess.ID(), keys)
		return ctx.Render("index", nil)
	})

	app.Get("/ideas", func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		user := sess.Get("name")
		provider := sess.Get("provider")
		if err != nil {
			return err
		}
		fmt.Println(user, provider)
		return ctx.Render("ideas", nil)
	})

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}

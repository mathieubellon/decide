package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"github.com/shareed2k/goth_fiber"
)

func main() {
	godotenv.Load()
	connectDB()
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	app.Use(logger.New())
	app.Use(csrf.New())
	app.Use(MakeSession)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:8000, http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Static("/static", "./static")

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://127.0.0.1:8000/auth/callback/google"), // TODO make BASE_URL an env variable
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), "http://127.0.0.1:8000/auth/callback/github"),
	)

	app.Get("/", Homepage) // Serves vue frontend
	app.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	app.Get("/auth/callback/:provider", Callback)
	app.Get("/logout", Logout)
	app.Get("/ideas", ListIdeas)
	app.Get("/me", Me)

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}

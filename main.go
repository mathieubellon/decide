package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"github.com/shareed2k/goth_fiber"
)
const discord = "https://discord.com/api/webhooks/363115391473680395/91nFHPPElzPRPS7RIYhTbCBS5I2yzCwpMO-_z9FCGb_cK4P4e74-bwOJGbNTNiZixPal"
const lin_api_key = "lin_api_fdYscgaOYjWWB3FmjPItEpb6REPTzrKZ5SIB3VNk"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	connectDB()
	initSession()
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	app.Use(logger.New())
	app.Use(csrf.New())
	app.Use(recover.New())
	app.Use(helmet.New())
	// 404 Handler

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:8000, http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://127.0.0.1:8000/auth/callback/google"), // TODO make BASE_URL an env variable
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), "http://127.0.0.1:8000/auth/callback/github"),
	)

	app.Static("/", "./index.html")

	app.Get("/home", Homepage).Name("index") // Serves vue frontend
	app.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	app.Get("/auth/callback/:provider", Callback)
	app.Get("/logout", Logout)

	api := app.Group("/api") // /api

	// Authentication required for all routes beginning with /api
	api.Use(Protected)

	api.Get("/me", Me)
	v1 := api.Group("/v1") // /api/v1
	v1.Get("/ideas", ListIdeas)

	// 404 are handled by frontend => see ./src/router/index.js
	app.Static("*", "./index.html")

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}

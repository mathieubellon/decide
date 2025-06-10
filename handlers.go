package main

import (
	"fmt"
	"log"

	"github.com/Pallinder/go-randomdata"
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
	"github.com/supabase-community/supabase-go"
)

func Homepage(ctx *fiber.Ctx) error {
	session, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	const supabase_url = "http://localhost:54321"
	const supabase_key = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Impva3N5dm9pbHlybmV2dGZtdGRvIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDk0NTc2MTUsImV4cCI6MjA2NTAzMzYxNX0.DXfWtvMECO1VtTgMwZkWMexbYR3A9ZWvMI8cfbT7PtA"

	supabase := supabase.NewClient(supabase_url, supabase_key)

	supabase.From("ideas").Select("*").ExecuteTo(&ideas)

	return ctx.Render("./index.html", fiber.Map{"Email": session.Get("userEmail")})
}

func ListIdeas(ctx *fiber.Ctx) error {
	store, err := globalSession.Get(ctx)
	if err != nil {
		return err
	}

	wid := store.Get("workspaceID").(uint)
	uid := store.Get("userID").(uint)

	idea := Idea{
		Title:       fmt.Sprintf("%s %s %s", randomdata.SillyName(), randomdata.Noun(), randomdata.Adjective()),
		Description: randomdata.Paragraph(),
		Votes:       randomdata.Number(1000),
		UserID:      uid,
		WorkspaceID: wid,
		Reach:       randomdata.Number(5),
		Priority:    randomdata.Number(10),
	}
	if err := db.Create(&idea).Error; err != nil {
		panic(err)
	}

	var workspace Workspace

	db.Model(&Workspace{}).Preload("Ideas").Find(&workspace, store.Get("workspaceID").(uint))

	return ctx.JSON(&fiber.Map{
		"page":  "ideas",
		"idea":  idea,
		"ideas": workspace.Ideas,
	})
}

func Me(ctx *fiber.Ctx) error {
	sess, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}

	var user User

	errow := db.Model(&User{}).Preload("SocialAccounts").Preload("Workspaces").Find(&user, sess.Get("userID").(uint)).Error
	if errow != nil {
		return err
	}

	return ctx.JSON(&fiber.Map{
		"page":          "me",
		"session":       sess.ID(),
		"user":          user.Email,
		"userUUID":      user.UUID,
		"userID":        user.ID,
		"workspaceID":   user.Workspaces[0].ID,
		"workspaceName": user.Workspaces[0].Name,
		"avatarURL":     user.SocialAccounts[0].AvatarURL,
		"provider":      user.SocialAccounts[0].Provider,
		"keys":          sess.Keys(),
	})
}

func Logout(ctx *fiber.Ctx) error {
	if err := goth_fiber.Logout(ctx); err != nil {
		log.Fatal(err)
	}
	// Destroy session
	sess, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	if err := sess.Destroy(); err != nil {
		panic(err)
	}

	return ctx.Redirect("/")
}

func Callback(ctx *fiber.Ctx) error {
	oauthResponse, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		log.Fatal(err)
	}

	user, err := FindOrCreateUser(oauthResponse)
	if err != nil {
		return ctx.Status(401).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	CreateUserSession(ctx, user.ID)

	return ctx.RedirectToRoute("index", nil)
}

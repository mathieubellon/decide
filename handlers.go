package main

import (
	"fmt"
	"log"

	"github.com/Pallinder/go-randomdata"
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func Homepage(ctx *fiber.Ctx) error {
	session, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	return ctx.Render("./index.html", fiber.Map{"Email": session.Get("userEmail")})
}

func ListIdeas(ctx *fiber.Ctx) error {
	store, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}

	idea := Idea{Title: fmt.Sprintf("%s %s %s", randomdata.SillyName(), randomdata.Noun(), randomdata.Adjective()), Description: randomdata.Paragraph(), Votes: randomdata.Number(1000), UserID: store.Get("userID").(uint), WorkspaceID: store.Get("workspaceID").(uint), Reach: randomdata.Number(5), Priority: randomdata.Number(10)}
	if err := db.Create(&idea).Error; err != nil {
		log.Println(err)
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

	if sess.Fresh() {
		return ctx.JSON(&fiber.Map{
			"page": "Unauthenticated",
		})
	}
	var user User

	errow := db.Model(&User{}).Preload("SocialAccounts").Preload("Workspaces").Find(&user, sess.Get("userID").(uint)).Error
	if errow != nil {
		return err
	}
	fmt.Println(user)

	// workspace := Workspace{}
	// db.Where("email = ?", sess.Get("userEmail")).First(&user)

	// db.Model(&workspace).Preload("Users").Find(&user)

	return ctx.JSON(&fiber.Map{
		"page":          "me",
		"session":       sess.ID(),
		"user":          sess.Get("userEmail"),
		"provider":      sess.Get("provider"),
		"userUUID":      user.UUID,
		"userID":        user.ID,
		"workspaceID":   user.Workspaces[0].ID,
		"workspaceName": user.Workspaces[0].Name,
		"avatarURL":     user.SocialAccounts[0].AvatarURL,
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
	if err != nil {
		panic(err)
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
		return ctx.SendStatus(500)
	}

	sess, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	sess.Fresh()
	sess.Set("userEmail", user.Email)
	//	sess.Set("provider", user.SocialAccounts[0].Provider)
	sess.Set("userUUID", user.UUID)
	sess.Set("userID", user.ID)
	//	sess.Set("avatarURL", user.SocialAccounts[0].AvatarURL)
	sess.Save()

	return ctx.JSON(oauthResponse)
}

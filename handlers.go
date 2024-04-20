package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func Homepage(ctx *fiber.Ctx) error {
	return ctx.Render("./index.html", nil)
}

func ListIdeas(ctx *fiber.Ctx) error {
	sess, err := store.Get(ctx)
	if err != nil {
		return err
	}
	return ctx.JSON(&fiber.Map{
		"page":    "ideas",
		"session": sess.ID(),
	})
}

func Me(ctx *fiber.Ctx) error {
	sess, err := store.Get(ctx)
	if err != nil {
		return err
	}

	return ctx.JSON(&fiber.Map{
		"page":     "me",
		"session":  sess.ID(),
		"user":     sess.Get("userEmail"),
		"provider": sess.Get("provider"),
		"keys":     sess.Keys(),
		"fresh":    sess.Fresh(),
	})
}

func Logout(ctx *fiber.Ctx) error {
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
}

func Callback(ctx *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		log.Fatal(err)
	}

	appuser := new(User)
	appuser.Email = user.Email
	db.Create(appuser)

	sess, err := store.Get(ctx)
	sess.Fresh()
	sess.Set("userEmail", user.Email)
	sess.Set("provider", "github")
	sess.Save()

	return ctx.JSON(user)
}

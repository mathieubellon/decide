package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func Homepage(ctx *fiber.Ctx) error {
	session, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	log.Println(session.ID())
	log.Println(session.Get("userEmail"))
	return ctx.Render("./index.html", fiber.Map{"Email": session.Get("userEmail")})
}

func ListIdeas(ctx *fiber.Ctx) error {
	store, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return ctx.JSON(&fiber.Map{
		"page":    "ideas",
		"session": store.ID(),
	})
}

func Me(ctx *fiber.Ctx) error {
	sess, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
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
	user, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		log.Fatal(err)
	}

	appuser := new(User)
	appuser.Email = user.Email
	db.Create(appuser)

	sess, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	sess.Fresh()
	sess.Set("userEmail", user.Email)
	sess.Set("provider", "github")
	sess.Save()

	log.Println(sess)

	return ctx.JSON(user)
}

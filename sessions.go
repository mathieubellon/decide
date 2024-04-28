package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres"
)

var globalSession *session.Store

// Init sessions store
func initSession() {
	globalSession = session.New(session.Config{
		Storage: postgres.New(postgres.Config{
			ConnectionURI: os.Getenv("POSTGRES_DATABASE_URL"),
			Database:      "godecide",
			Table:         "sessions",
		}),
		Expiration: 1 * time.Hour,
		KeyLookup:  "cookie:godecide_session",
	})
}

// Todo: Implement Session as struvct for validation
// type Session struct {
// 	SID    string
// 	UUID   string
// 	IP     string
// 	Login  string
// 	UA     string
// 	Expire time.Time
// }

func CreateUserSession(c *fiber.Ctx, uid uint) error {
	// Get or create session
	s, _ := globalSession.Get(c)

	// If this is a new session
	var user User

	err := db.Model(&User{}).Preload("SocialAccounts").Preload("Workspaces").Find(&user, uid).Error
	if err != nil {
		panic(err)
	}

	// Get session ID
	sid := s.ID()

	// Save session data
	s.Set("userID", uid)
	s.Set("workspaceID", user.Workspaces[0].ID)
	s.Set("sid", sid)
	s.Set("ip", c.Context().RemoteIP().String())
	s.Set("login", time.Unix(time.Now().Unix(), 0).UTC().String())
	s.Set("ua", string(c.Request().Header.UserAgent()))

	err = s.Save()
	if err != nil {
		// log.Println(err)
		return err
	}

	return nil
}

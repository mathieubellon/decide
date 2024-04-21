package main

import (
	"github.com/google/uuid"
	"github.com/markbates/goth"
)

func FindOrCreateUser(oauthResponse goth.User) (*User, error) {
	user := &User{}
	if err := db.Where("email = ?", oauthResponse.Email).First(user).Error; err != nil {
		user.Email = oauthResponse.Email
		user.UUID = uuid.Must(uuid.NewRandom()).String()
		if err := db.Create(user).Error; err != nil {
			return nil, err
		}
	}
	return user, nil
}

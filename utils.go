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
		user.SocialAccounts = append(user.SocialAccounts, SocialAccount{
			Provider:  oauthResponse.Provider,
			UID:       oauthResponse.UserID,
			Firstname: oauthResponse.FirstName,
			Lastname:  oauthResponse.LastName,
			Nickname:  oauthResponse.NickName,
			AvatarURL: oauthResponse.AvatarURL,
		})
		// TODO if email already exists add social account to user
		// TODO Force email presence (Github does not provide email)

		if err := db.Create(user).Error; err != nil {
			return nil, err
		}
	}
	return user, nil
}

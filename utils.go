package main

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/markbates/goth"
)

func FindOrCreateUser(oauthResponse goth.User) (*User, error) {
	user := User{}
	workspace := Workspace{}
	if err := db.Where("email = ?", oauthResponse.Email).First(&user).Error; err != nil {
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
		if err := db.Create(&user).Error; err != nil {
			return nil, err
		}
		workspace.Name = randomdata.SillyName()
		workspace.UUID = uuid.Must(uuid.NewRandom()).String()
		if err := db.Create(&workspace).Error; err != nil {
			return nil, err
		}
		db.Model(&workspace).Association("Users").Append(&user)
	}

	// TODO if email already exists add social account to user
	// TODO Force email presence (Github does not provide email in 100% of cases)

	return &user, nil
}

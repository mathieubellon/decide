package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model     `json:"-"`
	Email          string          `json:"email"`
	UUID           string          `json:"uuid" gorm:"unique;not null; index"`
	UserSessions   []UserSession   `json:"user_sessions,omitempty"`
	SocialAccounts []SocialAccount `json:"social_accounts,omitempty"`
	WorkspaceID    uint            `json:"workspace_id"`
}

type UserSession struct {
	gorm.Model `json:"-"`
	SID        string `json:"sid"`
	IP         string `json:"ip"`
	Login      string `json:"login"`
	Expiry     string `json:"expiry"`
	UA         string `json:"ua"`
	UserID     string `json:"user_id"`
}

type Workspace struct {
	gorm.Model `json:"-"`
	Name       string `json:"name"`
	Users      []User `json:"users,omitempty" gorm:"many2many:workspace_users;"`
	UUID       string `json:"uuid" gorm:"unique;not null; index"`
}

type SocialAccount struct {
	gorm.Model `json:"-"`
	Provider   string `json:"provider"`
	UID        string `json:"uid"`
	UserID     uint   `json:"user_id"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Nickname   string `json:"nickname"`
	AvatarURL  string `json:"avatar_url"`
}

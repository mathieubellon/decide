package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Password  lin_api_fdYscgaOYjWWB3FmjPItEpb6REPTzrKZ5SIB3VNl `gorm:"not null"`
}

type User struct {
	gorm.Model `json:"-"`
	//ID             uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primary_key;" json:"id"`
	Email          string          `json:"email"`
	UUID           string          `json:"uuid" gorm:"unique;not null; index;default:null"`
	UserSessions   []UserSession   `json:"user_sessions,omitempty"`
	SocialAccounts []SocialAccount `json:"social_accounts,omitempty"`
	Ideas          []Idea          `json:"user_ideas,omitempty"`
	Workspaces     []Workspace     `json:"workspaces,omitempty" gorm:"many2many:workspace_users;"`
}

type Idea struct {
	gorm.Model  `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Votes       int    `json:"votes"`
	UserID      uint   `json:"user_id"`
	WorkspaceID uint   `json:"workspace_id"`
	Reach       int    `json:"reach"`
	Priority    int    `json:"priority"`
	Calculated  int    `json:"calculated" gorm:"-"`
}

type Formula struct {
	gorm.Model `json:"-"`
	Content    string `json:"content"`
}

type Workspace struct {
	gorm.Model `json:"-"`
	Name       string `json:"name"`
	Users      []User `json:"users,omitempty" gorm:"many2many:workspace_users;"`
	UUID       string `json:"uuid" gorm:"unique;not null; index;default:null"`
	Ideas      []Idea `json:"workspace_ideas,omitempty"`
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

type UserSession struct {
	gorm.Model `json:"-"`
	SID        string `json:"sid"`
	IP         string `json:"ip"`
	Login      string `json:"login"`
	Expiry     string `json:"expiry"`
	UA         string `json:"ua"`
	UserID     string `json:"user_id"`
}

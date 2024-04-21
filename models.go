package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model   `json:"-"`
	Email        string        `json:"email"`
	UserSessions []UserSession `json:"user_sessions,omitempty"`
	UUID         string        `json:"uuid" gorm:"unique;not null; index"`
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

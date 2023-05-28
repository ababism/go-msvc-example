package core

import (
	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID
	Gmail        string
	Username     string
	Nickname     string
	IsRegistered bool
	Role         string
}

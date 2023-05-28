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

func (u *User) ToDTO() (user UserDTO) {
	user.Id = u.Id
	user.Gmail = u.Gmail
	user.Nickname = u.Nickname
	user.Username = u.Username
	user.IsRegistered = u.IsRegistered
	user.Role = u.Role
	return
}

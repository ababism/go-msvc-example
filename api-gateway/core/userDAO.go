package core

import (
	"database/sql"
	"github.com/google/uuid"
)

type UserDAO struct {
	Id           uuid.UUID `db:"id"`
	Gmail        string    `db:"gmail"`
	Username     string    `db:"username"`
	Nickname     string    `db:"nickname"`
	IsRegistered bool      `db:"is_registered"`
	Role         string    `db:"role"`
}
type UserNullableDAO struct {
	Id           uuid.UUID      `db:"id"`
	Gmail        string         `db:"gmail"`
	Username     sql.NullString `db:"username"`
	Nickname     sql.NullString `db:"nickname"`
	IsRegistered bool           `db:"is_registered"`
	Role         string         `db:"role"`
}

func (u *UserDAO) ToDomain() (user User) {
	user.Id = u.Id
	user.Gmail = u.Gmail
	user.Username = u.Username
	user.Nickname = u.Nickname
	user.IsRegistered = u.IsRegistered
	user.Role = u.Role
	return
}
func (u *UserNullableDAO) ToDomain() (user User) {
	user.Id = u.Id
	user.Gmail = u.Gmail
	if u.Username.Valid {
		user.Username = u.Username.String
	}
	if u.Nickname.Valid {
		user.Nickname = u.Nickname.String
	}
	user.IsRegistered = u.IsRegistered
	user.Role = u.Role
	return
}

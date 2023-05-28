package service

import (
	"github.com/google/uuid"
	"user-svc/internal/core"
	"user-svc/internal/repository"
)

type Token interface {
	GetJWT(gmail string) (core.JWT, error)
	ParseToken(token string) (uuid.UUID, bool, string, error)
	GenerateJWT(userId uuid.UUID, registered bool, role string) (core.JWT, error)
}

type User interface {
	RegisterUser(userId uuid.UUID, user core.User) (core.User, error)
	ChangeUsername(clientId uuid.UUID, username string) (core.User, error)
	ChangeNickname(clientId uuid.UUID, nickname string) (core.User, error)
	SearchUsers(query string, clientId uuid.UUID, limit int, offset int) ([]core.User, error)
	Exists(id uuid.UUID) bool
	GetById(id uuid.UUID) (core.User, error)
	AddPlayerID(clientId uuid.UUID, playerID uuid.UUID) error
	GetPlayerID(userId uuid.UUID) (string, error)
	InitPlayerID() error
}

type Service struct {
	Token
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Token: NewTokenService(repos.User),
		User:  NewUserService(repos.User),
	}
}

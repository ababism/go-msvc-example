package service

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"user-svc/internal/core"
	"user-svc/internal/repository"
)

const UsernameSymbols = "qwertyuiopasdfghjklzxcvbnm_-.0123456789"
const NicknameSymbols = " qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM_-.0123456789"

type UserService struct {
	r repository.User
}

func (s *UserService) InitPlayerID() error {
	return s.r.InitPlayerID()
}

func (s *UserService) GetPlayerID(userId uuid.UUID) (string, error) {
	res, err := s.r.GetPlayerID(userId)
	if err != nil {
		return "", core.ErrNotFound
	}
	return res, nil
}

func (s *UserService) AddPlayerID(clientId uuid.UUID, playerID uuid.UUID) error {
	if !s.r.ExistsWithId(clientId) {
		return core.ErrNotFound
	}
	err := s.r.InstallAppID(clientId, playerID)
	return err
}

func (s *UserService) GetById(id uuid.UUID) (core.UserPayload, error) {
	user, err := s.r.GetById(id)
	return user.ToPayload(), err
}

func (s *UserService) Exists(id uuid.UUID) bool {
	return s.r.ExistsWithId(id)
}

func (s *UserService) SearchUsers(query string, clientId uuid.UUID, limit int, offset int) (res []core.UserPayload, err error) {
	users, err := s.r.SearchUsers(query, clientId, limit, offset)
	if err != nil {
		return
	}
	res = make([]core.UserPayload, len(users))
	for i, user := range users {
		uDomain := user.ToDomain()
		res[i] = uDomain.ToPayload()
	}
	return res, nil
}

func (s *UserService) ChangeNickname(clientId uuid.UUID, nickname string) (core.User, error) {
	if !s.validateNickname(nickname) {
		return core.User{}, errors.New("username is invalid")
	}
	newClient, err := s.r.ChangeNickname(clientId, nickname)
	if err != nil {
		return core.User{}, core.ErrInternal
	}
	return newClient.ToDomain(), nil
}

func (s *UserService) ChangeUsername(clientId uuid.UUID, username string) (core.User, error) {
	tmpUser, err := s.r.GetByUsername(username)
	if err == nil && tmpUser.IsRegistered {
		return core.User{}, errors.New("user with this nickname already exists")
	}
	if !s.validateUsername(username) {
		return core.User{}, errors.New("username is invalid")
	}
	newClient, err := s.r.ChangeUsername(clientId, username)
	if err != nil {
		return core.User{}, core.ErrInternal
	}
	return newClient.ToDomain(), nil
}

func NewUserService(repository repository.User) *UserService {
	return &UserService{r: repository}
}
func (s *UserService) RegisterUser(id uuid.UUID, user core.User) (core.User, error) {
	if ok, err := s.validateUserFields(user); err != nil {
		return core.User{}, err
	} else if !ok {
		return core.User{}, err
	}

	user.Id = id
	u, err := s.r.Register(user)
	return u.ToDomain(), err
}

func (s *UserService) validateUserFields(user core.User) (bool, error) {
	_, err := s.r.GetByUsername(user.Username)
	if err == nil {
		return false, errors.New("username already exists")
	}
	if !s.validateUsername(user.Username) {
		return false, errors.New("invalid username")
	}
	if !s.validateNickname(user.Nickname) {
		return false, errors.New("invalid nickname")
	}
	if !s.validateRole(user.Role) {
		return false, errors.New("invalid role")
	}
	return true, nil
}

func (s *UserService) validateUsername(username string) bool {
	if len(username) < 5 || len(username) > 30 {
		return false
	}
	for _, s := range username {
		if !strings.ContainsRune(UsernameSymbols, s) {
			return false
		}
	}
	return true
}
func (s *UserService) validateNickname(nickname string) bool {
	if len(nickname) < 5 || len(nickname) > 30 {
		return false
	}
	for _, s := range nickname {
		if !strings.ContainsRune(NicknameSymbols, s) {
			return false
		}
	}
	return true
}
func (s *UserService) validateRole(role string) bool {
	if role != core.ManagerRole && role != core.DefaultRole {
		return false
	}
	return true
}

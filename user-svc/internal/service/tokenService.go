package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"time"
	"user-svc/internal/core"
	"user-svc/internal/repository"
)

// TokenService -- сервис для работы с jwt, он привязан к репозиторию пользователей для проверок
type TokenService struct {
	userRepo repository.User
}

var (
	signingKey = ""
)

const (
	registeredTTLYears     = 1
	unregisteredTTLMinutes = 15
)

// TODO

type tokenClaims struct {
	jwt.StandardClaims
	UserId       uuid.UUID `json:"user_id"`
	IsRegistered bool      `json:"is_registered"`
	Role         string    `json:"role"`
}

func initializeSecret() {
	signingKey = viper.GetString("signing_key")
}

func NewTokenService(userRepo repository.User) *TokenService {
	initializeSecret()
	return &TokenService{userRepo: userRepo}
}

// GetJWT дает jwt, в обмен на почту пользователя.
func (s *TokenService) GetJWT(gmail string) (core.JWT, error) {
	if !s.userRepo.Exists(gmail) {
		userId, err := s.userRepo.Create(gmail)
		if err != nil {
			return core.JWT{}, err
		}
		return s.GenerateJWT(userId, false, core.DefaultRole)
	}
	user, err := s.userRepo.GetByGmail(gmail)
	if err != nil {
		return core.JWT{}, err
	}
	return s.GenerateJWT(user.Id, user.IsRegistered, user.Role)
}

// GenerateJWT генерирует jwt, из нужных параметров
func (s *TokenService) GenerateJWT(userId uuid.UUID, registered bool, role string) (core.JWT, error) {
	if registered {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().AddDate(registeredTTLYears, 0, 0).Unix(),
			},
			userId,
			true,
			role,
		})
		str, err := token.SignedString([]byte(signingKey))
		return core.JWT{Token: str, IsRegistered: registered}, err
	} else {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * unregisteredTTLMinutes).Unix(),
			},
			userId,
			false,
			role,
		})
		str, err := token.SignedString([]byte(signingKey))
		return core.JWT{Token: str, IsRegistered: registered}, err
	}
}

// ParseToken проверяет jwt, в случае успешной проверки отдает актульную информацию о пользователе
func (s *TokenService) ParseToken(token string) (uuid.UUID, bool, string, error) {
	t, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.UUID{0}, false, core.DefaultRole, err
	}

	claims, ok := t.Claims.(*tokenClaims)
	if !ok {
		return uuid.UUID{0}, false, core.DefaultRole, errors.New("token claims are not of type *tokenClaims")
	}
	usr, err := s.userRepo.GetById(claims.UserId)
	if claims.IsRegistered != usr.IsRegistered {
		if err != nil {
			return uuid.UUID{0}, false, core.DefaultRole, err
		}
		return uuid.UUID{0}, false, core.DefaultRole, errors.New("registration flags with token and user do not match")
	}
	return claims.UserId, claims.IsRegistered, claims.Role, nil
}

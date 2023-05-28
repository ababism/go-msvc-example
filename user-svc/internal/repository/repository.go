package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"user-svc/internal/core"
)

//type ReleaseItem interface {
//	GetSongById(song uuid.UUID) (core.SongDAO, error)
//	GetAlbumById(album uuid.UUID) (core.AlbumDAO, error)
//}

type User interface {
	// Create returns id of new user, and changes his id
	Create(gmail string) (uuid.UUID, error)
	Exists(gmail string) bool
	GetById(userId uuid.UUID) (user core.User, err error)
	GetByUsername(username string) (core.User, error)
	GetByGmail(gmail string) (user core.User, err error)
	Register(u core.User) (user core.User, err error)
	ChangeUsername(id uuid.UUID, username string) (user core.User, err error)
	ChangeNickname(id uuid.UUID, nickname string) (user core.User, err error)
	SearchUsers(query string, clientId uuid.UUID, limit int, offset int) ([]core.User, error)
	ExistsWithId(id uuid.UUID) bool
	InstallAppID(clientId uuid.UUID, playerID uuid.UUID) error
	GetPlayerID(userId uuid.UUID) (string, error)
	InitPlayerID() error
}

type Repository struct {
	//ReleaseItem
	User
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}

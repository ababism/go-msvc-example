package repository

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"user-svc/internal/core"
	"user-svc/internal/repository/dao"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) InitPlayerID() (err error) {
	q := `
	CREATE TABLE user_appid
	(
    user_id UUID        not null PRIMARY KEY,
    app_id  UUID UNIQUE not null,
    constraint user_fk foreign key (user_id) references users (id)
	);
	`
	_, err = r.db.Exec(q)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (r UserRepository) GetPlayerID(userId uuid.UUID) (res string, err error) {
	q := `
	SELECT app_id FROM user_appid WHERE user_id = $1
	`
	err = r.db.Get(&res, q, userId)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (r UserRepository) InstallAppID(clientId uuid.UUID, playerID uuid.UUID) error {
	q := `
	INSERT INTO user_appid VALUES 
	 ($1, $2)
	`
	_, err := r.db.Exec(q, clientId, playerID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (r UserRepository) ExistsWithId(id uuid.UUID) (res bool) {
	q := `
	SELECT EXISTS(SELECT
	FROM users
	WHERE id = $1);
	`
	logrus.Trace(formatQuery(q))
	err := r.db.Get(&res, q, id)
	if err != nil {
		return false
	}
	return res
}

func (r UserRepository) SearchUsers(query string, clientId uuid.UUID, limit int, offset int) (res []core.User, err error) {
	var users []dao.UserDAO
	q := `
	SELECT *
		FROM users
	WHERE users.username ILIKE $1 || '%'
	ORDER BY username
	LIMIT $2 OFFSET $3;
	`
	logrus.Trace(formatQuery(q))
	err = r.db.Select(&users, q, query, limit, offset)
	res = make([]core.User, len(users))
	for i, v := range users {
		res[i] = v.ToDomain()
	}
	return
}

func (r UserRepository) Create(gmail string) (uuid.UUID, error) {
	q := `
	INSERT INTO users (gmail, is_registered)
	VALUES ($1, false)
	RETURNING id
	`
	logrus.Trace(formatQuery(q))
	row := r.db.QueryRow(q, gmail)
	var (
		id  uuid.UUID
		err error
	)
	err = row.Scan(&id)
	return id, err
}

func (r UserRepository) Exists(gmail string) (res bool) {
	q := `
	SELECT EXISTS(SELECT
	FROM users
	WHERE gmail = $1);
	`
	logrus.Trace(formatQuery(q))
	err := r.db.Get(&res, q, gmail)
	if err != nil {
		return false
	}
	return res
}

func (r UserRepository) GetById(userId uuid.UUID) (core.User, error) {
	q := `
	SELECT * FROM users WHERE id = $1
	`
	logrus.Trace(formatQuery(q))
	var user dao.UserNullableDAO
	err := r.db.Get(&user, q, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.User{}, core.ErrNotFound
		}
		return core.User{}, err
	}
	return user.ToDomain(), nil
}

func (r UserRepository) GetByUsername(username string) (res core.User, err error) {
	var user dao.UserDAO
	q := `
	SELECT * FROM users WHERE username = $1
	`
	logrus.Trace(formatQuery(q))
	err = r.db.Get(&user, q, username)
	if err != nil {
		return core.User{}, err
	}
	return user.ToDomain(), nil
}

func (r UserRepository) GetByGmail(gmail string) (core.User, error) {
	q := `
	SELECT * FROM users WHERE gmail = $1
	`
	logrus.Trace(formatQuery(q))
	var user dao.UserNullableDAO
	err := r.db.Get(&user, q, gmail)
	//row := r.db.QueryRow(q, gmail)
	//err = row.Scan(&user)
	if err != nil {
		return core.User{}, err
	}
	return user.ToDomain(), nil
}

func (r UserRepository) Register(u core.User) (res core.User, err error) {
	var user dao.UserDAO
	q := `
	UPDATE users
	SET (username, nickname, is_registered, role) = ($1, $2, true, $4)
	WHERE id = $3
	RETURNING *
	`
	logrus.Trace(formatQuery(q))
	//row := r.db.QueryRow(q, u.Username, u.Nickname, u.Id)
	//err = row.Scan(&user)
	err = r.db.Get(&user, q, u.Username, u.Nickname, u.Id, u.Role)
	res = user.ToDomain()
	return
}

func (r UserRepository) ChangeUsername(userId uuid.UUID, username string) (res core.User, err error) {
	var user dao.UserDAO
	q := `
	UPDATE users
	SET username = $1
	WHERE id = $2
	RETURNING *
	`
	logrus.Trace(formatQuery(q))
	err = r.db.Get(&user, q, username, userId)
	res = user.ToDomain()
	return
}
func (r UserRepository) ChangeNickname(userId uuid.UUID, nickname string) (res core.User, err error) {
	var user dao.UserDAO
	q := `
	UPDATE users
	SET nickname = $1
	WHERE id = $2
	RETURNING *
	`
	logrus.Trace(formatQuery(q))
	err = r.db.Get(&user, q, nickname, userId)
	res = user.ToDomain()
	return
}

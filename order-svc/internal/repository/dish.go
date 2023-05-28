package repository

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"order-svc/internal/core"
	"order-svc/internal/repository/dao"
)

const (
	nameMaxSize int = 100
)

func NewDishRepository(db *sqlx.DB) *DishRepository {
	return &DishRepository{db: db}
}

type DishRepository struct {
	db *sqlx.DB
}

func (dr DishRepository) GetAll(limit, offset int64) ([]core.Dish, error) {
	var (
		dishes []dao.DishDAO
	)
	q := `
	SELECT *
	FROM dishes
	ORDER BY dishes.name
	LIMIT $1 OFFSET $2;
	`
	logrus.Trace(formatQuery(q))
	err := dr.db.Select(&dishes, q, limit, offset)
	if err != nil {
		logrus.Error(err)
		return []core.Dish{}, core.ErrInternal
	}
	res := make([]core.Dish, len(dishes))
	for i, dish := range dishes {
		res[i] = dish.ToDomain()
	}
	return res, nil
}

func (dr DishRepository) Create(d core.Dish) (core.Dish, error) {
	var dish dao.DishDAO
	if len(d.Name) > 50 || d.Quantity < 0 || d.Price < 0 {
		return core.Dish{}, core.ErrIncorrectBody
	}
	q := `
	INSERT INTO  dishes ("name", description, price, quantity, is_available)
	VALUES ($1, $2, $3, $4, $5) RETURNING *;
	`
	logrus.Trace(formatQuery(q))
	err := dr.db.Get(&dish, q, d.Name, d.Description, d.Price, d.Quantity, d.IsAvailable)
	if err != nil {
		logrus.Error(err)
		if err.Error() == "pq: duplicate key value violates unique constraint \"name_unique\"" {
			return core.Dish{}, core.ErrAlreadyExists
		}
		return core.Dish{}, core.ErrInternal
	}
	return dish.ToDomain(), nil
}

func (dr DishRepository) Get(id uuid.UUID) (core.Dish, error) {
	var dish dao.DishDAO
	q := `
	SELECT * FROM dishes WHERE id = $1;
	`
	logrus.Trace(formatQuery(q))
	err := dr.db.Get(&dish, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Dish{}, core.ErrNotFound
		}
		logrus.Error(err)
		return core.Dish{}, core.ErrInternal
	}
	return dish.ToDomain(), nil
}

func (dr DishRepository) Update(d core.Dish) (core.Dish, error) {
	var dish dao.DishDAO
	if len(d.Name) > nameMaxSize || d.Quantity < 0 || d.Price < 0 {
		return core.Dish{}, core.ErrIncorrectBody
	}
	q := `
	UPDATE dishes SET ("name", description, price, quantity, is_available) = 
	 ($1, $2, $3, $4, $5) WHERE id = $6 RETURNING *;
	`
	logrus.Trace(formatQuery(q))
	err := dr.db.Get(&dish, q, d.Name, d.Description, d.Price, d.Quantity, d.IsAvailable, d.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Dish{}, core.ErrNotFound
		}
		logrus.Error(err)
		return core.Dish{}, core.ErrInternal
	}
	return dish.ToDomain(), nil
}

func (dr DishRepository) Delete(id uuid.UUID) (core.Dish, error) {
	var dish dao.DishDAO
	q := `
	UPDATE dishes
		SET is_available = false
			WHERE id = $1 RETURNING *;
	`
	logrus.Trace(formatQuery(q))
	err := dr.db.Get(&dish, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Dish{}, core.ErrNotFound
		}
		logrus.Error(err)
		return core.Dish{}, core.ErrInternal
	}
	return dish.ToDomain(), nil
}

package service

import (
	"errors"
	"github.com/google/uuid"
	"order-svc/internal/core"
	"order-svc/internal/repository"
)

func NewDishService(dr repository.Dish) *DishService {
	return &DishService{r: dr}
}

type DishService struct {
	r repository.Dish
}

func (ds DishService) GetAll(limit, offset int64) ([]core.Dish, error) {
	if limit > 200 {
		return []core.Dish{}, errors.New("incorrect limit (<= 200) value")
	}
	return ds.r.GetAll(limit, offset)
}

func (ds DishService) Create(d core.Dish) (core.Dish, error) {
	if len(d.Name) > 100 {
		return core.Dish{}, errors.New("name should be less than 100 symbols")
	}
	if d.Price < 0 || d.Quantity < 0 {
		return core.Dish{}, errors.New("quantity and price should be >= 0")
	}
	return ds.r.Create(d)
}

func (ds DishService) Get(id uuid.UUID) (core.Dish, error) {
	return ds.r.Get(id)
}

func (ds DishService) Update(d core.Dish) (core.Dish, error) {
	if len(d.Name) > 100 {
		return core.Dish{}, errors.New("name should be less than 100 symbols")
	}
	if d.Price < 0 || d.Quantity < 0 {
		return core.Dish{}, errors.New("quantity and price should be >= 0")
	}
	return ds.r.Update(d)
}

func (ds DishService) Delete(id uuid.UUID) (core.Dish, error) {
	return ds.r.Delete(id)
}

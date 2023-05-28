package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"order-svc/internal/core"
	"order-svc/internal/repository/dao"
	"time"
)

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type OrderRepository struct {
	db *sqlx.DB
}

func (or OrderRepository) Create(order core.Order) (core.Order, error) {
	var (
		t     time.Time
		currT time.Time
	)
	// TODO написать присвоение выданного id заказу
	// считаем время, когда будет готов заказ
	q := `
	SELECT ready_at FROM orders ORDER BY ready_at DESC LIMIT 1 OFFSET 0;
	`
	logrus.Trace(formatQuery(q))
	row := or.db.QueryRow(q)
	err := row.Scan(&t)
	if err != nil {
		q = `
		SELECT current_timestamp;
		`
		logrus.Trace(formatQuery(q))
		row = or.db.QueryRow(q)
		err = row.Scan(&currT)
		t = currT.Add(time.Minute * 3)
	} else {
		t = t.Add(time.Minute * 3)
	}
	// записываем заказ как отмененный
	if len(order.SpecialRequests) > 2000 {
		return core.Order{}, core.ErrIncorrectBody
	}
	// TODO написать присвоение выданного id заказу
	q = `
	INSERT INTO  orders (user_id, special_requests, ready_at)
	VALUES ($1, $2, $3) RETURNING  id, user_id, status, special_requests, ready_at, created_at, updated_at;
	`
	logrus.Trace(formatQuery(q))
	//_, err = or.db.Exec(q, order.UserId, order.SpecialRequests, t)
	row = or.db.QueryRow(q, order.UserId, order.SpecialRequests, t)
	var o dao.OrderDAO
	err = row.Scan(&o.Id, &o.UserId, &o.Status, &o.SpecialRequests, &o.ReadyAt, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		logrus.Error(err)
		return core.Order{}, core.ErrInternal
	}
	temp := order.Dishes
	order = o.ToDomain()
	order.Dishes = temp
	// смотрим в транзакции можем ли мы сделать заказ по кол-во dish, если не можем то выходим
	tx, err := or.db.Begin()
	if err != nil {
		return core.Order{}, core.ErrInternal
	}
	for _, d := range order.Dishes {
		var ok bool
		q = `
		SELECT EXISTS(SELECT
		FROM dishes
		WHERE id = $1 and quantity >= $2 and is_available);
		`
		logrus.Trace(formatQuery(q))
		row = tx.QueryRow(q, d.Id, d.Quantity)
		err = row.Scan(&ok)
		if err != nil {
			_ = tx.Rollback()
			logrus.Error(err)
			order.Status = core.OrderStatusRejected
			return order, core.ErrNotFound
		}
		if !ok {
			_ = tx.Rollback()
			order.Status = core.OrderStatusRejected
			return order, core.ErrNotFound
		}
	}
	// (если можем, то) списываем продукты и добавляем их в таблцу order_dish
	for i, d := range order.Dishes {
		var dish dao.DishDAO
		q = `
		UPDATE dishes SET quantity = quantity - $2 
		              WHERE id = $1 
		              RETURNING id, "name", description, price, quantity, is_available, created_at, updated_at;
		`
		logrus.Trace(formatQuery(q))
		row = tx.QueryRow(q, d.Id, d.Quantity)
		err = row.Scan(&dish.Id, &dish.Name, &dish.Description, &dish.Price,
			&dish.Quantity, &dish.IsAvailable, &dish.CreatedAt, &dish.UpdatedAt)

		if err != nil {
			_ = tx.Rollback()
			logrus.Error(err)
			order.Status = core.OrderStatusRejected
			return order, core.ErrInternal
		}
		q = `
		INSERT INTO order_dish (order_id, dish_id, quantity, price) 
						VALUES ($1, $2, $3, $4);
		`
		logrus.Trace(formatQuery(q))
		_, err = tx.Exec(q, order.Id, d.Id, d.Quantity, dish.Price)
		if err != nil {
			_ = tx.Rollback()
			logrus.Error(err)
			order.Status = core.OrderStatusRejected
			return order, core.ErrInternal
		}
		// Присваиваем вытянутую из меню на момент заказа блюдо
		qTemp := d.Quantity
		order.Dishes[i] = dish.ToDomain()
		order.Dishes[i].Quantity = qTemp
	}
	// обновляем OrderStatusWait и отправляем пользователю заказ
	var orderDAO dao.OrderDAO
	q = `
	UPDATE orders 
   	SET status = $1
		WHERE id = $2 RETURNING id, ready_at, created_at, updated_at;
	`
	logrus.Trace(formatQuery(q))
	row = tx.QueryRow(q, core.OrderStatusWait, order.Id)
	err = row.Scan(&orderDAO.Id, &orderDAO.ReadyAt, &orderDAO.CreatedAt, &orderDAO.UpdatedAt)
	if err != nil {
		errRollback := tx.Rollback()
		logrus.Error(errRollback)
		order.Status = core.OrderStatusRejected
		return order, core.ErrInternal
	}
	order.Status = core.OrderStatusWait
	order.UpdatedAt = orderDAO.UpdatedAt
	return order, tx.Commit()
}

func (or OrderRepository) Get(id uuid.UUID) (core.Order, error) {
	var o dao.OrderDAO
	q := `
	UPDATE orders 
   	SET status = $1
		WHERE id = $2 AND ready_at <= current_timestamp;
	`
	logrus.Trace(formatQuery(q))
	//ordTemp := dao.OrderDAO{Id: id, Status: core.OrderStatusDone}
	_, err := or.db.Exec(q, core.OrderStatusDone, id)
	if err != nil {
		logrus.Error(err)
		return core.Order{}, core.ErrInternal
	}
	q = `
	SELECT * FROM orders 
   	WHERE id = $1;
	`
	logrus.Trace(formatQuery(q))
	err = or.db.Get(&o, q, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return core.Order{}, core.ErrNotFound
		}
		logrus.Error(err)
		return core.Order{}, core.ErrInternal
	}
	var (
		dishes []dao.DishDAO
	)
	q = `
	SELECT dishes.id,
       dishes.name,
       dishes.description,
       dishes.price,
       dishes.quantity,
       dishes.is_available,
       dishes.created_at,
       dishes.updated_at,
       order_dish.quantity,
       order_dish.price
	FROM dishes
         INNER JOIN order_dish on dishes.id = order_dish.dish_id and order_dish.order_id = $1
	`
	logrus.Trace(formatQuery(q))
	err = or.db.Select(&dishes, q, id)
	if err != nil {
		logrus.Error(err)
		return core.Order{}, core.ErrInternal
	}
	order := o.ToDomain()
	order.Dishes = make([]core.Dish, len(dishes))
	for i, dish := range dishes {
		order.Dishes[i] = dish.ToDomain()
	}
	return order, nil
}

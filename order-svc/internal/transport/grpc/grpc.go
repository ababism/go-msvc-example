package grpc_transport

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"order-svc/internal/core"
	"order-svc/internal/service"
	"order-svc/pkg/pb"
)

type Handler struct {
	services *service.Service
	pb.UnimplementedOrderServiceServer
}

func (h Handler) GetAllDishes(ctx context.Context, r *pb.GetAllDishRequest) (*pb.GetAllDishResponse, error) {
	if r == nil {
		return &pb.GetAllDishResponse{
			Status: http.StatusBadRequest,
			Error:  "empty request",
		}, nil
	}
	d, err := h.services.Dish.GetAll(r.Limit, r.Offset)
	if err != nil {
		if errors.Is(err, core.ErrInternal) {
			return &pb.GetAllDishResponse{
				Status: http.StatusInternalServerError,
				Error:  core.ErrInternal.Error(),
			}, nil
		}
		return &pb.GetAllDishResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}
	res := pb.GetAllDishResponse{
		Status: http.StatusOK,
	}
	res.Dishes = make([]*pb.Dish, len(d))
	for i, cd := range d {
		res.Dishes[i] = &pb.Dish{
			Id:          cd.Id.String(),
			Name:        cd.Name,
			Description: cd.Description,
			Price:       cd.Price,
			Quantity:    cd.Quantity,
			IsAvailable: cd.IsAvailable,
			CreatedAt:   timestamppb.New(cd.CreatedAt),
			UpdatedAt:   timestamppb.New(cd.UpdatedAt),
		}
	}
	return &res, nil
}

func (h Handler) CreateDish(ctx context.Context, r *pb.CreateDishRequest) (*pb.CreateDishResponse, error) {
	if r == nil {
		return &pb.CreateDishResponse{
			Status: http.StatusBadRequest,
			Error:  "empty request",
		}, nil
	}
	rd := core.Dish{
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Quantity:    r.Quantity,
		IsAvailable: r.IsAvailable,
	}
	d, err := h.services.Dish.Create(rd)
	if err != nil {
		if errors.Is(err, core.ErrInternal) {
			return &pb.CreateDishResponse{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}, nil
		}
		return &pb.CreateDishResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}
	return &pb.CreateDishResponse{
		Status:      http.StatusOK,
		Id:          d.Id.String(),
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Quantity:    d.Quantity,
		IsAvailable: d.IsAvailable,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}, nil
}

func (h Handler) GetDish(ctx context.Context, r *pb.GetDishRequest) (*pb.GetDishResponse, error) {
	if r == nil {
		return &pb.GetDishResponse{
			Status: http.StatusBadRequest,
			Error:  "empty request",
		}, nil
	}
	id, err := uuid.Parse(r.Id)
	if err != nil {
		return &pb.GetDishResponse{
			Status: http.StatusBadRequest,
			Error:  "couldn't parse uuid",
		}, nil
	}
	d, err := h.services.Dish.Get(id)
	if err != nil {
		return &pb.GetDishResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}
	return &pb.GetDishResponse{
		Status:      http.StatusOK,
		Id:          d.Id.String(),
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Quantity:    d.Quantity,
		IsAvailable: d.IsAvailable,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}, nil
}

func (h Handler) UpdateDish(ctx context.Context, r *pb.UpdateDishRequest) (*pb.UpdateDishResponse, error) {
	if r == nil {
		return &pb.UpdateDishResponse{
			Status: http.StatusBadRequest,
			Error:  "empty request",
		}, nil
	}
	id, err := uuid.Parse(r.Id)
	if err != nil {
		return &pb.UpdateDishResponse{
			Status: http.StatusBadRequest,
			Error:  core.ErrIncorrectBody.Error(),
		}, nil
	}
	rd := core.Dish{
		Id:          id,
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Quantity:    r.Quantity,
		IsAvailable: r.IsAvailable,
	}
	d, err := h.services.Dish.Update(rd)
	if err != nil {
		if errors.Is(err, core.ErrInternal) {
			return &pb.UpdateDishResponse{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}, nil
		}
		return &pb.UpdateDishResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}
	return &pb.UpdateDishResponse{
		Status:      http.StatusOK,
		Id:          d.Id.String(),
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Quantity:    d.Quantity,
		IsAvailable: d.IsAvailable,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}, nil
}

func (h Handler) DeleteDish(ctx context.Context, r *pb.DeleteDishRequest) (*pb.DeleteDishResponse, error) {
	if r == nil {
		return &pb.DeleteDishResponse{
			Status: http.StatusBadRequest,
			Error:  "empty request",
		}, nil
	}
	id, err := uuid.Parse(r.Id)
	if err != nil {
		return &pb.DeleteDishResponse{
			Status: http.StatusBadRequest,
			Error:  "couldn't parse uuid",
		}, nil
	}
	d, err := h.services.Dish.Delete(id)
	if err != nil {
		if errors.Is(err, core.ErrInternal) {
			return &pb.DeleteDishResponse{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}, nil
		}
		return &pb.DeleteDishResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}
	return &pb.DeleteDishResponse{
		Status:      http.StatusOK,
		Id:          d.Id.String(),
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Quantity:    d.Quantity,
		IsAvailable: d.IsAvailable,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}, nil
}

func (h Handler) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	if r == nil {
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  "empty request",
		}, nil
	}
	if r.Dishes == nil || len(r.Dishes) == 0 {
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  "order must have a least one dish",
		}, nil
	}
	dishes := make([]core.Dish, len(r.Dishes))
	for i, v := range r.Dishes {
		if v == nil {
			break
		}
		id, err := uuid.Parse(v.Id)
		if err != nil {
			return &pb.CreateOrderResponse{
				Status: http.StatusBadRequest,
				Error:  "incorrect dish id in order",
			}, nil
		}
		dishes[i] = core.Dish{
			Id:          id,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			Quantity:    v.Quantity,
			IsAvailable: v.IsAvailable,
			CreatedAt:   v.CreatedAt.AsTime(),
			UpdatedAt:   v.UpdatedAt.AsTime(),
		}
	}
	id, err := uuid.Parse(r.UserId)
	o := core.Order{
		UserId:          id,
		Dishes:          dishes,
		SpecialRequests: r.SpecialRequests,
	}
	createdOrd, err := h.services.Order.Create(o)
	if err != nil {
		if errors.Is(err, core.ErrInternal) {
			return &pb.CreateOrderResponse{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}, nil
		}
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}
	resDishes := make([]*pb.Dish, len(createdOrd.Dishes))
	for i, cd := range createdOrd.Dishes {
		resDishes[i] = &pb.Dish{
			Id:          cd.Id.String(),
			Name:        cd.Name,
			Description: cd.Description,
			Price:       cd.Price,
			Quantity:    cd.Quantity,
			IsAvailable: cd.IsAvailable,
			CreatedAt:   timestamppb.New(cd.CreatedAt),
			UpdatedAt:   timestamppb.New(cd.UpdatedAt),
		}
	}
	return &pb.CreateOrderResponse{
		Status:          http.StatusOK,
		Id:              createdOrd.Id.String(),
		Dishes:          resDishes,
		UserId:          createdOrd.UserId.String(),
		OrderStatus:     createdOrd.Status,
		SpecialRequests: createdOrd.SpecialRequests,
		ReadyAt:         timestamppb.New(createdOrd.ReadyAt),
		CreatedAt:       timestamppb.New(createdOrd.CreatedAt),
		UpdatedAt:       timestamppb.New(createdOrd.UpdatedAt),
	}, nil
}

func (h Handler) GetOrder(ctx context.Context, r *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	if r == nil {
		return &pb.GetOrderResponse{
			Status: http.StatusBadRequest,
			Error:  "empty request",
		}, nil
	}
	id, err := uuid.Parse(r.Id)
	order, err := h.services.Order.Get(id)
	if err != nil {
		if errors.Is(err, core.ErrInternal) {
			return &pb.GetOrderResponse{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}, nil
		}
		return &pb.GetOrderResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}
	resDishes := make([]*pb.Dish, len(order.Dishes))
	for i, cd := range order.Dishes {
		resDishes[i] = &pb.Dish{
			Id:          cd.Id.String(),
			Name:        cd.Name,
			Description: cd.Description,
			Price:       cd.Price,
			Quantity:    cd.Quantity,
			IsAvailable: cd.IsAvailable,
			CreatedAt:   timestamppb.New(cd.CreatedAt),
			UpdatedAt:   timestamppb.New(cd.UpdatedAt),
		}
	}
	return &pb.GetOrderResponse{
		Status:          http.StatusOK,
		Id:              order.Id.String(),
		Dishes:          resDishes,
		UserId:          order.UserId.String(),
		OrderStatus:     order.Status,
		SpecialRequests: order.SpecialRequests,
		ReadyAt:         timestamppb.New(order.ReadyAt),
		CreatedAt:       timestamppb.New(order.CreatedAt),
		UpdatedAt:       timestamppb.New(order.UpdatedAt),
	}, nil
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

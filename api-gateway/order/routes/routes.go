package routes

import (
	"api-gateway/core"
	h "api-gateway/http_tools"
	"api-gateway/order/pkg/pb"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllDishes(c *gin.Context, s pb.OrderServiceClient) {
	req := pb.GetAllDishRequest{}
	//if !h.BindRequestBody(c, &req) {
	//	return
	//}
	limitP := c.Query("limit")
	offsetP := c.Query("offset")
	//limit, err := strconv.Atoi(limitP)
	limit, err := strconv.ParseInt(limitP, 10, 64)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "couldn't get limit from parameter")
		return
	}
	//offset, err := strconv.Atoi(offsetP)
	offset, err := strconv.ParseInt(offsetP, 10, 64)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "couldn't get offset from parameter")
		return
	}
	if limit > 50 {
		h.ErrorResponse(c, http.StatusBadRequest, "limit is too big")
		return
	}
	req.Limit = limit
	req.Offset = offset
	res, err := s.GetAllDishes(context.Background(), &req)
	if err != nil {
		h.ErrorResponse(c, res.Status, res.Error)
		return
	}
	c.JSON((int)(res.Status), &res)
}

// TODO mb check if role good
func CreateDish(c *gin.Context, s pb.OrderServiceClient) {
	var req pb.CreateDishRequest
	if !h.BindRequestBody(c, &req) {
		return
	}
	res, err := s.CreateDish(context.Background(), &req)
	if err != nil {
		h.ErrorResponse(c, res.Status, res.Error)
		return
	}
	c.JSON((int)(res.Status), &res)
}

func GetDish(c *gin.Context, s pb.OrderServiceClient) {
	id := c.Param("id")
	res, err := s.GetDish(context.Background(), &pb.GetDishRequest{Id: id})
	if err != nil {
		h.ErrorResponse(c, res.Status, res.Error)
		return
	}
	c.JSON((int)(res.Status), &res)
}

func UpdateDish(c *gin.Context, s pb.OrderServiceClient) {
	var req pb.UpdateDishRequest
	if !h.BindRequestBody(c, &req) {
		return
	}
	res, err := s.UpdateDish(context.Background(), &req)
	if err != nil || res.Error != "" {
		h.ErrorResponse(c, res.Status, res.Error)
		return
	}
	c.JSON((int)(res.Status), &res)
}

func DeleteDish(c *gin.Context, s pb.OrderServiceClient) {
	id := c.Param("id")
	res, err := s.DeleteDish(context.Background(), &pb.DeleteDishRequest{Id: id})
	if err != nil {
		h.ErrorResponse(c, res.Status, res.Error)
		return
	}
	c.JSON((int)(res.Status), &res)
}

func CreateOrder(c *gin.Context, s pb.OrderServiceClient) {
	var req pb.CreateOrderRequest
	id, err := h.GetClientId(c)
	if err != nil {
		h.ErrorResponse(c, http.StatusUnauthorized, core.ErrTokenInvalid.Error())
		return
	}
	if !h.BindRequestBody(c, &req) {
		return
	}
	req.UserId = id.String()
	res, err := s.CreateOrder(context.Background(), &req)
	if err != nil {
		h.ErrorResponse(c, res.Status, res.Error)
		return
	}
	c.JSON((int)(res.Status), &res)
}

func GetOrder(c *gin.Context, s pb.OrderServiceClient) {
	id := c.Param("id")
	res, err := s.GetOrder(context.Background(), &pb.GetOrderRequest{Id: id})
	if err != nil {
		h.ErrorResponse(c, res.Status, res.Error)
		return
	}
	c.JSON((int)(res.Status), &res)
}

package grpc

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"user-svc/internal/core"
	"user-svc/internal/service"
	"user-svc/pkg/pb"
)

const (
	authHeader     = "Authorization"
	userContextKey = "clientId"
)

type Handler struct {
	services *service.Service
	pb.UnimplementedUserServiceServer
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h Handler) Auth(c context.Context, authReq *pb.AuthRequest) (*pb.AuthResponse, error) {
	gmail, err := service.GetGmail(authReq.IdToken)
	if err != nil {
		return &pb.AuthResponse{
			Status: http.StatusUnauthorized,
			Error:  core.ErrTokenInvalid.Error(),
		}, nil
	}
	
	jwt, err := h.services.GetJWT(gmail)
	if err != nil {
		logrus.Errorf("while generating JWT" + err.Error())
		return &pb.AuthResponse{
			Status: http.StatusInternalServerError,
			Error:  core.ErrInternal.Error(),
		}, nil
	}

	return &pb.AuthResponse{
		Status:       http.StatusOK,
		Jwt:          jwt.Token,
		IsRegistered: jwt.IsRegistered,
	}, nil
}

func (h Handler) SignUp(c context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	id, registered, _, err := h.services.ParseToken(req.Jwt)
	if err != nil {
		return &pb.SignUpResponse{
			Status: http.StatusUnauthorized,
			Error:  core.ErrTokenInvalid.Error(),
		}, nil
	}
	if registered {
		return &pb.SignUpResponse{
			Status: http.StatusBadRequest,
			Error:  "user already registered",
		}, nil
	}
	newUser, err := h.services.User.RegisterUser(id, core.User{Username: req.Username, Nickname: req.Nickname, Role: req.Role})
	if err != nil {
		return &pb.SignUpResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	jwt, err := h.services.GenerateJWT(newUser.Id, true, newUser.Role)
	if err != nil {
		logrus.Errorf("while generating JWT" + err.Error())
		return &pb.SignUpResponse{
			Status: http.StatusInternalServerError,
			Error:  core.ErrInternal.Error(),
		}, nil
	}
	return &pb.SignUpResponse{
		Status:       http.StatusOK,
		Jwt:          jwt.Token,
		Id:           newUser.Id.String(),
		Gmail:        newUser.Gmail,
		Username:     newUser.Username,
		Nickname:     newUser.Nickname,
		IsRegistered: true,
		Role:         newUser.Role,
	}, nil
}

func (h Handler) Get(c context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userId, err := uuid.Parse(req.IdToken)
	if err != nil {
		return &pb.GetUserResponse{
			Status: http.StatusForbidden,
			Error:  "could not parse uuid from id parameter",
		}, nil
	}
	ok := h.services.User.Exists(userId)
	if !ok {
		return &pb.GetUserResponse{
			Status: http.StatusNotFound,
			Error:  "couldn't find user",
		}, nil
	}
	user, err := h.services.User.GetById(userId)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return &pb.GetUserResponse{
				Status: http.StatusNotFound,
				Error:  "couldn't find user",
			}, nil
		}
		return &pb.GetUserResponse{
			Status: http.StatusInternalServerError,
			Error:  core.ErrInternal.Error(),
		}, nil
	}
	return &pb.GetUserResponse{
		Status:       http.StatusOK,
		Id:           user.Id.String(),
		Gmail:        user.Gmail,
		Username:     user.Username,
		Nickname:     user.Nickname,
		IsRegistered: user.IsRegistered,
		Role:         user.Role,
	}, nil
}

func (h Handler) Validate(c context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	header := req.Jwt
	id, registered, role, err := h.services.ParseToken(header)
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		}, nil

	}
	if !registered {
		return &pb.ValidateResponse{
			Status: http.StatusForbidden,
			Error:  "Client isn't registered",
		}, nil
	}
	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: id.String(),
		Role:   role,
	}, nil
}

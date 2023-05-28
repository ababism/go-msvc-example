package routes

import (
	"api-gateway/core"
	h "api-gateway/http_tools"
	"api-gateway/user/pkg/pb"
	"context"
	"github.com/gin-gonic/gin"
)

func Start(ctx *gin.Context, c pb.UserServiceClient) {
	type idToken struct {
		IdToken string `json:"id_token" binding:"required"`
	}
	var t idToken
	if !h.BindRequestBody(ctx, &t) {
		return
	}
	res, err := c.Auth(context.Background(), &pb.AuthRequest{
		IdToken: t.IdToken,
	})

	if err != nil {
		h.NewErrorResponse(ctx, res.Status, core.CodeAccessDenied, res.Error)
		return
	}
	ctx.JSON((int)(res.Status), &res)
}
func SignUp(ctx *gin.Context, c pb.UserServiceClient) {
	type req struct {
		Username string `json:"username" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
		Role     string `json:"role" binding:"required"`
		core.JWT
	}
	var r req
	if !h.BindRequestBody(ctx, &r) {
		return
	}
	res, err := c.SignUp(context.Background(), &pb.SignUpRequest{
		Username: r.Username,
		Nickname: r.Nickname,
		Jwt:      r.JWT.Token,
		Role:     r.Role,
	})
	if err != nil {
		h.NewErrorResponse(ctx, res.Status, core.CodeIncorrectBody, res.Error)
		return
	}
	ctx.JSON((int)(res.Status), &res)
}

func GetUser(c *gin.Context, s pb.UserServiceClient) {
	type req struct {
		Username string `json:"username" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
		core.JWT
	}
	userId := c.Param("id")

	res, err := s.Get(context.Background(), &pb.GetUserRequest{IdToken: userId})
	if err != nil {
		h.NewErrorResponse(c, res.Status, core.CodeIncorrectBody, res.Error)
		return
	}
	c.JSON((int)(res.Status), &res)
}

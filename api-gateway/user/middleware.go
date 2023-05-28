package user

import (
	"api-gateway/core"
	h "api-gateway/http_tools"
	"api-gateway/user/pkg/pb"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthMiddlewareConfig struct {
	svc *ServiceClient
}

func InitAuthMiddleware(svc *ServiceClient) AuthMiddlewareConfig {
	return AuthMiddlewareConfig{svc}
}

func (c *AuthMiddlewareConfig) AuthRequired(ctx *gin.Context) {
	//authorization := ctx.Request.Header.Get("authorization")
	header := ctx.GetHeader(h.AuthHeader)

	res, err := c.svc.Client.Validate(context.Background(), &pb.ValidateRequest{
		Jwt: header,
	})

	if err != nil || res.Status != http.StatusOK {
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		h.NewErrorResponse(ctx, res.Status, core.CodeAccessDenied, core.ErrTokenInvalid.Error())
		return
	}
	//if role != res.Role {
	//	//ctx.AbortWithStatus(http.StatusUnauthorized)
	//	h.NewErrorResponse(ctx, http.StatusForbidden, core.CodeAccessDenied, "client role do not match requirements")
	//	return
	//}
	ctx.Set(h.ClientIdCtxtKey, res.UserId)
	ctx.Set(h.ClientRoleCtxKey, res.Role)
	ctx.Next()
}
func (c *AuthMiddlewareConfig) ManagerRoleRequired(ctx *gin.Context) {
	//authorization := ctx.Request.Header.Get("authorization")
	header := ctx.GetHeader(h.AuthHeader)

	res, err := c.svc.Client.Validate(context.Background(), &pb.ValidateRequest{
		Jwt: header,
	})

	if err != nil || res.Status != http.StatusOK {
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		h.NewErrorResponse(ctx, res.Status, core.CodeAccessDenied, core.ErrTokenInvalid.Error())
		return
	}
	if res.Role != core.ManagerRole {
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		h.NewErrorResponse(ctx, http.StatusForbidden, core.CodeAccessDenied, "client role do not match requirements")
		return
	}
	ctx.Set(h.ClientIdCtxtKey, res.UserId)
	ctx.Set(h.ClientRoleCtxKey, res.Role)
	ctx.Next()
}

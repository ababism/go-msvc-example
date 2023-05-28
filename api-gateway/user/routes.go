package user

import (
	"api-gateway/config"
	"api-gateway/user/routes"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config) *ServiceClient {
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	auth := r.Group("/auth")
	{
		auth.POST("/start", svc.Start)
		auth.POST("/sign_up", svc.SignUp)
	}
	user := r.Group("/user")
	{
		user.GET("/:id", svc.GetUser)
	}
	return svc
}

func (svc *ServiceClient) Start(ctx *gin.Context) {
	routes.Start(ctx, svc.Client)
}

func (svc *ServiceClient) SignUp(ctx *gin.Context) {
	routes.SignUp(ctx, svc.Client)
}

func (svc *ServiceClient) GetUser(ctx *gin.Context) {
	routes.GetUser(ctx, svc.Client)
}

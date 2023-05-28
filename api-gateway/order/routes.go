package order

import (
	"api-gateway/config"
	"api-gateway/order/routes"
	"api-gateway/user"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *user.ServiceClient) *ServiceClient {
	a := user.InitAuthMiddleware(authSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	dish := r.Group("/dish")
	dish.Use(a.AuthRequired)
	{
		dish.GET("/menu", svc.GetAllDishes)
		dish.GET("/:id", svc.GetDish)
	}
	dishManager := r.Group("/dish")
	dishManager.Use(a.ManagerRoleRequired)
	{
		dishManager.POST("/create", svc.CreateDish)
		dishManager.PUT("/update", svc.UpdateDish)
		dishManager.DELETE("/delete/:id", svc.DeleteDish)
	}
	//dish := r.Group("/dish")
	//dish.Use(a.AuthRequired)
	//{
	//	dish.GET("/menu", svc.GetAllDishes)
	//	dish.GET("/:id", svc.GetDish)
	//}
	//manager := r.Group("/manage")
	//manager.Use(a.ManagerRoleRequired)
	//{
	//	manageDish := r.Group("/dish")
	//	{
	//		manageDish.POST("/create", svc.CreateDish)
	//		manageDish.PATCH("/update", svc.UpdateDish)
	//		manageDish.DELETE("/delete/:id", svc.DeleteDish)
	//	}
	//}
	order := r.Group("/order")
	order.Use(a.AuthRequired)
	{
		order.POST("/create", svc.CreateOrder)
		order.GET("/:id", svc.GetOrder)
	}
	return svc
}

func (svc *ServiceClient) GetAllDishes(ctx *gin.Context) {
	routes.GetAllDishes(ctx, svc.Client)
}

func (svc *ServiceClient) CreateDish(ctx *gin.Context) {
	routes.CreateDish(ctx, svc.Client)
}

func (svc *ServiceClient) GetDish(ctx *gin.Context) {
	routes.GetDish(ctx, svc.Client)
}

func (svc *ServiceClient) UpdateDish(ctx *gin.Context) {
	routes.UpdateDish(ctx, svc.Client)
}

func (svc *ServiceClient) DeleteDish(ctx *gin.Context) {
	routes.DeleteDish(ctx, svc.Client)
}

func (svc *ServiceClient) CreateOrder(ctx *gin.Context) {
	routes.CreateOrder(ctx, svc.Client)
}

func (svc *ServiceClient) GetOrder(ctx *gin.Context) {
	routes.GetOrder(ctx, svc.Client)
}

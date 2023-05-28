package main

import (
	"api-gateway/config"
	"api-gateway/order"
	"api-gateway/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

//func (s *Server) Run(port string, handler http.Handler) error {
//	s.httpServer = &http.Server{
//		Addr:           ":" + port,
//		Handler:        handler,
//		MaxHeaderBytes: 1 << 20, // 1 MB
//		ReadTimeout:    10 * time.Second,
//		WriteTimeout:   10 * time.Second,
//	}
//
//	return s.httpServer.ListenAndServe()
//}
//
//func (s *Server) Shutdown(ctx context.Context) error {
//	return s.httpServer.Shutdown(ctx)
//}

func main() {
	// Логирование и viper
	logrus.SetFormatter(new(logrus.JSONFormatter))
	c := initViperConfig()

	//srv := new(Server)
	r := gin.Default()
	fmt.Println(c.UserSvcUrl)
	fmt.Println(c.OrderSvcUrl)
	fmt.Println(viper.GetString("port"))
	userSvc := *user.RegisterRoutes(r, &c)
	order.RegisterRoutes(r, &c, &userSvc)

	err := r.Run(":" + viper.GetString("port"))
	if err != nil {
		logrus.Fatalf("error starting gin engine (server): %v", err.Error())
	}

	//if err = srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
	//	logrus.Fatalf("error loading env variables: %s", err.Error())
	//}
}

func initViperConfig() config.Config {
	c, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Error reading config file, %s", err)
	}
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	// Enable VIPER to read Environment Variables 😱🥹🔥
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	if err = viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file, %s", err)
	}
	return c
}

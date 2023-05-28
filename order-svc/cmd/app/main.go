package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"order-svc/internal/repository"
	"order-svc/internal/service"
	grpc_t "order-svc/internal/transport/grpc"
	"order-svc/pkg/pb"
)

func main() {
	// Логирование и viper
	logrus.SetFormatter(new(logrus.JSONFormatter))
	initViperConfig()

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: viper.GetString("db.password"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	lis, err := net.Listen("tcp", viper.GetString("port"))
	if err != nil {
		logrus.Fatalln("Failed to listing:", err)
	}
	fmt.Println("Auth Svc on", viper.GetString("port"))

	repos := repository.New(db)
	services := service.NewService(repos)
	handlers := grpc_t.NewHandler(services)

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, handlers)

	if err = grpcServer.Serve(lis); err != nil {
		logrus.Fatalln("Failed to serve:", err)
	}
	//if err = srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
	//	logrus.Fatalf("error loading env variables: %s", err.Error())
	//}
}

func initViperConfig() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	// Enable VIPER to read Environment Variables 😱🥹🔥
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file, %s", err)
	}
}

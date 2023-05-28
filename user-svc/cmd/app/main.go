package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"user-svc/internal/repository"
	"user-svc/internal/service"
	grpct "user-svc/internal/transport/grpc"
	"user-svc/pkg/pb"
)

func main() {
	// Логирование и viper
	logrus.SetFormatter(new(logrus.JSONFormatter))
	initViperConfig()
	// подключаем постгрес
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

	// Инциализируем подключение к Google Oauth2 API
	service.InitializeOAuthGoogle()

	lis, err := net.Listen("tcp", viper.GetString("port"))
	if err != nil {
		logrus.Fatalln("Failed to listing:", err)
	}
	fmt.Println("Auth Svc on", viper.GetString("port"))

	// инициализируем трехслойную чистую архитектуру
	// репозиторий
	repos := repository.New(db)
	// бизнес-логика
	services := service.NewService(repos)
	// транспортный слой
	handlers := grpct.NewHandler(services)

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, handlers)

	if err = grpcServer.Serve(lis); err != nil {
		logrus.Fatalln("Failed to serve:", err)
	}
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

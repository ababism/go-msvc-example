package user

import (
	"api-gateway/config"
	"api-gateway/user/pkg/pb"
	"fmt"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.UserServiceClient
}

func InitServiceClient(c *config.Config) pb.UserServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.UserSvcUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewUserServiceClient(cc)
}

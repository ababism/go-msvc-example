package order

import (
	"api-gateway/config"
	proto "api-gateway/order/pkg/pb"
	"fmt"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client proto.OrderServiceClient
}

func InitServiceClient(c *config.Config) proto.OrderServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.OrderSvcUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return proto.NewOrderServiceClient(cc)
}

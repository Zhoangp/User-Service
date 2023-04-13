package client

import (
	"fmt"
	"github.com/Zhoangp/User-Service/config"
	"github.com/Zhoangp/User-Service/pb"
	"google.golang.org/grpc"
)
func InitServiceClient(c *config.Config) (pb.FileServiceClient, error) {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.OtherServices.FileUrl, grpc.WithInsecure())
	//grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(c.Service.MaxSizeMess), grpc.MaxCallSendMsgSize(c.Service.MaxSizeMess))
	if err != nil {
		fmt.Println("Could not connect:", err)
		return nil, err
	}
	return pb.NewFileServiceClient(cc), nil
}
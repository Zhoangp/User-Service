package main

import (
	"fmt"
	"github.com/Zhoangp/User-Service/config"
	userhttp "github.com/Zhoangp/User-Service/internal/delivery/http"
	"github.com/Zhoangp/User-Service/internal/repo"
	"github.com/Zhoangp/User-Service/internal/usecase"
	"github.com/Zhoangp/User-Service/pb"
	"github.com/Zhoangp/User-Service/pkg/database/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cf, err := config.LoadConfig("config/config.yml")
	if err != nil {
		panic(err)
	}
	gormDb, err := mysql.NewMysql(cf)
	if err != nil {
		fmt.Println(err)
		return
	}
	lis, err := net.Listen("tcp", ":"+cf.Service.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Auth Svc on", cf.Service.Port)

	repoUser := repo.NewUserRepository(gormDb)
	useCaseUser := usecase.NewUserUseCase(repoUser, cf)
	hdlUser := userhttp.NewUserHandler(useCaseUser)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, hdlUser)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

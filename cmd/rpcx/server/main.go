package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"github.com/rc452860/vnet/cmd/rpcx/server/service"

	"github.com/rc452860/vnet/cmd/rpcx"
	"google.golang.org/grpc"
)

type Args struct {
	Port   uint
	Source string
}

func main() {
	arg := &Args{}
	flag.UintVar(&arg.Port, "port", 5050, "server port")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(int(arg.Port)))
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	server := grpc.NewServer()
	rpcx.RegisterUserServiceServer(server, &service.UserService{})
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

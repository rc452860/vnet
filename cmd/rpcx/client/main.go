package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/rc452860/vnet/cmd/rpcx"
	"google.golang.org/grpc"
)

type Config struct {
	NodeId     int
	Token      string
	RpcAddress string
}

func main() {
	config := &Config{}
	flag.StringVar(&config.RpcAddress, "RpcAddress", "localhost:5050", "rpc address")
	flag.Parse()

	conn, err := grpc.Dial(config.RpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rpcx.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.PullEnableUsers(ctx, &rpcx.PullEnableUsersRequest{})
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	for _, item := range r.GetEnableUsers() {
		fmt.Println(item.Method)
	}
}

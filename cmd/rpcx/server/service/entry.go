package service

import (
	"net"

	"github.com/rc452860/vnet/cmd/rpcx"
	"github.com/rc452860/vnet/cmd/rpcx/config"
	"github.com/rc452860/vnet/common/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// Start 程序入口
func Start() {
	log.Info("server started...")
	err := InitDS(viper.GetString(config.S_DS))
	if err != nil {
		log.Error("init database connection failed, please check ds config: %v", err)
		return
	}

	lis, err := net.Listen("tcp", ":"+viper.GetString(config.S_RPCPort))
	if err != nil {
		log.Error("failed to serve: %v", err)
	}
	log.Info("server listen on:%s", viper.GetString(config.S_RPCPort))
	server := grpc.NewServer()
	rpcx.RegisterUserServiceServer(server, &UserService{})
	if err := server.Serve(lis); err != nil {
		log.Error("failed to serve: %v", err)
	}
}

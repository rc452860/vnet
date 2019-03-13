package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rc452860/vnet/cmd/rpcx"
)

func ExampleUserService_PullEnableUsers() {
	err := InitDS("root:killer@tcp(localhost:3306)/ssrpanel_dev?parseTime=true&loc=UTC")
	if err != nil {
		log.Fatalln(err)
		return
	}

	service := &UserService{}
	// secondsEastOfUTC := int((8 * time.Hour).Seconds())
	// beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	start, _ := time.Parse("2006-01-02 15:04:05", "2018-11-19 00:07:56")
	log.Println(start.Location().String())
	log.Println(start.Unix())
	// start, _ = time.ParseInLocation("2006-01-02 15:04:05", "2018-11-19 00:07:57", beijing)
	log.Println(start.Unix())
	result, err := service.PullEnableUsers(context.Background(), &rpcx.PullEnableUsersRequest{
		NodeId:               2,
		Token:                "abc",
		UpdateTime:           start.Unix(),
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(len(result.EnableUsers))
	rln := len(result.EnableUsers)
	if rln > 5 {
		rln = 5
	}
	fmt.Println(result.EnableUsers[:rln])
	//Output:
}

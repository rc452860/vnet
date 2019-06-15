package client

import (
	"fmt"
	"github.com/rc452860/vnet/model"
	"github.com/sirupsen/logrus"
)

func ExampleGetNodeInfo() {
	HOST = "http://localhost"
	logrus.SetLevel(logrus.DebugLevel)
	result := GetNodeInfo(2,"txsnvhghmrmg4pjm")
	fmt.Printf("value: %+v \n",result)
	//Output:
}

func ExampleGetUserList(){
	HOST = "http://localhost"
	logrus.SetLevel(logrus.DebugLevel)
	result := GetUserList(1,"txsnvhghmrmg4pjm")
	fmt.Printf("value: %+v\n",result)
	//Output:
}

func ExamplePostAllUserTraffic() {
	HOST = "http://localhost"
	logrus.SetLevel(logrus.DebugLevel)
	PostAllUserTraffic([]*model.UserTraffic{
		{
			1,200,200,
		},
	},1,"txsnvhghmrmg4pjm")
	//Output:
}

func ExamplePostNodeOnline() {
	HOST = "http://localhost"
	logrus.SetLevel(logrus.DebugLevel)
	PostNodeOnline([]*model.NodeOnline{
		{
			1,
			"192.168.1.1",
		},
	},1,"txsnvhghmrmg4pjm")
	//Output:
}

func ExamplePostNodeStatus() {
	HOST = "http://localhost"
	logrus.SetLevel(logrus.DebugLevel)
	PostNodeStatus(model.NodeStatus{
		"10%",
		"10%",
		"10%",
	},1,"txsnvhghmrmg4pjm")

	//Output:
}
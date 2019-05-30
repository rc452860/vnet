package api

import (
	"fmt"
	"github.com/rc452860/vnet/model"
	"github.com/sirupsen/logrus"
)

func ExampleGetNodeInfo() {
	logrus.SetLevel(logrus.DebugLevel)
	result := GetNodeInfo(1,"txsnvhghmrmg4pjm")
	fmt.Printf("value: %+v \n",result)
	//Output:
}

func ExampleGetUserList(){
	logrus.SetLevel(logrus.DebugLevel)
	result := GetUserList(1,"txsnvhghmrmg4pjm")
	fmt.Printf("value: %+v\n",result)
	//Output:
}

func ExamplePostAllUserTraffic() {
	logrus.SetLevel(logrus.DebugLevel)
	PostAllUserTraffic([]model.UserTraffic{
		{
			1,200,200,
		},
	},1,"txsnvhghmrmg4pjm")
	//Output:
}

func ExamplePostNodeOnline() {
	logrus.SetLevel(logrus.DebugLevel)
	PostNodeOnline([]model.NodeOnline{
		{
			1,
			"192.168.1.1",
		},
	},1,"txsnvhghmrmg4pjm")
	//Output:
}

func ExamplePostNodeStatus() {
	logrus.SetLevel(logrus.DebugLevel)
	PostNodeStatus(model.NodeStatus{
		"10%",
		"10%",
		"10%",
	},1,"txsnvhghmrmg4pjm")

	//Output:
}
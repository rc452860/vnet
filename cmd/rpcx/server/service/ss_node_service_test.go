package service

import "fmt"

func ExampleSsNodeService_VerifyNode() {
	InitDBTest()
	sns := GetSsNodeServiceInstance()
	result := sns.VerifyNode(3, "abc")
	fmt.Println(result)
	result = sns.VerifyNode(2, "abc")
	fmt.Println(result)
	//Output:
	//
}

func InitDBTest() {
	err := InitDS("root:killer@tcp(localhost:3306)/ssrpanel_dev?parseTime=true&loc=UTC")
	if err != nil {
		panic(err)
	}
}

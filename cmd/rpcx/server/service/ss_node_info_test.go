package service

func ExampleRecordSsNodeInfo() {
	InitDBTest()

	s := GetSsNodeInfoInstance()
	s.RecordSsNodeInfo(2, "hello")
	//Output:
}

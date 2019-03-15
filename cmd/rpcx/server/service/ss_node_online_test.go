package service

func ExampleRecordSsNodeOnline() {
	InitDBTest()

	s := GetSsNodeOnlineServiceInstance()
	s.RecordSsNodeOnline(2, 100)
	//Output:
}

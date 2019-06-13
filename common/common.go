package common

type TrafficReport interface{
	Upload(uid int,n int64)
	Download(uid int,n int64)
}


type OnlineReport interface{
	Online(uid int,ip string)
}
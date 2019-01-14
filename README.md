[![Build Status](https://travis-ci.org/rc452860/vnet.svg?branch=master)](https://travis-ci.org/rc452860/vnet)
[![Go Report Card](https://goreportcard.com/badge/github.com/rc452860/vnet)](https://goreportcard.com/report/github.com/rc452860/vnet)

![you have me](./assert/donate.png)

## 功能介绍
Vnet是一个代理工具,在某些网络条件受到限制的情况先提供突破服务.

## 开发计划
- [ x ] shadowsocsk代理协议
- [  ] kcp自定义协议
- [ x ] 代理服务流量统计
- [  ] 代理服务速度监控(正在测试)
- [  ] restful api
- [  ] 服务器cpu 内存 硬盘 上传下载速度监控

## 已知问题
- [ ] log formatter setdepth 多线程问题待改进

## 编译方式
```
go get -u -d github/rc452860/vnet
```

进入$gopath/rc452860/vnet目录

```
go build cmd/server/server.go
```

## 运行
linux:
```
chmod +x server && ./server
```

windows:
运行server.exe
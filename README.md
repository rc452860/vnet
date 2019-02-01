[![Build Status](https://travis-ci.org/rc452860/vnet.svg?branch=master)](https://travis-ci.org/rc452860/vnet)
[![Go Report Card](https://goreportcard.com/badge/github.com/rc452860/vnet)](https://goreportcard.com/report/github.com/rc452860/vnet)
[![Join the chat at https://gitter.im/rc452860/vnet](https://badges.gitter.im/rc452860/vnet.svg)](https://gitter.im/rc452860/vnet?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

<img src="./assert/donate.png" width="300" title="you have me!">

## 功能介绍
Vnet是一个代理工具,在某些网络条件受到限制的情况先提供突破服务.

## 开发计划
- [ x ] shadowsocsk代理协议
- [  ] kcp自定义协议
- [ x ] 代理服务流量统计
- [ x ] 代理服务速度监控
- [  ] restful api(进行中)
- [ x ] 服务器cpu 内存 硬盘 上传下载速度监控

## 已知问题
- [ ] log formatter setdepth 多线程问题待改进



## 运行
linux去release页面下载对应的指令集二进制文件给运行权限直接运行,根据提示输入对应的配置好

window直接运行exe

linx:
```
chmod +x server && ./server
```

## 编译方式
```
go get -u -d github.com/rc452860/vnet/...
```

进入$gopath/rc452860/vnet目录

```
go build cmd/server/server.go
```

## 直接使用方式(无需编译)
在release页面下载最新的对应的可执行文件并赋予可执行权限
列如64位linux系统
```
wget https://github.com/rc452860/vnet/releases/download/v0.0.4/vnet_linux_amd64 -O vnet &&chmod +x vnet
./vnet
```
按照提示输入数据库等配置信息即可完成

后台运行使用`nohup`工具辅助
```
nohup ./vnet >vnet.log 2>&1 &
```

## 支持加密方式
```
aes-256-cfb
bf-cfb
chacha20
chacha20-ietf
aes-128-cfb
aes-192-cfb
aes-128-ctr
aes-192-ctr
aes-256-ctr
cast5-cfb
des-cfb
rc4-md5
salsa20
aes-256-gcm
aes-192-gcm
aes-128-gcm
chacha20-ietf-poly1305
```

## 注意事项
config.json配置文件中的所有时间单位都为毫秒
升级后续删除原有config.json重新生成
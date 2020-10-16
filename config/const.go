package config

import "time"

const (
	IvReloadConfig     = time.Second * 3  //客户端重新加载配置文件
	IvSent             = time.Second * 3  //客户端发送消息
	IvGetTraffic       = time.Second * 1  //客户端更新流量
	IvCheckIPvNSupport = time.Second * 10 //客户端检查IP支持
	IvReconnect        = time.Second * 5  //客户端重连
)

const (
	DnsIpv4 = "223.5.5.5"    //Aliyun DNS
	DnsIpv6 = "2400:3200::1" //Aliyun DNS
)

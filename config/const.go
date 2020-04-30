package config

import "time"

const (
	IntervalReloadConfig     = time.Second * 5
	IntervalSent             = time.Second * 3
	IntervalHeartbeat        = time.Minute * 1
	IntervalGetTraffic       = time.Second * 1
	IntervalCheckIPvNSupport = time.Second * 30
)

const (
	DnsIpv4 = "199.85.126.10"        //Norton ConnectSafe
	DnsIpv6 = "2606:4700:4700::1111" //Cloudflare
)

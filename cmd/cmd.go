package cmd

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	Version    = "0.1.6"
	Host       = ""
	Port       = 9001
	HTTPPort   = 9002
	MultiCore  = true
	Tick       = time.Second * 3
	Filename   = "./config.json"
	ID         = ""
	HasConvStr = true
)

var (
	serverCmd, clientCmd, uuidCmd, systemCmd, trafficCmd *flag.FlagSet

	host       string
	port       int
	httpPort   int
	multicore  bool
	interval   time.Duration
	filename   string
	id         string
	hasConvStr bool
)

type Params struct {
	Host       string
	Port       int
	HTTPPort   int
	Multicore  bool
	Interval   time.Duration
	Filename   string
	Id         string
	HasConvStr bool
}

func NewParams(host string, port, httpPort int, multicore bool,
	interval time.Duration, filename string, id string, hasConvStr bool) (p *Params) {
	return &Params{
		Host:       host,
		Port:       port,
		HTTPPort:   httpPort,
		Multicore:  multicore,
		Interval:   interval,
		Filename:   filename,
		Id:         id,
		HasConvStr: hasConvStr,
	}
}

type Cmd struct {
	T string
	*Params
}

func NewCmd(t string, p *Params) *Cmd {
	return &Cmd{t, p}
}

func init() {
	serverCmd = flag.NewFlagSet("server", flag.ExitOnError)
	clientCmd = flag.NewFlagSet("client", flag.ExitOnError)
	uuidCmd = flag.NewFlagSet("uuid", flag.ExitOnError)
	systemCmd = flag.NewFlagSet("client", flag.ExitOnError)
	trafficCmd = flag.NewFlagSet("traffic", flag.ExitOnError)
}

func usage() {
	fmt.Fprintf(os.Stderr, `ServerStatus version: ServerStatus/%s
Usage: ServerStatus <command>

Available commands:
    server               {启动服务端 [ServerStatus server [-h host] [-p TCPPort] [-hp HTTPPort] [-m multicore] [-c filename]]}
    client               {启动客户端, -s 转换相关字段为字符串(默认开启) [ServerStatus client [-s=false] [-h host] [-p port] [-m multicore] [-t tick] [-id uuid]]}
    system               {输出系统当前的参数 [ServerStatus system]}
    uuid                 {生成uuid [ServerStatus uuid]}
    traffic              {监听网卡实时流量 [ServerStatus traffic]}
    help                 {帮助 [ServerStatus help [--help]]}
`, Version)
	flag.PrintDefaults()
}

func unknownCommand() {
	fmt.Printf("Unknown command: \n\n")
	usage()
}

func handlerServer(args []string) {
	serverCmd.StringVar(&host, "h", Host, "listen host")
	serverCmd.IntVar(&port, "p", Port, "listen tcp port")
	serverCmd.IntVar(&httpPort, "hp", HTTPPort, "listen http port")
	serverCmd.BoolVar(&multicore, "m", MultiCore, "multicore")
	serverCmd.DurationVar(&interval, "t", Tick, "pushing tick")
	serverCmd.StringVar(&filename, "c", Filename, "use config.json")
	serverCmd.Parse(args)
}

func handlerClient(args []string) {
	clientCmd.StringVar(&host, "h", Host, "listen host")
	clientCmd.IntVar(&port, "p", Port, "listen port")
	clientCmd.BoolVar(&multicore, "m", MultiCore, "multicore")
	clientCmd.DurationVar(&interval, "t", Tick, "pushing tick")
	clientCmd.StringVar(&id, "id", ID, "uuid")
	clientCmd.BoolVar(&hasConvStr, "s", HasConvStr, "convert string")
	clientCmd.Parse(args)
}

func handlerUUID(args []string) {
	uuidCmd.Parse(args)
}

func handlerSystem(args []string) {
	systemCmd.Parse(args)
}

func handlerTraffic(args []string) {
	trafficCmd.Parse(args)
}

func Run() *Cmd {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	t := os.Args[1]
	args := os.Args[2:]

	switch t {
	case "server":
		handlerServer(args)
	case "client":
		handlerClient(args)
	case "uuid":
		handlerUUID(args)
	case "system":
		handlerSystem(args)
	case "traffic":
		handlerTraffic(args)
	case "help":
		usage()
	default:
		unknownCommand()
	}

	return NewCmd(t, NewParams(host, port, httpPort, multicore, interval, filename, id, hasConvStr))
}

# ServerStatus
一个实时监听服务器状态的命令行工具，使用go实现

## 一个例子
[https://status.flxxyz.com/](https://status.flxxyz.com/) 这只是一个分离出来的web界面，其中数据汇总均是通过该项目完成。

## 使用(服务端)
1. 前往 [Release](https://github.com/flxxyz/ServerStatus/releases) 页，下载系统对应的版本
2. 解压**ServerStatus**压缩包，进入**ServerStatus**文件夹
3. 执行`./ServerStatus server`, 默认监听端口`tcp: 9001; http: 9002`

## 使用(客户端)
1. 前往 [Release](https://github.com/flxxyz/ServerStatus/releases) 页，下载系统对应的版本
2. 解压**ServerStatus**压缩包，进入**ServerStatus**文件夹
3. `./ServerStatus uuid` 复制生成的唯一id
4. 将客户端信息填入服务端`config.json`中
5. 执行`./ServerStatus client -h [服务器地址] -p 服务器端口 -id [客户端id]`接入，默认参数可忽略
6. 服务端控制台收到当前客户端消息即为连接成功

## 读取数据
- [x] http: [server ip]:9002
- [ ] websocket: [server ip]:9002

## 更多命令
执行 `ServerStatus help` 命令获取更多参数信息
```text
ServerStatus version: ServerStatus/0.2.0
Usage: ServerStatus <command>

Available commands:
    server               {启动服务端 [ServerStatus server [-h host] [-p TCPPort] [-hp HTTPPort] [-m multicore] [-c filename]]}
    client               {启动客户端, -s [ServerStatus client [-h host] [-p port] [-m multicore] [-t tick] [-id uuid]]}
    system               {输出系统当前的参数 [ServerStatus system]}
    uuid                 {生成uuid [ServerStatus uuid]}
    traffic              {监听网卡实时流量 [ServerStatus traffic]}
    help                 {帮助 [ServerStatus help [--help]]}
```
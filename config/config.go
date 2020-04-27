package config

import (
    "ServerStatus/msg"
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
    "sync"
    "time"
)

var (
    conf *Config
)

type Config struct {
    Filename       string
    LastModifyTime int64
    Lock           *sync.RWMutex
    Data           []msg.Node
}

func NewConfig(filename string, data []msg.Node) *Config {
    conf = &Config{
        Filename: filename,
        Data:     data,
        Lock:     &sync.RWMutex{},
    }

    conf.parse()

    go conf.reload()

    return conf
}

func (c *Config) parse() bool {
    //记录最后修改时间
    fileInfo, _ := os.Stat(c.Filename)
    c.LastModifyTime = fileInfo.ModTime().Unix()

    //读取文件内容
    file, err := ioutil.ReadFile(c.Filename)
    if err != nil {
        log.Fatal("read config file fail, ", err.Error())
    }

    //解json
    c.Lock.Lock()
    if err = json.Unmarshal(file, &c.Data); err != nil {
        log.Fatal("JSON Unmarshal error, ", err.Error())
        return false
    }
    c.Lock.Unlock()

    return true
}

func (c *Config) reload() {
    ticker := time.NewTicker(time.Second * IntervalReloadConfig)
    for range ticker.C {
        fileInfo, _ := os.Stat(c.Filename)
        currModifyTime := fileInfo.ModTime().Unix()
        if currModifyTime > c.LastModifyTime {
            if c.parse() {
                log.Println("重新加载配置文件conf.json")
            }
        }
    }
}

func GetConf(key string) (node msg.Node, ok bool) {
    var data = make(map[string]msg.Node, 0)
    for _, node := range conf.Data {
        data[node.Id] = node
    }
    node, ok = data[key]
    return
}

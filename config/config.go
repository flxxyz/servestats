package config

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	conf *Config
)

type Config struct {
	Filename       string
	LastModifyTime int64
	Lock           *sync.RWMutex
	Data           []interface{}
}

func NewConfig(filename string, data []interface{}) *Config {
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
		log.Println("JSON Unmarshal error, ", err.Error())
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

func (c *Config) Get(key string) (node interface{}, ok bool) {
	data := make(map[string]interface{}, 0)
	for i, _ := range conf.GetData() {
		n := conf.Data[i].(map[string]interface{})
		data[n["id"].(string)] = n
	}
	node, ok = data[key]
	return
}

func (c *Config) GetData() (data []interface{}) {
	c.Lock.Lock()
	data = conf.Data
	c.Lock.Unlock()
	return
}

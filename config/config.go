package config

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	Conf *Config
)

// ConfigNode 配置文件中单个节点
type Node struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Enable   bool   `json:"enable"`
	Region   string `json:"region"`
}

// Config 配置文件
type Config struct {
	Filename       string
	LastModifyTime int64
	Lock           *sync.RWMutex
	Data           []*Node
	C              chan bool
}

func NewConfig(filename string) {
	c := &Config{
		Filename: filename,
		Data:     make([]*Node, 0),
		Lock:     &sync.RWMutex{},
		C:        make(chan bool, 10),
	}
	c.parse()
	go c.reload()

	Conf = c
}

func (c *Config) parse() {
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
		return
	}
	c.Lock.Unlock()

	c.C <- true
}

func (c *Config) reload() {
	t := time.NewTicker(IvReloadConfig)
	for range t.C {
		fileInfo, _ := os.Stat(c.Filename)
		currModifyTime := fileInfo.ModTime().Unix()
		if currModifyTime > c.LastModifyTime {
			c.parse()
		}
	}
}

// Get 获取单个节点
func (c *Config) Get(key string) *Node {
	for _, node := range c.Data {
		if node.Id == key {
			return node
		}
	}
	return nil
}

// All 获取所有节点
func (c *Config) All() map[string]*Node {
	nodes := make(map[string]*Node, 0)
	for _, node := range c.Data {
		nodes[node.Id] = node
	}
	return nodes
}

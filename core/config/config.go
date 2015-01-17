package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"
)

type IConfig interface {
	GetConfig() interface{}
	loadConfig() interface{}
}

func NewConfig(configFile string, conf interface{}) config {
	return config{configFile, conf, nil}
}

type config struct {
	configFile string
	config     interface{}
	configData interface{}
}

func (c *config) GetConfig() interface{} {
	if reflect.TypeOf(c.configData) == nil {
		c.configData = c.loadConfig()
	}
	return c.configData
}

func (c *config) loadConfig() interface{} {
	log.Println("loading config from file...")
	file, err := ioutil.ReadFile(c.configFile)
	if err != nil {
		log.Println("open config: ", err)
	}
	temp := reflect.New(reflect.TypeOf(c.config).Elem()).Interface()
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config: ", err)
	}
	return temp
}

type ApiConfig struct {
	Hockeystream struct {
		Key string
	}
	Flightstats struct {
		AppId  string
		AppKey string
	}
	Openweathermap struct {
		AppId string
	}
}

var apiConfig config
var apiInit bool

func GetConfig() ApiConfig {
	if !apiInit {
		apiConfig = NewConfig("config.json", new(ApiConfig))
		apiInit = true
	}
	if r, ok := apiConfig.GetConfig().(*ApiConfig); ok {
		return *r
	} else {
		log.Println("boo!")
		var config ApiConfig
		return config
	}
}

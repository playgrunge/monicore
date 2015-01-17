package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"
)

func New(configFile string, conf interface{}) Config {
	return Config{configFile, conf, nil}
}

type Config struct {
	configFile string
	config     interface{}
	configData interface{}
}

func (c *Config) GetConfig() interface{} {
	if reflect.TypeOf(c.configData) == nil {
		c.configData = c.loadConfig()
	}
	return c.configData
}

func (c *Config) loadConfig() interface{} {
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

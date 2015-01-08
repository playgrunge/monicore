package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Hockeystream struct {
		Key string
	}
	Flightstats struct {
		AppId  string
		AppKey string
	}
}

var (
	c          Config
	configFile = "apiconfig.json"
)

func GetConfig() Config {
	if c == (Config{}) {
		loadConfig()
	}
	return c
}

func loadConfig() {
	log.Println("loading config from file...")
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("open config: ", err)
	}
	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config: ", err)
	}
	c = *temp
}

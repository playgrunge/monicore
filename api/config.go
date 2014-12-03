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

var c Config

func getConfig() Config {
	if c == (Config{}) {
		c = loadConfig()
	}
	return c
}

func loadConfig() Config {
	file, err := ioutil.ReadFile("apiconfig.json")
	if err != nil {
		log.Println("open config: ", err)
	}

	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config: ", err)
	}

	return *temp
}

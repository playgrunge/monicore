package control

import (
	"github.com/playgrunge/monicore/core/config"
)

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
	Twitter struct {
		ConsumerKey    string
		ConsumerSecret string
	}
}

var apiConfig config.Config
var apiInit bool

func GetConfig() ApiConfig {
	if !apiInit {
		apiConfig = config.New("config.json", new(ApiConfig))
		apiInit = true
	}
	if r, ok := apiConfig.GetConfig().(*ApiConfig); ok {
		return *r
	} else {
		var config ApiConfig
		return config
	}
}

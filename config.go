package binance_api

import (
	"github.com/jinzhu/configor"
	"path/filepath"
)

var BasePath = "."

var Config = struct {
	Binance struct {
		ApiKey    string
		SecretKey string
	}
}{}

func Initialize() {
	config := "config.yml"

	config = filepath.Join(BasePath, config)
	//fmt.Println("Test mode config", config)

	err := configor.New(&configor.Config{Silent: true, ErrorOnUnmatchedKeys: true}).Load(&Config, config)
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
}

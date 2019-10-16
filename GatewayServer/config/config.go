package config

import (
	"fmt"

	"github.com/robfig/config"
)

const configFile = "config.ini"

var listenPort string

func init() {
	conf, err := config.ReadDefault(configFile)
	if nil != err {
		fmt.Println("load config failed")
	} else {
		listenPort, _ = conf.String("base", "listenPort")
	}
}

//GetListenPort 本地监听端口
func GetListenPort() string {
	return listenPort
}

package conf

import (
	"log"
	"os"

	"code.google.com/p/gcfg"
)

type Configurer interface {
	MakeConfig(string)
	GetGottpConfig() *GottpSettings
}

func ReadConfig(configString string, cfg Configurer) error {
	err := gcfg.ReadStringInto(cfg, configString)
	if err != nil {
		log.Println("Error Loading configuration", err)
	}
	return err
}

func MakeConfig(configPath string, cfg Configurer) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("no such file or directory: " + configPath)
	}

	err := gcfg.ReadFileInto(cfg, configPath)
	if err != nil {
		panic(err)
	}
}

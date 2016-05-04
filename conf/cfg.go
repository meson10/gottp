package conf

import (
	"log"
	"os"

	"gopkg.in/gcfg.v1"
)

// Configurer interface implements two methods
type Configurer interface {
	MakeConfig(string)
	GetGottpConfig() *GottpSettings
}

//ReadConfig takes a configString, and a Configurer, it dumps the data from the string
//to corresponding feilds in the Configurer.
func ReadConfig(configString string, cfg Configurer) error {
	err := gcfg.ReadStringInto(cfg, configString)
	if err != nil {
		log.Println("Error Loading configuration", err)
	}
	return err
}

//MakeConfig takes a configPath and a Congigurer,  it dumps the data from the file
//to corresponding feilds in the Configurer.
func MakeConfig(configPath string, cfg Configurer) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("no such file or directory: " + configPath)
	}

	err := gcfg.ReadFileInto(cfg, configPath)
	if err != nil {
		panic(err)
	}
}

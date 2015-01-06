package main

import (
	"github.com/Simversity/gottp/conf"
)

type config struct {
	Gottp conf.GottpSettings
}

func (self *config) MakeConfig(configPath string) {
	if configPath != "" {
		conf.MakeConfig(configPath, self)
	}
}

func (self *config) GetGottpConfig() *conf.GottpSettings {
	return &self.Gottp
}

var settings config

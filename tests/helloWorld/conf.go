package main

import "github.com/Simversity/gottp"

type config struct {
	Gottp gottp.SettingsMap
}

func (self *config) MakeConfig(configPath string) {
	gottp.Settings = self.Gottp
}

var settings config

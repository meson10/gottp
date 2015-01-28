package main

import "gopkg.in/simversity/gottp.v0"

type config struct {
	Gottp gottp.SettingsMap
}

func (self *config) MakeConfig(configPath string) {
	gottp.Settings = self.Gottp
}

var settings config

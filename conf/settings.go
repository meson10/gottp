package conf

type SettingsMap struct {
	EmailHost     string
	EmailPort     string
	EmailUsername string
	EmailPassword string
	EmailSender   string
	EmailFrom     string
	ErrorTo       []string
	EmailDummy    bool
	Listen        string
}

const BaseConfig = `;Sample Configuration File
[gottp]
listen=""`

var Settings SettingsMap

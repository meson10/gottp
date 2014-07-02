package gottp

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

const baseConfig = `;Sample Configuration File
[gottp]
listen="127.0.0.1:8005";`

var Settings SettingsMap

package conf

type GottpSettings struct {
	EmailHost     string
	EmailPort     string
	EmailUsername string
	EmailPassword string
	EmailSender   string
	EmailFrom     string
	ErrorTo       string
	EmailDummy    bool
	Listen        string
}

const baseConfig = `;Sample Configuration File
[gottp]
listen="127.0.0.1:8005";`

type Config struct {
	Gottp GottpSettings
}

func (self *Config) MakeConfig(configPath string) {
	ReadConfig(baseConfig, self)
	if configPath != "" {
		MakeConfig(configPath, self)
	}
}

func (self *Config) GetGottpConfig() *GottpSettings {
	return &self.Gottp
}

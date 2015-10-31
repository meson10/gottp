package conf

//GottpSettings is a structure representatng all the setting, including listening address and port number
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

//Config is a structure that wraps GottpSettings for Configurations
//It implements both MakeConfig and GetConfig making it Configurer
type Config struct {
	Gottp GottpSettings
}

//MakeConfig takes the file path as configPath and returns with data filled
//into corresponding feilds of the Config struct.
//After a call to this function, Config.Gottp is populated with appropriate
//values.
func (c *Config) MakeConfig(configPath string) {
	ReadConfig(baseConfig, c)
	if configPath != "" {
		MakeConfig(configPath, c)
	}
}

//GetGottpConfig returns pointer to GottpSettings
func (c *Config) GetGottpConfig() *GottpSettings {
	return &c.Gottp
}

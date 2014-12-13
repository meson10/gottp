package conf

import (
	"flag"
	"log"
)

func CliArgs() (string, string) {
	var unixSocketptr = flag.String(
		"UNIX_SOCKET",
		"",
		"Use Unix Socket, default is None",
	)

	var config = flag.String(
		"config",
		"",
		"Config [.ini format] file to Load the configurations from",
	)

	flag.Parse()

	if *config == "" {
		log.Println("No config file supplied. Using defauls.")
	}

	return *config, *unixSocketptr
}

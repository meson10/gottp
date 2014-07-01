// +build !appengine

package gottp

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func cleanAddr(addr string) {
	err := os.Remove(addr)
	if err != nil {
		log.Panic(err)
	}
}

func interrupt_cleanup(addr string) {
	if strings.Index(addr, "/") != 0 {
		return
	}

	sigchan := make(chan os.Signal, 10)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	//NOTE: Capture every Signal right now.
	//signal.Notify(sigchan)

	s := <-sigchan
	log.Println("Exiting Program. Got Signal: ", s)

	// do last actions and wait for all write operations to end
	cleanAddr(addr)
	os.Exit(0)
}

func FlagArgs(cfg Configurer) map[string]*string {
	args := map[string]*string{}

	unixSocketptr := flag.String("UNIX_SOCKET", "", "Use Unix Socket, default is None")
	args["unix_socket"] = unixSocketptr

	configPtr := flag.String("config", "", "Config [.ini format] file to Load the configurations from")
	args["config"] = configPtr

	//Must be called after all flags are defined and before flags are accessed by the program.
	flag.Parse()

	cfg.MakeConfig(*args["config"])

	return args
}

type Configurer interface {
	MakeConfig(string)
	GetAddr() string
}

var SysInitChan = make(chan bool, 1)

func MakeServer(cfg Configurer) {
	var addr string
	ret := FlagArgs(cfg)

	if *ret["unix_socket"] != "" {
		addr = *ret["unix_socket"]
	} else {
		addr = cfg.GetAddr()
	}

	SysInitChan <- true

	var serverError error
	log.Println("Listening on " + addr)
	if strings.Index(addr, "/") == 0 {
		listener, err := net.Listen("unix", addr)
		if err != nil {
			c, err := net.Dial("unix", addr)

			if c != nil {
				defer c.Close()
			}

			if err != nil {
				log.Println("The socket does not look like consumed. Erase ?")
				cleanAddr(addr)
				listener, err = net.Listen("unix", addr)
			} else {
				log.Fatal("Cannot start Server. Address is already in Use.", err)
				os.Exit(0)
			}
		}

		go interrupt_cleanup(addr)
		serverError = http.Serve(listener, nil)
	} else {
		serverError = http.ListenAndServe(addr, nil)
	}

	if serverError != nil {
		log.Println(serverError)
	}
}

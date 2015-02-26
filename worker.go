package gottp

import (
	"gopkg.in/simversity/gotracer.v1"
	"log"
	"sync"
	"time"
)

var worker func(chan bool)

var errChan = make(chan bool)

var exitChan = make(chan bool, 1)

var wg = new(sync.WaitGroup)

func spawner() {
	go workerWrapper()

	s := <-errChan
	if s {
		go spawner()
	}
}

func workerWrapper() {

	wg.Add(1)

	defer wg.Done()
	defer gotracer.Tracer{
		Dummy:         settings.Gottp.EmailDummy,
		EmailHost:     settings.Gottp.EmailHost,
		EmailPort:     settings.Gottp.EmailPort,
		EmailPassword: settings.Gottp.EmailPassword,
		EmailUsername: settings.Gottp.EmailUsername,
		EmailSender:   settings.Gottp.EmailSender,
		EmailFrom:     settings.Gottp.EmailFrom,
		ErrorTo:       settings.Gottp.ErrorTo,
	}.Notify(func() string {
		errChan <- true
		return "Exception in worker"
	})

	worker(exitChan)
	errChan <- false
}

func RunWorker(wk func(chan bool)) {
	if worker != nil {
		panic("Worker already running.")
	}
	worker = wk
	go spawner()
}

func waitForWorker(compl chan<- bool) {
	wg.Wait()
	compl <- true
}

func StopWorker() {
	if worker != nil {
		waitChannel := make(chan bool)
		go waitForWorker(waitChannel)
		exitChan <- true

		select {
		case <-waitChannel:
		case <-time.After(10 * time.Second):
			log.Println("Timing out")

		}
		worker = nil
	}
}

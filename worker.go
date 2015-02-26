package gottp

import (
	"gopkg.in/simversity/gotracer.v1"
	"sync"
)

var workerRunning bool

var worker func(chan bool, chan bool)

var errChan = make(chan bool)

var exitChan = make(chan bool)

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
	defer gotracer.Tracer{Dummy: true}.Notify(func() string {
		errChan <- true
		return "Exception in worker"
	})

	worker(errChan, exitChan)
	errChan <- false
}

func RunWorker(wk func(chan bool, chan bool)) {
	if workerRunning {
		panic("Worker already running.")
	}
	worker = wk
	workerRunning = true
	go spawner()
}

func StopWorker() {
	if workerRunning {
		exitChan <- true
		wg.Wait()
		workerRunning = false
	}
}

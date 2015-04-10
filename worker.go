package gottp

import (
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

	_, ok := <-errChan
	if !ok {
		log.Println("Timing out in 10 seconds")
		time.Sleep(10 * time.Second)
		wg.Done()
		return
	}

	go spawner()
}

func workerWrapper() {
	wg.Add(1)
	defer wg.Done()

	// Trusing the fact that Tracer will always execute callback.
	defer Tracer.Notify(func(reason string) string {
		errChan <- true
		return "Exception in worker: " + reason
	})

	worker(exitChan)
}

func RunWorker(wk func(chan bool)) {
	if worker != nil {
		panic("Worker already running.")
	}
	worker = wk
	go spawner()
}

func shutdownWorker() {
	if worker == nil {
		return
	}

	log.Println("Preparing for shutdown")
	exitChan <- true
	close(errChan)
	wg.Wait()
}

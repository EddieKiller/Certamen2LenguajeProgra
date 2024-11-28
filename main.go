package main

import (
	"sync"
)

func main() {
	readyQueue := make(chan *Process, 10)
	blockedQueue := make(chan *Process, 10)
	cpu := make(chan *Process)
	output := make(chan string)

	dispatcher := &Dispatcher{
		ReadyQueue:   readyQueue,
		BlockedQueue: blockedQueue,
		CPU:          cpu,
		Output:       output,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go dispatcher.Dispatch(4, &wg)

	// Simula la creación de procesos y ejecución
	process1 := &Process{
		ID:            "Proceso_1",
		State:         "Listo",
		ProgramCounter: 0,
		Instructions:   []string{"I", "I", "ES 3", "I", "F"},
	}
	readyQueue <- process1

	wg.Wait()
}

package main

import (
	"fmt"
	"strings"
	"time"
	"sync"
)

type Dispatcher struct {
	ReadyQueue   chan *Process
	BlockedQueue chan *Process
	CPU          chan *Process
	Output       chan string
}

func (d *Dispatcher) Dispatch(m int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case process := <-d.ReadyQueue:
			d.CPU <- process
			d.Output <- fmt.Sprintf("LOAD %s", process.ID)
			// Simula la ejecución de m instrucciones
			for i := 0; i < m && process.ProgramCounter < len(process.Instructions); i++ {
				instruction := process.Instructions[process.ProgramCounter]
				d.Output <- fmt.Sprintf("EXEC %s: %s", process.ID, instruction)
				process.ProgramCounter++
				time.Sleep(100 * time.Millisecond) // Simula el tiempo de CPU
				if strings.HasPrefix(instruction, "ES") {
					process.IOState = extractIOTime(instruction)
					process.State = "Bloqueado"
					d.Output <- fmt.Sprintf("ST %s -> Bloqueado", process.ID)
					d.BlockedQueue <- process
					break
				}
			}

			if process.ProgramCounter >= len(process.Instructions) {
				process.State = "Terminado"
				d.Output <- fmt.Sprintf("F %s", process.ID)
			} else if process.State != "Bloqueado" {
				process.State = "Listo"
				d.ReadyQueue <- process
			}
		}
	}
}

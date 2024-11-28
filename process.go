package main

import (
	"strconv"
	"strings"
)

type Process struct {
	ID            string
	State         string // "Listo", "Bloqueado", "Ejecutando", "Terminado"
	ProgramCounter int
	IOState        int // Tiempo restante para desbloquear (si aplica)
	Instructions   []string
}

func extractIOTime(instruction string) int {
	parts := strings.Split(instruction, " ")
	if len(parts) < 2 {
		return 0
	}
	time, _ := strconv.Atoi(parts[1])
	return time
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var AsignarID int = 1 // Identifiant du processus
var currentTime int = 1
var executionLog []string

func AddToLog(time int, instruction, component string, pc int) {
	logEntry := fmt.Sprintf("%d\t%s\t%s\t%d", time, instruction, component, pc)
	executionLog = append(executionLog, logEntry)
}

// Structure représentant un ordre de création
type CreationOrder struct {
	Time  int      // Temps de création
	Files []string // Liste des fichiers de processus
}

func LoadCreationOrder(filename string) ([]CreationOrder, error) {
	fmt.Println("Lecture du fichier d'ordre :", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de l'ouverture du fichier : %w", err)
	}
	defer file.Close()

	var orders []CreationOrder
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fmt.Printf("Ligne brute lue : '%s'\n", line) // Débogage

		// Ignorer les lignes vides ou les commentaires
		if line == "" || strings.HasPrefix(line, "#") {
			fmt.Println("Ligne ignorée :", line)
			continue
		}

		// Séparer le temps de création et les fichiers
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Println("Ligne mal formatée, ignorée :", line)
			continue
		}

		time, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Printf("Erreur de conversion du temps : '%s'\n", parts[0])
			continue
		}

		files := parts[1:]                                     // Le reste sont les fichiers de processus
		fmt.Printf("Temps : %d, Fichiers : %v\n", time, files) // Débogage
		orders = append(orders, CreationOrder{Time: time, Files: files})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Erreur lors de la lecture du fichier : %w", err)
	}

	fmt.Println("Ordres chargés :", orders) // Résultat final
	return orders, nil
}

// Gérer la création des processus en fonction du fichier d'ordre
func (d *Dispatcher) HandleCreationOrders(orders []CreationOrder, currentTime int) {
	for _, order := range orders {
		if order.Time == currentTime {
			for _, file := range order.Files {
				process, err := LoadProcessFile(AsignarID, file)
				if err != nil {
					fmt.Printf("Erreur lors du chargement du fichier de processus %s : %s\n", file, err)
					continue
				}
				d.AddToReadyQueue(process)
				AddToLog(currentTime, fmt.Sprintf("LOAD %s", file), "Dispatcher", process.ProgramCounter)
				AsignarID++
				currentTime++
			}
		}
	}
}

// Dispatcher structure
type Dispatcher struct {
	ReadyQueue   []Process // File des processus prêts
	BlockedQueue []Process // File des processus bloqués
}

// Créer un nouveau dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		ReadyQueue:   []Process{},
		BlockedQueue: []Process{},
	}
}

// Ajouter un processus à la file Ready
func (d *Dispatcher) AddToReadyQueue(process Process) {
	fmt.Printf("Ajout du processus %d à la file Ready\n", process.ID)
	d.ReadyQueue = append(d.ReadyQueue, process)
}

// Ajouter un processus à la file Blocked
func (d *Dispatcher) AddToBlockedQueue(process Process) {
	fmt.Printf("Ajout du processus %d à la file Blocked\n", process.ID)
	d.BlockedQueue = append(d.BlockedQueue, process)
}

// Retirer un processus de la file Ready
func (d *Dispatcher) PullFromReadyQueue() (Process, bool) {
	if len(d.ReadyQueue) == 0 {
		return Process{}, false
	}
	process := d.ReadyQueue[0]
	d.ReadyQueue = d.ReadyQueue[1:]

	// Log du PULL
	AddToLog(currentTime, "PULL", "Dispatcher", process.ProgramCounter)
	currentTime++

	return process, true
}

// Déplacer les processus débloqués de BlockedQueue vers ReadyQueue
func (d *Dispatcher) PullFromBlockedQueue() {
	for i := 0; i < len(d.BlockedQueue); {
		process := &d.BlockedQueue[i]
		if process.IOState > 0 {
			process.IOState-- // Réduire le temps de blocage
			fmt.Printf("Processus %d: Temps de blocage restant %d\n", process.ID, process.IOState)
		}
		if process.IOState == 0 {
			fmt.Printf("Processus %d débloqué, déplacement vers ReadyQueue\n", process.ID)
			d.AddToReadyQueue(*process)
			// Supprimer le processus débloqué de BlockedQueue
			d.BlockedQueue = append(d.BlockedQueue[:i], d.BlockedQueue[i+1:]...)
		} else {
			i++
		}
	}
}

// Gérer les transitions des processus bloqués
func (d *Dispatcher) HandleBlockedProcesses() {
	for i := 0; i < len(d.BlockedQueue); {
		process := &d.BlockedQueue[i]
		if process.IOState > 0 {
			process.IOState-- // Réduire le temps de blocage
			if process.IOState == 0 {
				fmt.Printf("Le processus %d est débloqué\n", process.ID)
				d.AddToReadyQueue(*process)
				d.BlockedQueue = append(d.BlockedQueue[:i], d.BlockedQueue[i+1:]...)
				continue
			}
		}
		i++
	}
}

func (d *Dispatcher) ExecuteProcesses(cycles int, orders []CreationOrder) {
	for cycles > 0 {
		fmt.Println("Cycle de processeur :", cycles)

		// Gérer les ordres de création à l'instant courant
		d.HandleCreationOrders(orders, currentTime)

		// Gérer les processus bloqués
		d.PullFromBlockedQueue()

		// Retirer un processus de la file Ready
		process, ok := d.PullFromReadyQueue()
		if !ok {
			fmt.Println("Aucun processus prêt à exécuter")
			cycles--
			currentTime++
			continue
		}

		// Exécuter le processus
		fmt.Printf("Exécution du processus %d\n", process.ID)

		AddToLog(currentTime, "EXEC", fmt.Sprintf("nombre_proceso_%d", process.ID), process.ProgramCounter)
		currentTime++
		for cycles > 0 && process.ProgramCounter < len(process.Instructions) {
			instruction := process.Instructions[process.ProgramCounter]
			fmt.Printf("Instruction exécutée : %s\n", instruction)
			process.ProgramCounter++

			// Ajouter l'instruction au log
			AddToLog(currentTime, instruction, fmt.Sprintf("nombre_proceso_%d", process.ID), process.ProgramCounter)

			if instruction == "F" {
				fmt.Printf("Le processus %d est terminé\n", process.ID)
				AddToLog(currentTime, "END", fmt.Sprintf("nombre_proceso_%d", process.ID), process.ProgramCounter)
				cycles--

				break
			} else if len(instruction) >= 2 && instruction[:2] == "ES" {
				ioState := extractDelay(instruction)
				fmt.Printf("Le processus %d est bloqué pour %d cycles\n", process.ID, ioState)
				process.IOState = ioState
				d.AddToBlockedQueue(process)
				AddToLog(currentTime, fmt.Sprintf("PUSH_Bloqueado nombre_proceso_%d", process.ID), "Dispatcher", process.ProgramCounter)
				currentTime++
				break
			}

			cycles--
			currentTime++
		}

		// Ajouter le processus de retour à Ready si non terminé
		if process.ProgramCounter < len(process.Instructions) && process.IOState == 0 {
			d.AddToReadyQueue(process)
			AddToLog(currentTime, fmt.Sprintf("PUSH_Listo nombre_proceso_%d", process.ID), "Dispatcher", process.ProgramCounter)
			currentTime++
		}

		// Incrémenter le temps courant
		currentTime++
	}
}

// Extraire le délai d'une instruction ES
func extractDelay(instruction string) int {
	var delay int
	if len(instruction) > 3 {
		fmt.Sscanf(instruction, "ES %d", &delay)
	}
	return delay
}

func WriteLogToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier de log :", err)
		return
	}
	defer file.Close()

	// Ajouter un en-tête
	_, _ = file.WriteString("# Tiempo de CPU\tTipo Instrucción\tProceso/Despachador\tValor CP\n")

	// Écrire chaque log
	for _, logEntry := range executionLog {
		_, _ = file.WriteString(logEntry + "\n")
	}

	fmt.Println("Fichier de log généré :", filename)
}
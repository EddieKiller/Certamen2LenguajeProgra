package main

import (
	"fmt"
)

func main() {
	// Charger le fichier d'ordre de création
	orders, err := LoadCreationOrder("order.txt")
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	// Créer le dispatcher
	dispatcher := NewDispatcher()

	// Exécuter les processus avec les ordres de création
	for {

		dispatcher.ExecuteProcesses(5, orders)

		if len(dispatcher.ReadyQueue) == 0 && len(dispatcher.BlockedQueue) == 0 {
			fmt.Println("Toutes les files sont vides. Fin de la simulation.")
			break
		}
	}

	WriteLogToFile("execution_log.txt")

}

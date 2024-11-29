package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Structure d'un processus
type Process struct {
	ID             int      // Identifiant du processus
	State          string   // État du processus : "Ready", "Blocked", etc.
	ProgramCounter int      // Ligne actuelle d'exécution
	IOState        int      // Temps restant pour débloquer (si applicable)
	Instructions   []string // Liste des instructions
}

// Créer un nouveau processus
func NewProcess(id int, instructions []string) Process {
	return Process{
		ID:             id, // Pour simplifier, on fixe l'ID à 1 (modifiable selon vos besoins)
		State:          "Ready",
		ProgramCounter: 0,
		IOState:        0,
		Instructions:   instructions,
	}
}

func LoadProcessFile(id int, filename string) (Process, error) {
	fmt.Println("Lecture du fichier :", filename)
	file, err := os.Open(filename)
	if err != nil {
		return Process{}, fmt.Errorf("Erreur lors de l'ouverture du fichier : %w", err)
	}
	defer file.Close()

	var instructions []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Printf("Ligne brute lue : '%s'\n", line)

		// Nettoyer la ligne
		line = strings.TrimSpace(line)
		//fmt.Printf("Ligne nettoyée : '%s'\n", line)

		if line == "" {
			//fmt.Println("Ligne vide ignorée")
			continue
		}

		instructions = append(instructions, line)
		//fmt.Println("Instruction ajoutée :", line)
	}

	if err := scanner.Err(); err != nil {
		return Process{}, fmt.Errorf("Erreur lors de la lecture du fichier : %w", err)
	}

	//fmt.Println("Instructions lues :", instructions)
	return NewProcess(id, instructions), nil
}

// Fonction pour afficher les détails du processus
func PrintProcessDetails(process Process) {
	fmt.Printf("Processus ID: %d\n", process.ID)
	fmt.Printf("État: %s\n", process.State)
	fmt.Println("Liste des instructions:")
	for i, instruction := range process.Instructions {
		fmt.Printf("  Instruction %d: %s\n", i+1, instruction)
	}
}

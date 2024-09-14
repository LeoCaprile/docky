package main

import (
	"docky/models"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {

	program := tea.NewProgram(models.GetInitialModel())

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running programgetInitialModel", err)
		os.Exit(1)
	}

}

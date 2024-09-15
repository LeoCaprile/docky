package main

import (
	"docky/docker/client"
	"docky/models"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	client.TestConnection()

	program := tea.NewProgram(models.GetInitialModel())

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running programgetInitialModel", err)
		os.Exit(1)
	}

}

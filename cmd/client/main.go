package main

import (
	"fmt"
	"os"


	"github.com/Gjones747/goChat/internal/client/views"
	tea "github.com/charmbracelet/bubbletea"
)


func main() {
	fmt.Println("here")
	tui := tea.NewProgram(
		views.InitialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),


		)
	if _, err := tui.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
	
}

package views

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// this is the first view that a user is greated with it should show a nice box to input a room code that the user wants to join

type enterRoomView struct {
	textInput     textinput.Model
	err           error
	fieldComplete bool

	windowWidth  int
	windowHeight int
}

type roomCodeInput struct {
	input string
}

func (m enterRoomView) isComplete() tea.Cmd {
	return func() tea.Msg {
		return roomCodeInput{input: m.textInput.Value()}
	}
}

func initialEnterRoomView() enterRoomView {
	ti := textinput.New()
	ti.Placeholder = "roomcode"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return enterRoomView{
		textInput:     ti,
		err:           nil,
		fieldComplete: false,
	}
}

func (model enterRoomView) Init() tea.Cmd {
	return textinput.Blink
}

func (m enterRoomView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			log.Println("here")
			m.fieldComplete = true
			return m, m.isComplete()
		}

		// We handle errors just like any other message
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (model enterRoomView) View() string {
	return fmt.Sprintf(
		"Welcome to goChat! \n\n %s %s",
		"Please enter the roomcode for the room you want to join! \n",
		model.textInput.View()) + "\n"
}

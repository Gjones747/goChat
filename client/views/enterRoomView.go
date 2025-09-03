package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
			m.fieldComplete = true
			return m, m.isComplete()
		}

		// We handle errors just like any other message
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m enterRoomView) introText() string {
	boldStyle := lipgloss.NewStyle().Bold(true).Italic(true).Foreground(lipgloss.Color("5"))

	paddingBetweenSize := 2
	paddingTopFirst := ((m.windowHeight) / 6)


	title := `
 ______     ______     ______     __  __     ______     ______  
/\  ___\   /\  __ \   /\  ___\   /\ \/ /    /\  ___\   /\__  _\ 
\ \___  \  \ \ \/\ \  \ \ \____  \ \  _"-.  \ \  __\   \/_/\ \/ 
 \/\_____\  \ \_____\  \ \_____\  \ \_\ \_\  \ \_____\    \ \_\ 
  \/_____/   \/_____/   \/_____/   \/_/\/_/   \/_____/     \/_/ 
                                                                `
	
    paddingBetween := strings.Repeat("\n", paddingBetweenSize)
    paddingFirst := strings.Repeat("\n", paddingTopFirst)


	windowStyle := lipgloss.NewStyle().
        Align(lipgloss.Center).
        Width(m.windowWidth-5).
        Height(m.windowHeight-5).
        BorderStyle(lipgloss.ASCIIBorder())

	inputStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		BorderStyle(lipgloss.ASCIIBorder())

	return windowStyle.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		paddingFirst,
		title,
		paddingBetween,
		boldStyle.Render("Welcome to Socket, a lightweight terminal based chat application!"),
		"To begin enter the roomcode for the chat room you want to join! \n",
		inputStyle.Render(m.textInput.View())))

}

func (model enterRoomView) View() string {
	return model.introText()
}










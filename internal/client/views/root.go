package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

// this is the root or main view of the cli chat interface
type rootModel struct {
	currentPage   int
	enterRoomView enterRoomView
}

func InitialModel() rootModel {
	view := rootModel{
		currentPage:   0,
		enterRoomView: initEnterRoomView(),
	}

	return view
}

func (model rootModel) Init() tea.Cmd {
	return nil
}

func (model rootModel) View() string {
	switch {
	case model.currentPage == 0:
		return model.enterRoomView.View()
	default:
		return ""
	}
}

func (model rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return model, tea.Quit

		}

	}

	return model, nil
}

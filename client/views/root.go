package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

// this is the root or main view of the cli chat interface
type rootModel struct {
	currentPage   int

	roomView roomView
	enterRoomView enterRoomView
}

func InitialModel() rootModel {
	view := rootModel{
		currentPage:   1,

		roomView: initialRoomView(),
		enterRoomView: initialEnterRoomView(),
	}

	return view
}

func (model rootModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (model rootModel) View() string {
	switch {
	case model.currentPage == 0:
		return model.enterRoomView.View()
	case model.currentPage == 1:
		return model.roomView.View()
	default:
		return ""
	}
}

func (model rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch model.currentPage {
	case 0:
		newView, cmd := model.enterRoomView.Update(msg)
		return newView, cmd
	case 1:
		newView, cmd := model.roomView.Update(msg)
		return newView, cmd
	}

	return model, nil
}

package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

// this is the root or main view of the cli chat interface
type rootModel struct {
	currentPage int

	roomView      roomView
	enterRoomView enterRoomView
}

func InitialModel() rootModel {
	view := rootModel{
		currentPage: 0,

		roomView:      initialRoomView(),
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

type initializeWindow struct{}

func (model rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.roomView.windowWidth = msg.Width
		model.roomView.windowHeight = msg.Height
		model.enterRoomView.windowWidth = msg.Width
		model.enterRoomView.windowHeight = msg.Height
	case roomCodeInput:
		model.currentPage = 1
		model.roomView.roomCode = msg.input
		return model, func() tea.Msg {
			return initializeWindow{}
		}
	}

	switch model.currentPage {
	case 0:
		newView, cmd := model.enterRoomView.Update(msg)
		model.enterRoomView = newView.(enterRoomView)
		return model, cmd
	case 1:
		newView, cmd := model.roomView.Update(msg)
		model.roomView = newView.(roomView)
		return model, cmd
	}

	return model, nil
}

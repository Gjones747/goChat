package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


const gap = "\n\n"

			
type roomView struct {
	messages []string

	viewport viewport.Model
	textInput textinput.Model
	senderStyle lipgloss.Style

	ready bool
	
	windowWidth int 
	windowHwight int
}


func initialRoomView() roomView {
	ti := textinput.New()
	ti.Width = 100
	ti.Placeholder = " \x1b[3menter your message here        \x1b[0m "
	ti.Focus()
	ti.CharLimit = 500

	vp := viewport.New(40, 5)

	return roomView {
		messages: []string{},
		viewport: vp,
		textInput: ti,
	}
}


func (model roomView) Init() tea.Cmd {
	return textinput.Blink
}


func (m roomView) Update (msg tea.Msg) (tea.Model, tea.Cmd) {

    var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)


	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "j" || keyMsg.String() == "k" {
			m.textInput, tiCmd = m.textInput.Update(msg)
		} else {
			m.textInput, tiCmd = m.textInput.Update(msg)
			m.viewport, vpCmd = m.viewport.Update(msg)
		}

	} else {
		m.textInput, tiCmd = m.textInput.Update(msg)
		m.viewport, vpCmd = m.viewport.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textInput.Width = msg.Width
		m.viewport.Height = msg.Height - lipgloss.Height(m.headerView()) - lipgloss.Height(gap) - 2

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textInput.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textInput.Value())
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textInput.Reset(    )
			m.viewport.GotoBottom()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}


func (m roomView) View() string  {
	return fmt.Sprintf(
		"%s%s%s%s%s%s",
		m.headerView(),
		m.viewport.View(),
		m.renderFooter(),
		"\n\n",
		m.textInput.View(),
		"\n\n",
	)
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)


func (m roomView) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}


func (m roomView) renderFooter() string {
	info := "\n"
	line := strings.Repeat("─", max(0, m.viewport.Width))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

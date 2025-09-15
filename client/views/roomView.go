package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"os/user"
	"strings"

	"github.com/Gjones747/goChat/api"
	socketcontroller "github.com/Gjones747/goChat/webSocket"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

type roomView struct {
	messages []string
	users    []string

	viewport    viewport.Model
	textInput   textinput.Model
	senderStyle lipgloss.Style

	userName string

	ready    bool
	roomCode string

	windowWidth  int
	windowHeight int

	connection net.Conn

	incomingTeaMsg chan incomingMessage
	incomingUsers  chan []byte
}

var incomingMessages chan api.Message = make(chan api.Message, 10)

//var (
//	statusMessageStyle = lipgloss.NewStyle().
//	Foreground(lipgloss.Color("5")).
//	Bold(true).
//	Italic(true)
//)

type incomingMessage struct {
	message string
}

func (m *roomView) startIncomingRelay() {
	go func() {
		for msg := range incomingMessages { // incomingMessages is your package channel of api.Message
			formatted := fmt.Sprintf("%s-%s: %s", msg.DateTime, msg.SenderName, msg.Contents)
			m.incomingTeaMsg <- incomingMessage{message: formatted}
		}
		// When incomingMessages is closed, this goroutine exits cleanly
	}()
}

func (m *roomView) watchIncoming() tea.Cmd {
	return func() tea.Msg {
		return <-m.incomingTeaMsg // blocks until relay pushes a formatted incomingMessage
	}
}

func (m *roomView) sendMessage(contents string) {
	newMessage := api.NewMessage(m.userName, []byte(contents))

	socketcontroller.SendFramesToServer(m.connection, newMessage)
}

func (m *roomView) watchConnection() {
	for {
		data, err := socketcontroller.ClientFrameReader(m.connection)

		if err != nil {
			log.Println("connection error")
			return
		}

		var dataEnvelope api.Envelope
		if err := json.Unmarshal(data, &dataEnvelope); err != nil {
			log.Println("failed to decode server response Yikes!")
			return
		}

		switch dataEnvelope.DataType {
		case "message":
			var msgData api.Message
			if err := json.Unmarshal(dataEnvelope.Data, &msgData); err != nil {
				log.Println("failed to parse message json")
				return
			}
			incomingMessages <- msgData

		}
	}

}

func initialRoomView() roomView {
	ti := textinput.New()
	ti.Width = 100
	ti.Placeholder = " \x1b[3menter your message here        \x1b[0m "
	ti.Focus()
	ti.CharLimit = 500

	vp := viewport.New(40, 5)

	return roomView{
		messages:  []string{},
		viewport:  vp,
		textInput: ti,
		users:     []string{},

		incomingTeaMsg: make(chan incomingMessage, 10),
		incomingUsers:  make(chan []byte, 10),
	}
}

func (model roomView) Init() tea.Cmd {
	return nil
}

func (m roomView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
	case initializeWindow:

		currentUser, err := user.Current()
		if err != nil {
		}

		m.userName = currentUser.Username

		connection, err := m.initializeConnection()

		if err != nil {
			// todo RETURN ERROR SCREEN VIEW or go back to join room or smth just show an error ahh handling
		}

		m.connection = connection

		m.messages = append(m.messages, fmt.Sprintf("Socket Connection Established! Hostname: %s just joined: %s", m.userName, m.roomCode))

		go m.watchConnection()
		m.startIncomingRelay()

		m.viewport.Width = m.windowWidth
		m.textInput.Width = m.windowWidth
		m.viewport.Height = m.windowHeight - lipgloss.Height(m.headerView()) - lipgloss.Height(gap) - 3

		if len(m.messages) == 1 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Foreground(lipgloss.Color("5")).Bold(true).Italic(true).Render(strings.Join(m.messages, "\n")))
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		return m, m.watchIncoming()

	case incomingMessage:
		m.messages = append(m.messages, msg.message)
		// update viewport content and scroll
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
		return m, m.watchIncoming()

	case tea.WindowSizeMsg:
		m.viewport.Width = m.windowWidth
		m.textInput.Width = m.windowWidth
		m.viewport.Height = m.windowHeight - lipgloss.Height(m.headerView()) - lipgloss.Height(gap) - 3

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
			userInput := m.textInput.Value()
			m.sendMessage(userInput)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textInput.Reset()
			m.viewport.GotoBottom()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m roomView) View() string {
	return lipgloss.PlaceHorizontal(m.windowWidth, lipgloss.Center, m.renderMessageViewPort())
}

func (m roomView) renderMessageViewPort() string {

	borderStyle := lipgloss.ASCIIBorder()

	windowStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(m.windowWidth - 5).
		Height(m.windowHeight - 5).
		BorderTop(true).BorderBottom(true).BorderRight(true).BorderLeft(true).
		BorderStyle(borderStyle).
		Padding(1)

	display := fmt.Sprintf(
		"%s%s%s%s",
		m.viewport.View(),
		m.renderFooter(),
		"\n\n",
		m.textInput.View())

	return windowStyle.Render(display)

}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.ASCIIBorder()
		b.Right = "|"
		b.Left = "|"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

func (m roomView) headerView() string {
	title := titleStyle.Render(m.roomCode)
	return lipgloss.JoinHorizontal(lipgloss.Center, title)
}

func (m roomView) renderFooter() string {
	info := "\n"
	line := strings.Repeat("-", max(0, m.viewport.Width-7))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// TODO dot env
func (m roomView) initializeConnection() (net.Conn, error) {

	url := url.URL{Scheme: "ws", Host: "kaolun.site:8081", Path: "/joinRoom"}
	q := url.Query()
	q.Add("room_code", m.roomCode)
	q.Add("user_name", m.userName)

	url.RawQuery = q.Encode()

	connection, err := socketcontroller.SendUpgrade(url)

	if err != nil {
		return nil, err
	}

	return connection, nil

}

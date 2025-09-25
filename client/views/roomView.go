package views

import (
	"bytes"
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

	viewport         viewport.Model
	userListViewPort viewport.Model
	textInput        textinput.Model
	senderStyle      lipgloss.Style

	userName  string
	sessionID []byte
	stringSessionID string

	ready    bool
	roomCode string

	windowWidth  int
	windowHeight int

	connection net.Conn

	incomingTeaMsg      chan incomingMessage
	incomingUserListMsg chan incomingUserListMessage
}

var (
	boldOrangeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))

	boldStyle = lipgloss.NewStyle().
			Bold(true)

	roomCodeStyle = lipgloss.NewStyle().Bold(true).Italic(true).Foreground(lipgloss.Color("5"))

	youSentStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	otherSendStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
)

var incomingMessages chan api.Message = make(chan api.Message, 10)
var incomingUserList chan api.UserList = make(chan api.UserList, 10)

//var (
//	statusMessageStyle = lipgloss.NewStyle().
//	Foreground(lipgloss.Color("5")).
//	Bold(true).
//	Italic(true)
//)

type incomingMessage struct {
	message string
}

type incomingUserListMessage struct {
	users []string
}

func (m *roomView) startIncomingRelay() {
	go func() {
		for msg := range incomingMessages { // incomingMessages is your package channel of api.Message
			if msg.Type == "USER" {
				var header string
				if bytes.Equal(msg.SessionID, m.sessionID) {
					name := youSentStyle.Render("You sent: ")
					header = fmt.Sprintf("%s %s", msg.DateTime.Format("3:04pm"), name)
				} else {
					name := otherSendStyle.Render(msg.SenderName)
					header = fmt.Sprintf("%s %s: ", msg.DateTime.Format("3:04pm"), name)
				}
				formatted := header + string(msg.Contents)
				m.incomingTeaMsg <- incomingMessage{message: formatted}
			}

			if msg.Type == "JOIN_LEAVE" {
				formatted := fmt.Sprintf("%s %s", msg.DateTime.Format("3:04pm"), msg.Contents)
				m.incomingTeaMsg <- incomingMessage{message: formatted}	
			}

		}
		// When incomingMessages is closed, this goroutine exits cleanly
	}()
}
func (m *roomView) startIncomingUserListRelay() {
	go func() {
		for msg := range incomingUserList { // incomingMessages is your package channel of api.Message
			var formattedUsers []string
			for _, val := range msg.Users {
				if m.stringSessionID == val[1]{
					formattedUsers = append(formattedUsers, youSentStyle.Render(val[0]))
				} else {
					formattedUsers = append(formattedUsers, val[0])
				}
			}
			m.incomingUserListMsg <- incomingUserListMessage{users: formattedUsers}
		}
		// When incomingMessages is closed, this goroutine exits cleanly
	}()
}

func (m *roomView) watchIncoming() tea.Cmd {
	return func() tea.Msg {
		return <-m.incomingTeaMsg // blocks until relay pushes a formatted incomingMessage
	}
}

func (m *roomView) watchIncomingUserList() tea.Cmd {
	return func() tea.Msg {
		return <-m.incomingUserListMsg
	}
}

func (m *roomView) sendMessage(contents string) {
	newMessage := api.NewMessage("USER", m.userName, m.sessionID, []byte(contents))

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

		case "userList":
			var userListData api.UserList
			if err := json.Unmarshal(dataEnvelope.Data, &userListData); err != nil {
				log.Println("failed to parse userList json")
				return
			}
			incomingUserList <- userListData

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

	userListVP := viewport.New(40, 5)
	uniqueSessionID, err := api.MakeSessionID()
	stringSessionID := bytesToBinaryString(uniqueSessionID)

	if err != nil {
		log.Println(err)
	}

	return roomView{
		messages:         []string{},
		userListViewPort: userListVP,
		viewport:         vp,
		textInput:        ti,
		users:            []string{},

		sessionID: uniqueSessionID,
		stringSessionID: stringSessionID,

		incomingTeaMsg:      make(chan incomingMessage, 10),
		incomingUserListMsg: make(chan incomingUserListMessage, 10),
	}
}

func (model roomView) Init() tea.Cmd {
	return nil
}

func (m roomView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		tiCmd     tea.Cmd
		vpCmd     tea.Cmd
		userVPCmd tea.Cmd
	)

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "j" || keyMsg.String() == "k" {
			m.textInput, tiCmd = m.textInput.Update(msg)
		} else {
			m.textInput, tiCmd = m.textInput.Update(msg)
			m.viewport, vpCmd = m.viewport.Update(msg)
			m.userListViewPort, userVPCmd = m.userListViewPort.Update(msg)
		}

	} else {
		m.textInput, tiCmd = m.textInput.Update(msg)
		m.viewport, vpCmd = m.viewport.Update(msg)
		m.userListViewPort, userVPCmd = m.userListViewPort.Update(msg)
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
		m.startIncomingUserListRelay()

		m.viewport.Width = m.windowWidth
		m.textInput.Width = m.windowWidth
		m.viewport.Height = m.windowHeight - lipgloss.Height(m.headerView()) - lipgloss.Height(gap) - 3

		if len(m.messages) == 1 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Foreground(lipgloss.Color("5")).Bold(true).Italic(true).Render(strings.Join(m.messages, "\n")))
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		return m, tea.Batch(m.watchIncoming(), m.watchIncomingUserList())

	case incomingMessage:
		m.messages = append(m.messages, msg.message)
		// update viewport content and scroll
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
		return m, m.watchIncoming()

	case incomingUserListMessage:
		m.users = msg.users
		m.userListViewPort.SetContent(lipgloss.NewStyle().Width(m.userListViewPort.Width).Render(strings.Join(msg.users, "\n")))
		m.userListViewPort.GotoBottom()
		return m, m.watchIncomingUserList()

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

	return m, tea.Batch(tiCmd, vpCmd, userVPCmd)
}

func (m roomView) renderUserList(width, height int) string {
	// Room info at the top
	roomInfoStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.ASCIIBorder()).
		PaddingLeft(2).
		Width(width)
	commandsHeaderText := roomCodeStyle.Render("\nCOMMANDS:")
	quitText := boldOrangeStyle.Render(`"q"`) + boldStyle.Render(" or ") + boldOrangeStyle.Render(`"Ctrl+c"`) + boldStyle.Render(" to quit")
	roomCodeText := roomCodeStyle.Render("\nROOM CODE: ") + boldStyle.Render(m.roomCode)
	title := lipgloss.JoinVertical(lipgloss.Left, boldOrangeStyle.Render(`
   ____         __       __ 
  / __/__  ____/ /_____ / /_
 _\ \/ _ \/ __/  '_/ -_) __/
/___/\___/\__/_/\_\\__/\__/ `), roomCodeText, commandsHeaderText, quitText, "\n")
	var renderedInfo string
	if width-2 >= lipgloss.Width(title) {
		renderedInfo = roomInfoStyle.Render(title)
	} else {
		title = lipgloss.JoinVertical(lipgloss.Left, roomCodeText, commandsHeaderText, quitText, "\n")
		renderedInfo = roomInfoStyle.Render(title)

	}

	// User list viewport below the room info
	userListHeight := height - lipgloss.Height(renderedInfo)
	userListStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.ASCIIBorder()).
		Width(width).
		Height(userListHeight)

	renderedList := userListStyle.Render(m.userListViewPort.View())

	// Combine room info + user list vertically
	return lipgloss.JoinVertical(lipgloss.Left, renderedInfo, renderedList)
}

func (m roomView) View() string {
	// Compute split sizes
	messageWidth := int(float64(m.windowWidth) * 0.8)
	userListWidth := m.windowWidth - messageWidth
	contentHeight := m.windowHeight - 5

	// Build main content row (messages + user list)
	contentRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderMessageViewPort(messageWidth, contentHeight),
		m.renderUserList(userListWidth, contentHeight+1),
	)
	return lipgloss.PlaceHorizontal(m.windowWidth, lipgloss.Center, contentRow)
}

func (m roomView) renderMessageViewPort(width, height int) string {

	width = width - 7
	borderStyle := lipgloss.ASCIIBorder()

	windowStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(width).
		Height(height).
		BorderTop(true).BorderBottom(true).BorderRight(true).BorderLeft(true).
		BorderStyle(borderStyle).
		Padding(1)

	display := fmt.Sprintf(
		"%s%s%s%s",
		m.viewport.View(),
		m.renderFooter(width),
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

func (m roomView) renderFooter(width int) string {
	info := "\n"
	line := strings.Repeat("-", max(0, width-2))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// TODO dot env
func (m roomView) initializeConnection() (net.Conn, error) {

	url := url.URL{Scheme: "wss", Host: "chat.kaolun.site:8081", Path: "/joinRoom"}
	q := url.Query()
	q.Add("room_code", m.roomCode)
	q.Add("user_name", m.userName)
	q.Add("session_id", string(m.stringSessionID))

	url.RawQuery = q.Encode()

	connection, err := socketcontroller.SendUpgrade(url)

	if err != nil {
		return nil, err
	}

	return connection, nil

}




func bytesToBinaryString(data []byte) string {
	var sb strings.Builder // Use strings.Builder for efficient string concatenation
	for _, b := range data {
		// Convert each byte to its 8-bit binary string representation
		// Use fmt.Sprintf("%08b", b) to ensure leading zeros for 8 bits
		sb.WriteString(fmt.Sprintf("%08b", b))
	}
	return sb.String()
}

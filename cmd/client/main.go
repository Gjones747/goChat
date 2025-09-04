package main

import (
	"log"
	"net/url"

	socketcontroller "github.com/Gjones747/goChat/webSocket"
	//"os"
	//"github.com/Gjones747/goChat/client/views"
	//tea "github.com/charmbracelet/bubbletea"
)

func main() {
	
	url := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/makeRoom"}
	q := url.Query()
	q.Add("room_code", "sixseven")
	q.Add("user_name", "joe")

	url.RawQuery = q.Encode()

	connection, err := socketcontroller.SendUpgrade(url)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(connection)





	//tui := tea.NewProgram(
	//	views.InitialModel(),
	//	tea.WithAltScreen(),
	//	tea.WithMouseCellMotion(),
	//)
	//if _, err := tui.Run(); err != nil {
	//	fmt.Printf("Alas, there's been an error: %v", err)
	//	os.Exit(1)
	//}

}


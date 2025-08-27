package views


// this is the first view that a user is greated with it should show a nice box to input a room code that the user wants to join

type enterRoomView struct {
	input []byte
}


func initEnterRoomView() enterRoomView {
	enterRoomView := enterRoomView {
		input: []byte{},
	}
	return enterRoomView
}

func (model enterRoomView) View() string {
	return "Type in the room you want to visit"
}

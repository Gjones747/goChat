package views


// this is the first view that a user is greated with it should show a nice box to input a room code that the user wants to join

type enterRoomView struct {
	input []byte
}


func makeEnterRoomView() enterRoomView {
	enterRoomView := enterRoomView {
		input: []byte{},
	}
	return enterRoomView
}

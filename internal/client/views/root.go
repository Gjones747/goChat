package views

// this is the root or main view of the cli chat interface
type rootView struct {
	currentPage   int
	enterRoomView enterRoomView
}

func MakeRootView() rootView {
	view := rootView{
		currentPage:   0,
		enterRoomView: makeEnterRoomView(),
	}

	return view
}

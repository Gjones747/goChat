package models

// essentially this needs to be a current list of all the open rooms

type RoomHub struct {
	Rooms map[string]*Room
}

func MakeRoomHub() *RoomHub {
	return &RoomHub{
		Rooms: make(map[string]*Room),
	}
}

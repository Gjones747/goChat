package api

//this will be the logic for updating the listed users in a room

type UserList struct {
	Users [][]string // each list includes a [username, sessionID]
}

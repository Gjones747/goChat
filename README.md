# goChat/Socket

A **lightweight terminal chat client** built in Go.  

- **Minimal & session-based**: no persistent storage, no databases.  
- **Fast & simple**: rooms spin up goroutines when active and shut down when empty.  

![socketdemo mov](https://github.com/user-attachments/assets/5e401632-eceb-44d8-bf03-c30925d64232)  
*(low FPS gif, sorry!)*

---

## Disclaimer

This is **not production-ready**.  
It’s a learning project where I explored:  
- Go’s concurrency model with goroutines & channels  
- WebSocket internals (RFC 6455 implemented with raw `net.Conn`)  
- Go project structure (`cmd/server` and `cmd/client`)  

---

## Getting Started

### There are two ways to install the client:

**Download the socket-client executable from [releases](https://github.com/Gjones747/goChat/releases),**

or

Clone and run the client (server is already deployed):  

```bash
git clone https://github.com/Gjones747/goChat.git
cd goChat
git checkout deploy

# Run
go run cmd/client/main.go

# Build
go build cmd/client/main.go
```

## Server

The server is where most of the goroutine logic lives. The control flow works like this:

1. A client sends a request to the `/joinRoom` endpoint with query params:  
   - `user_name`  
   - `room_code`  
   - `session_id`  
2. A new `User` object is created with the provided `username` and `session_id`.  
3. The server checks its main map of active rooms (`roomHub`).  
   - If the room doesn’t exist, it spins up a new one and adds it to the hub.  
4. The user is added to the room’s map of active users.  
5. Go channels are established on both the user and room side:  
   - When a user sends a message → it’s pushed to the room’s message channel.  
   - The room then broadcasts the message to all user channels.  
   - When a user channel receives a message → it’s sent back to the client via the stored `net.Conn`.  

## SocketController

All WebSocket logic lives here. Instead of relying on a library, I worked directly with Go’s raw `net.Conn` and the RFC 6455 WebSocket protocol doc (https://www.rfc-editor.org/rfc/pdfrfc/rfc6455.txt.pdf).

- Learned about byte-level operations when encoding/decoding WebSocket frames.  
- Great reminder that “there are levels to this” you don’t usually see abstracted away.  
- Skipped ping/pong frames: the server closes `net.Conn` on errors, which was good enough for my use case.  

This was easily the most valuable part of the project from a learning perspective.  

---

## Deploy

I wanted to make the server accessible, so I deployed it using my brother’s old PC:  

- Hosted through a **Cloudflare tunnel**.  
- Workflow:  
  - Develop on `dev` → merge into `main` → rebase `deploy`.  
  - Only “working and tested” code makes it into `deploy`.  
- Learned some DevOps fundamentals while setting this up.  
- Next step: **automated testing**—as the project grew, manual testing became harder to manage.  

---

## Client

The client was my way of avoiding a full web frontend while still learning WebSockets in Go.  

- Built as a **terminal UI** using [Bubble Tea](https://github.com/charmbracelet/bubbletea).  
- Features:  
  - Auto-updates when users join or messages are sent/received.  
  - Highlights your own messages with a unique `session_id`.  
- `session_id` keys are generated using the current timestamp + random bytes (ensuring uniqueness even if usernames collide).  

This approach let me stay in Go-land while still delivering an interactive UI experience.  








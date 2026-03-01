package main

import (
	"errors"
	"final-game-server/internal/engine"
	"final-game-server/internal/games/info"
	"final-game-server/internal/games/paint"
	"final-game-server/internal/shared"
	"log"
	"time"

	"github.com/google/uuid"
)

var validGameTypes = map[string]struct{}{
	"info":  {},
	"paint": {},
}

type Hub struct {
	Rooms map[string]map[string]*engine.Room
	Join  chan *shared.ClientRequest
	// mu     sync.Mutex
	nextId int
}

func NewHub() *Hub {
	hub := &Hub{
		Rooms: make(map[string]map[string]*engine.Room),
		Join:  make(chan *shared.ClientRequest),
		// mu:     sync.Mutex{},
		nextId: 0,
	}

	hub.Rooms["info"] = make(map[string]*engine.Room)
	room, _ := engine.NewRoom("all", "info", 128)
	room.Game = info.NewInfo(room)
	hub.Rooms["info"]["all"] = room

	return hub
}

func (h *Hub) ConnectClientToRoom(client *engine.Client, cr *shared.ClientRequest) error {
	// h.mu.Lock()
	// defer h.mu.Unlock()

	if cr.GType == "info" {
		room := h.Rooms["info"]["all"]
		if room.AddClient(client) != nil {
			log.Println("Unable to add client to info room")
		}
		return
	}

	_, ok := validGameTypes[Type]
	if !ok {
		return errors.New("Invalid RoomType")
	}

	if h.Rooms[cr.GType] == nil {
		h.Rooms[cr.GType] = make(map[string]*engine.Room)
	}

	room, exists := h.Rooms[cr.GType][cr.RID]
	var roomId string
	if (!exists) || (len(room.Clients) >= room.Capacity) {
		foundRoom := false
		for _, existingRoom := range h.Rooms[cr.GType] {
			if existingRoom.AddClient(client) != nil {
				room = existingRoom
				roomId = room.Id
				foundRoom = true
				break
			}
		}
		if !foundRoom {
			if exists {
				roomId = uuid.NewString()
			}
			room = h.CreateRoom(cr.GType, roomId)
			h.Rooms[cr.GType][roomId] = room

			room.AddClient(client)
		}
	}
}

func (h *Hub) CreateRoom(GType string, Id string) *engine.Room {
	if GType == "paint" {
		room := engine.NewRoom(Id, GType, 4)
		room.Game = paint.NewPaintGame()
		return room
	}
	panic("Unable to find Roomconstructor")
}

func (h *Hub) Run() {
	hubTicker := time.NewTicker(1000 * time.Millisecond)

	for {
		select {
		case cr := <-h.Join:

			client := engine.NewClient(h.getNextId(), cr.Conn)

			h.ConnectClientToRoom(client, cr)

			// 		if clientRequest.HInfo == "info" || clientRequest.HInfo == "status" {
			// 			room := h.Info[clientRequest.HInfo]
			// 			client := &Client{Conn: conn, Send: make(chan []byte, 256), Room: room}
			// 			if room.AddPlayer(client) != nil {
			// 				conn.Close()
			// 				log.Printf("Failed to connect client to hub")
			// 			} else {
			// 				go client.WritePump()
			// 				client.ReadPump()
			// 			}
			// 			return
			// 		}

			// 		_, ok := validGameTypes[clientRequest.GType]
			// 		if !ok {
			// 			log.Printf("Clientconnection failed due to invalid game type")
			// 			conn.Close()
			// 			return
			// 		}

			// 		client := &Client{Conn: conn, Send: make(chan []byte, 256), Room: nil}
			// 		// TODO: figure out how to send clientrequest to this function which should use game.AddPlayer to dynamically init the correct Player from clientrequest data
			// 		h.ConnectClientToRoom(clientRequest.GType, clientRequest.RID, client)
			// 		go client.WritePump()
			// 		client.ReadPump()

			// 		// room := hub.GetOrCreateRoom(gType, rID)
			// 		// conn, err := upgrader.Upgrade(w, r, nil)
			// 		// if err != nil {
			// 		// 	return
			// 		// }

			// 		// use clientconsturction NewClient instead!!!!!!!
			// 		// client := &Client{Conn: conn, Send: make(chan []byte, 256), Room: room}
			// 		// room.Join <- client

			// 		// go client.WritePump()
			// 		// client.ReadPump()

		case <-hubTicker.C:
			GSList := make([]*shared.ShareableGameState, 0)
			for Type := range h.Rooms {
				for _, room := range h.Rooms[Type] {
					GSList = append(GSList, room.Game.GetShareableGameState())
				}
			}
			if infoGame, ok := h.Rooms["info"]["all"].Game.(*info.Info); ok {
				infoGame.SetGSList(GSList)
			}
		}

	}
}

func (h *Hub) getNextId() int {
	h.nextId++
	return h.nextId
}

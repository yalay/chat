package controllers

import (
	"sync"
	"time"
)

const (
	kMaxRoomNum        = 16
	kBroadcastDuration = time.Second
	kMaxMsgBuf         = 16
)

var allRooms map[int]*Room
var roomIdsFlag []bool

type Room struct {
	Id      int
	clients map[*Client]bool // all clients in this room
	msgBuf  []string
	sync.RWMutex
}

func init() {
	allRooms = make(map[int]*Room, kMaxRoomNum)
	roomIdsFlag = make([]bool, kMaxRoomNum)
}

func NewRoom() *Room {
	roomId := getAvailableRoomId()
	if roomId == 0 {
		return nil
	}

	newRoom := &Room{
		Id:      roomId,
		clients: make(map[*Client]bool),
		msgBuf:  make([]string, 0, kMaxMsgBuf),
	}
	allRooms[roomId] = newRoom
	go newRoom.run()
	return newRoom
}

func GetRoom(id int) *Room {
	return allRooms[id]
}

func DelRoom(id int) {
	room, ok := allRooms[id]
	if ok {
		room.Lock()
		for client, _ := range room.clients {
			room.DelClient(client)
		}
		room.Unlock()
		delete(allRooms, id)
	}
}

func (r *Room) run() {
	ticker := time.NewTicker(kBroadcastDuration)
	for {
		select {
		case <-ticker.C:
			r.broadcast()
		default:
			if len(r.msgBuf) >= kMaxMsgBuf {
				r.broadcast()
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func (r *Room) AddClient(client *Client) {
	r.Lock()
	r.clients[client] = true
	r.Unlock()
}

func (r *Room) DelClient(client *Client) {
	r.Lock()
	if _, ok := r.clients[client]; ok {
		delete(r.clients, client)
	}
	r.Unlock()
}

func (r *Room) Broadcast(msg string) {
	r.msgBuf = append(r.msgBuf, msg)
}

func (r *Room) broadcast() {
	r.RLock()
	curMsgs := make([]string, len(r.msgBuf))
	copy(curMsgs, r.msgBuf)
	r.msgBuf = nil
	for client, _ := range r.clients {
		go func(curClient *Client) {
			curClient.SendMessage(curMsgs)
		}(client)
	}
	r.RUnlock()
}

// room id start from 1. 0 is invalid.
func getAvailableRoomId() int {
	for i, used := range roomIdsFlag {
		if !used {
			roomIdsFlag[i] = true
			return i + 1
		}
	}
	return 0
}

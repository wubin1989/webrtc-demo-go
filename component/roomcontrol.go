package RoomControl

import (
	"github.com/googollee/go-socket.io"
	"log"
	"sync"
)

type RoomControl struct {
	clientCount map[string]int
	adaptor     socketio.BroadcastAdaptor
	mu          sync.Mutex
}

type clientsMediaDescription struct {
	clients struct{}
}

func NewRoomControl(adaptor socketio.BroadcastAdaptor) *RoomControl {
	return &RoomControl{
		clientCount: make(map[string]int),
		adaptor:     adaptor,
		mu:          sync.Mutex{},
	}
}

func (rc *RoomControl) Join(room string, socket socketio.Socket) error {
	err := rc.adaptor.Join(room, socket)
	if err == nil {
		rc.mu.Lock()
		rc.clientCount[room]++
		rc.mu.Unlock()
	}
	return err
}

func (rc *RoomControl) Leave(room string, socket socketio.Socket) error {
	err := rc.adaptor.Leave(room, socket)
	if err == nil {
		rc.mu.Lock()
		rc.clientCount[room]++
		rc.mu.Unlock()
	}
	return err
}

func (rc *RoomControl) GetClientsNumber(room string) int {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return rc.clientCount[room]
}

func (rc *RoomControl) Send(ignore socketio.Socket, room, message string, args ...interface{}) error {
	return rc.adaptor.Send(ignore, room, message, args...)
}

func (rc *RoomControl) DescribeRoom(name string) clientsMediaDescription {
	sResult := &clientsMediaDescription{}

	rc.adaptor.GetAllClients(name)
	if rc.adaptor[name] != nil {
		clientsMap := rc.adaptor[name]
	}
	for clientId, clientSocket := range clientsMap {
		log.Println("Id is %s", clientId)
		sResult.clients[clientId] = clientSocket.resources
	}
	return *sResult
}

package RoomControl

import (
	"github.com/wubin1989/go-socket.io"
	"log"
	// "reflect"
	"sync"
)

type RoomControl struct {
	clientCount map[string]int
	adaptor     socketio.BroadcastAdaptor
	mu          sync.Mutex
}

type clientsMediaDescription struct {
	Clients map[string]interface{}
}

func NewRoomControl(adaptor socketio.BroadcastAdaptor) *RoomControl {
	return &RoomControl{
		clientCount: make(map[string]int),
		adaptor:     adaptor,
		mu:          sync.Mutex{},
	}
}

func (rc *RoomControl) Join(room string, socket socketio.Socket) error {
	// log.Println("=====================================================")
	// log.Println(room)
	// log.Println(socket)
	err := rc.adaptor.Join(room, socket)
	if err == nil {
		rc.mu.Lock()
		rc.clientCount[room]++
		rc.mu.Unlock()
	}
	if err != nil {
		log.Fatal(err)
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

func (rc *RoomControl) GetAllClients(room string) (map[string]socketio.Socket, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	socketMap, err := rc.adaptor.GetAllClients(room)
	if err != nil {
		log.Fatal(err)
	}
	return socketMap, nil
}

func (rc *RoomControl) DescribeRoom(room string) clientsMediaDescription {
	sResult := &clientsMediaDescription{
		Clients: make(map[string]interface{}),
	}
	socketMap, err := rc.GetAllClients(room)

	if err != nil {
		log.Fatal(err)
	}

	for clientId, clientSocket := range socketMap {
		//log.Printf("Id is %s", clientId)
		sResult.Clients[clientId] = *clientSocket.GetResources()
	}
	return *sResult
}

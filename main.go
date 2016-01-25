package main

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-socket.io"
	//"github.com/satori/go.uuid"
	"log"
	"net/http"
	"strconv"
	"time"
	"webrtc-demo-go/component"
)

type ExtendedServer struct {
	server *socketio.Server
	//result
}

type handshakeDetail struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Sid      string `json:"sid"`
	RoomType string `json:"roomType"`
	Payload  struct {
		Type      string `json:"type"`
		Sdp       string `json:"sdp"`
		Candidate struct {
			Candidate     string `json:"candidate"`
			SdpMid        string `json:"sdpMid"`
			SdpMLineIndex int    `json:"sdpMLineIndex"`
		} `json:"candidate"`
	} `json:"payload"`
	Prefix string `json:"prefix"`
	Type   string `json:"type"`
}

//{"to":"FBZY_JLuq6vSZH6lk2bX","sid":"1453458111055","roomType":"video","type":"offer","prefix":"webkit"}

// type Payload struct {
// 	Candidate
// }

// type Candidate struct {
// 	Type string `json: "type"`
// 	Sdp  string `json: "sdp"`
// 	// candidate     string
// 	// sdpMid        string
// 	// sdpMLineIndex int
// }

/**
 *
{
    "payload":{
        "candidate":{
            "candidate":"candidate:1984149417 1 udp 2122260223 192.168.14.8 53365 typ host generation 0",
            "sdpMid":"audio",
            "sdpMLineIndex":0
        }
    }
}
*/

func main() {
	eServer := &ExtendedServer{}
	server, err := socketio.NewServer(nil)
	eServer.server = server
	if err != nil {
		log.Fatal(err)
	}

	defaultAdaptor := make(socketio.Broadcast)

	roomCtrl := RoomControl.NewRoomControl(defaultAdaptor)
	server.SetAdaptor(roomCtrl)

	// log.Printf("roomCtrl is %+v", roomCtrl)

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so = so.SetResources(false, true, false)

		// log.Printf("%+v", so)

		so.On("message", func(msg string) {

			convertedMsg, err := strconv.Unquote(msg)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("msg is coming ======================================================", convertedMsg)
			var m []byte

			hsDetail := &handshakeDetail{}

			err = json.Unmarshal([]byte(convertedMsg), hsDetail)

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("%+v", hsDetail)

			clientsInChat, _ := roomCtrl.GetAllClients("chat")
			toSocket := clientsInChat[hsDetail.To]

			// {
			// 	to: '/#PVVNhym5878Ldu0mAAAC',
			// 	sid: '1452784289306',
			// 	roomType: 'video',
			// 	type: 'candidate',
			// 	payload:
			// 	{ candidate:
			// 	  { candidate: 'candidate:3545837919 2 udp 2122129150 192.168.42.44 53907 typ host generation 0',
			// 	    sdpMid: 'video',
			// 	    sdpMLineIndex: 1 } },
			// 	prefix: 'webkit'
			// }

			hsDetail.From = so.Id()
			m, err = json.Marshal(hsDetail)

			if err != nil {
				log.Fatal(err)
			}

			// log.Printf("%+v", hsDetail)
			// log.Printf("%s", m)

			//so.Emit("message", m)
			toSocket.Emit("message", m)
		})

		so.On("disconnection", func() {
			log.Println("on disconnect", "   ", time.Now().UnixNano()/1e6)
		})

		so.On("create", func(name string) {
			// log.Printf("socket name is %s", name)

			// if (arguments.length == 2) {
			//     cb = (typeof cb == 'function') ? cb : function() {};
			//     name = name || uuid();
			// } else {
			//     cb = name;
			//     name = uuid();
			// }

			so.Emit("room factory message back", name)

			roomCtrl.Join(name, so)

			// check if exists
			// var room = io.nsps['/'].adapter.rooms[name];
			// if (room && room.length) {
			//     safeCb(cb)('taken');
			// } else {
			//     join(name);
			//     safeCb(cb)(null, name);
			// }
		})

		so.On("join", func(name string) {
			joinFeedBack := roomCtrl.DescribeRoom(name)
			//log.Printf("%+v", joinFeedBack)
			feedBackMsg, err := json.Marshal(joinFeedBack)
			if err != nil {
				log.Fatal(err)
			}
			// log.Printf("%s", feedBackMsg)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("===============================================")
			log.Println(fmt.Sprintf("%s", feedBackMsg))
			so.Emit("join feed back", fmt.Sprintf("%s", feedBackMsg))
			roomCtrl.Join(name, so)
		})

		//"stunservers": [{ "url": "stun:stun.l.google.com:19302" }]

		// stunServers := `[{ "url": "stun:stun.l.google.com:19302" }]`
		// so.Emit("stunservers", stunServers)

		// create shared secret nonces for TURN authentication
		// the process is described in draft-uberti-behave-turn-rest
		// var credentials = [];
		// // allow selectively vending turn credentials based on origin.
		// var origin = client.handshake.headers.origin;
		// if (!config.turnorigins || config.turnorigins.indexOf(origin) !== -1) {
		//     config.turnservers.forEach(function (server) {
		//         var hmac = crypto.createHmac('sha1', server.secret);
		//         // default to 86400 seconds timeout unless specified
		//         var username = Math.floor(new Date().getTime() / 1000) + (server.expiry || 86400) + "";
		//         hmac.update(username);
		//         credentials.push({
		//             username: username,
		//             credential: hmac.digest('base64'),
		//             urls: server.urls || server.url
		//         });
		//     });
		// }
		// client.emit('turnservers', credentials);

	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8889...")
	log.Fatal(http.ListenAndServe(":8889", nil))
}

// function clientsInRoom(name) {
//     return io.sockets.clients(name).length;
// }

// function removeFeed(type) {
//     if (client.room) {
//         io.sockets.in(client.room).emit('remove', {
//             id: client.id,
//             type: type
//         });
//         if (!type) {
//             client.leave(client.room);
//             client.room = undefined;
//         }
//     }
// }

//         // we don't want to pass "leave" directly because the
//         // event type string of "socket end" gets passed too.
//         client.on('disconnect', function() {
//             removeFeed();
//         });
//         client.on('leave', function() {
//             removeFeed();
//         });

// // function safeCb(cb) {
// //     if (typeof cb === 'function') {
// //         return cb;
// //     } else {
// //         return function() {};
// //     }
// // }

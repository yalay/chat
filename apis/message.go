package apis

import (
	"net/http"
	"controllers"
	"util"

	"github.com/gorilla/websocket"
	"path"
	"strconv"
)

// for test
var _ = controllers.NewRoom()

var wsUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ReqMessage struct {
	roomId int
	msg string
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	reqValues := r.URL.Query()
	referUrl := reqValues.Get("refer")
	if referUrl == "" {
		util.Render.Text(w, http.StatusBadRequest, "empty refer")
		return
	}
	roomId, err := strconv.Atoi(path.Base(referUrl))
	if err != nil {
		util.Render.Text(w, http.StatusBadRequest, "invalid room id")
		return
	}

	room := controllers.GetRoom(roomId)
	if room == nil {
		util.Render.Text(w, http.StatusNotFound, "no room")
		return
	}

	conn, err := wsUpgrade.Upgrade(w, r, nil)
	if err != nil {
		util.Log.Error(err.Error())
		return
	}

	client := controllers.NewClient(conn)
	client.Join(room)
	go client.ReadMessage()
}

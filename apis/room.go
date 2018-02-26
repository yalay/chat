package apis

import (
	"controllers"
	"net/http"
	"strconv"
	"util"

	"github.com/go-chi/chi"
)

type RspAddRoom struct {
	RoomId int `json:"room_id"`
}

func RoomHandler(w http.ResponseWriter, r *http.Request) {
	util.Render.HTML(w, http.StatusOK, "room", nil)
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	room := controllers.NewRoom()
	if room == nil {
		util.Render.Text(w, http.StatusConflict, "room not available")
		return
	}

	util.Render.JSON(w, http.StatusOK, RspAddRoom{
		RoomId: room.Id,
	})
}

func DelRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomId, _ := strconv.Atoi(chi.URLParam(r, "roomId"))
	if roomId == 0 {
		util.Render.Text(w, http.StatusBadRequest, "no room id.")
		return
	}

	room := controllers.GetRoom(roomId)
	if room == nil {
		util.Render.Text(w, http.StatusBadRequest, "room not exist")
		return
	}

	controllers.DelRoom(roomId)
	util.Render.JSON(w, http.StatusOK, "OK")
}

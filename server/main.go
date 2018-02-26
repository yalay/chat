package main

import (
	"flag"
	"net/http"
	"util"
	"middlewares"
	"strconv"
	"apis"

	"github.com/go-chi/chi"
	_ "github.com/go-chi/chi/middleware"
)

var logFile string
var tplPath string
var listenPort int

func init() {
	flag.StringVar(&logFile, "log", "./chat.log", "log file path")
	flag.StringVar(&tplPath, "tpl", "./views", "template file root path")
	flag.IntVar(&listenPort, "p", 8000, "server listen port")
	flag.Parse()

	util.InitLog(logFile)
	util.InitRender(tplPath)
}

func main() {
	router := chi.NewRouter()
	router.Use(
		middlewares.AccessMiddleware,
		middlewares.RecoverMiddleware,
	)

	router.Get("/", apis.HomeHandler)
	router.Get(`/rooms/{roomId:\d+}`, apis.HomeHandler)
	router.Post(`/rooms`, apis.CreateRoomHandler)
	router.Delete(`/rooms/{roomId:\d+}`, apis.DelRoomHandler)
	router.Get(`/messages`, apis.MessagesHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(listenPort), router)
	if err != nil {
		util.Log.Fatal(err)
	}
}

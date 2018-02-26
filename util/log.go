package util

import (
	"os"
	"path/filepath"

	"github.com/op/go-logging"
)

var Log *logging.Logger

func InitLog(logPath string) {
	openFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("open log file error:" + err.Error())
	}

	logBackend := logging.NewBackendFormatter(
		logging.NewLogBackend(openFile, "", 0),
		logging.MustStringFormatter(
			`[%{time:2006-01-02 15:04:05}] [%{level:.4s}] [%{shortfile}] - %{message}`,
		))

	Log = logging.MustGetLogger(filepath.Base(logPath))
	Log.SetBackend(logging.AddModuleLevel(logBackend))
}

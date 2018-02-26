package middlewares

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"net/http"
	"time"
	"util"
)

type LogResponseWriter struct {
	http.ResponseWriter
	rspBody *bytes.Buffer
	Status  int
}

func (r *LogResponseWriter) Write(p []byte) (int, error) {
	r.rspBody.Write(p)
	return r.ResponseWriter.Write(p)
}

func (r *LogResponseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *LogResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := r.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("chi/middleware: http.Hijacker is unavailable on the writer")
}

func AccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logRspWriter := &LogResponseWriter{w, bytes.NewBufferString(""), http.StatusOK}

		next.ServeHTTP(logRspWriter, r)

		clientIp := util.RealIP(r)
		elapsed := float64(time.Now().Sub(startTime).Nanoseconds()) / 1e6
		util.Log.Debugf("%s %s %d %.3fms %s %s", r.Method, r.URL.String(),
			logRspWriter.Status, elapsed, clientIp, logRspWriter.rspBody)
	})
}

package middlewares

import (
	"net/http"
	"net"
	"bufio"
)

type Hijacker struct {
	http.ResponseWriter
}

func (writer *Hijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func HijackMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&Hijacker{w}, r)
	})
}

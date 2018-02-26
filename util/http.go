package util

import (
	"net/http"
	"strings"
)

var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

func RealIP(r *http.Request) string {
	var realIp string
	if xffIp := r.Header.Get(xForwardedFor); xffIp != "" {
		i := strings.Index(xffIp, ", ")
		if i < 0 {
			i = len(xffIp)
		}
		realIp = xffIp[:i]
	} else if xrIp := r.Header.Get(xRealIP); xrIp != "" {
		realIp = xrIp
	}

	if realIp != "" {
		return realIp
	}
	return r.RemoteAddr
}

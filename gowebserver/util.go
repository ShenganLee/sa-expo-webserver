package gowebserver

import (
	"net"
	"net/http"
	"net/http/httputil"
	"regexp"
	"time"

	"github.com/vulcand/oxy/v2/forward"
)

// func getTime() string {
// 	// 格式化时间
// 	return time.Now().Format("2006-01-02 15:04:05")
// }

func isBlabk(str string) bool {
	if len(str) == 0 {
		return true
	}

	pattern := "^\\s*$"
	regexp := regexp.MustCompile(pattern)
	return regexp.MatchString(str)
}

func findAvailablePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return -1
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return -1
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func newForward() *httputil.ReverseProxy {
	fwd := forward.New(false)
	fwd.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		MaxIdleConnsPerHost: 100,
	}

	return fwd
}

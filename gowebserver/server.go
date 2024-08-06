package gowebserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/vulcand/oxy/v2/forward"
)

// 获取未被使用的端口
func GetPort(startPort int) int {
	if startPort < 80 {
		startPort = 80
	}

	for port := startPort; port < 65535; port++ {
		addr := fmt.Sprintf("localhost:%d", port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			// 端口已被使用
			continue
		}
		// 立即关闭监听，不进行实际的监听
		listener.Close()
		return port
	}

	return 0
}

func isBlabk(str string) bool {
	if len(str) == 0 {
		return true
	}

	pattern := "^\\s*$"
	regexp := regexp.MustCompile(pattern)
	return regexp.MatchString(str)
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}

type Proxy struct {
	Path    string   `json:"Path"`
	Target  string   `json:"Target"`
	Include []string `json:"Include"`
	Exclude []string `json:"Exclude"`
}

func noCacheHandle(w http.ResponseWriter) {
	// 设置HTTP头部，告知不缓存响应
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
}

func mainHandlerFunc(w http.ResponseWriter, r *http.Request, dir http.Dir, h http.Handler) {
	noCacheHandle(w)

	if r.URL.Path == "/" {
		h.ServeHTTP(w, r)
		return
	}

	isExist := fileExists(string(dir) + r.URL.Path)

	// 静态文件处理
	if isExist {
		h.ServeHTTP(w, r)
	} else {
		// 反向代理到 '/' 目录 兼容前端路由
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = "/"
		r2.URL.RawPath = "/"
		h.ServeHTTP(w, r2)
	}
}

var webServerClosePath = "/go-web-server-close"
var webServerRestartPath = "/go-web-server-restart"

func generateProxyServer(port int, fileDir string, proxys []Proxy) {
	mux := http.NewServeMux()

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	useFileServer := !isBlabk(fileDir)
	log.Printf("\nuseFileServer: %v\n", useFileServer)
	var (
		httpDir    http.Dir
		fileServer http.Handler
	)

	closed := false
	fmt.Printf("\nserver closed: %v\n", closed)

	mux.Handle(webServerClosePath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		closed = true
		server.Close()
	}))

	mux.Handle(webServerRestartPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		closed = false
		server.Close()
	}))

	if len(proxys) > 0 {
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

		for _, p := range proxys {
			target, err := url.Parse(p.Target)
			if err != nil {
				log.Printf("url.Parse Error: %s\n", p.Target)
				break
			}

			// http.Handle(p.Path, proxy)
			// http.Handle(p.Path, http.StripPrefix(p.Path, proxy))

			mux.Handle(p.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				urlPath := r.URL.Path

				if p.Include != nil && len(p.Include) > 0 {
					for _, include := range p.Include {

						regex, err := regexp.Compile(include)
						if err != nil {
							continue
						}

						if regex.MatchString(urlPath) {
							fmt.Printf("proxy: %s, target: %s\n", urlPath, p.Target)
							r.URL = target
							fwd.ServeHTTP(w, r)
							return
						}
					}
				}

				if useFileServer && p.Exclude != nil && len(p.Exclude) > 0 {
					for _, exclude := range p.Exclude {

						regex, err := regexp.Compile(exclude)
						if err != nil {
							// fmt.Printf("\nExclude Compile Error: %v\n", err)
							continue
						}

						// fmt.Printf("\n%s MatchString %s: %v\n", exclude, urlPath, regex.MatchString(urlPath))

						if regex.MatchString(urlPath) {
							log.Printf("exclude path: %s\n", urlPath)
							mainHandlerFunc(w, r, httpDir, fileServer)
							return
						}
					}
				}

				log.Printf("proxy: %s, target: %s\n", urlPath, p.Target)
				r.URL = target
				fwd.ServeHTTP(w, r)
			}))

		}
	}

	if useFileServer {
		httpDir = http.Dir(fileDir)
		fileServer = http.FileServer(httpDir)

		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fmt.Printf("path: %s\n", r.URL.Path)
			mainHandlerFunc(w, r, httpDir, fileServer)
		}))

		// http.Handle("/", http.FileServer(http.Dir(fileDir)))
	}

	// server.ListenAndServe()
	// if !closed {
	// 	generateProxyServer(port, fileDir, proxys)
	// }

	go func() {
		server.ListenAndServe()
		if !closed {
			generateProxyServer(port, fileDir, proxys)
		}
	}()
}

func Start(fileDir string, proxyStr string) string {
	var proxys []Proxy
	if !isBlabk(proxyStr) {
		json.Unmarshal([]byte(proxyStr), &proxys)
	}

	port := GetPort(9527)
	if port == 0 {
		return ""
	}

	// go generateProxyServer(port, fileDir, proxys)
	// time.Sleep(200 * time.Millisecond) // 延迟200 毫秒
	generateProxyServer(port, fileDir, proxys)

	return fmt.Sprintf("http://127.0.0.1:%d", port)
}

func Stop(addr string) {
	http.Get(addr + webServerClosePath)
}

func Restart(addr string) {
	http.Get(addr + webServerRestartPath)
}

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hello is a simple hello, world demonstration web server.
//
// It serves version information on /version and answers
// any other request like /name by saying "Hello, name!".
//
// See golang.org/x/example/outyet for a more sophisticated server.
package gowebserver

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

// func TestServer1(t *testing.T) {
// 	regex := regexp.MustCompile(`^/studio_api/(((publish_site|uploads)/.*)|public_components_manifest.json)$`)
// 	fmt.Printf("%v", regex.MatchString("/studio_api/check_iot_token"))
// }

func TestServer(t *testing.T) {
	// log.Println(getTime())
	// Run the server.
	serverConfig := &ServerConfig{
		Proxys: Proxys{
			{Path: "/cdn/", Target: "http://www.app.rc-ess.com"},
			{Path: "/minio/", Target: "http://www.app.rc-ess.com"},
			{Path: "/websocket/", Target: "http://www.app.rc-ess.com"},
			{Path: "/jowoiot-adapter/", Target: "http://www.app.rc-ess.com"},
			{Path: "/jowoiot-proxy/", Target: "http://www.app.rc-ess.com"},
			{
				Path:    "/studio_api/",
				Target:  "http://www.app.rc-ess.com",
				Skipper: []string{`^/studio_api/(((publish_site|uploads)/.*)|public_components_manifest.json)$`},
			},
		},
		Routers: Routers{
			{
				Path:     "/",
				FilePath: "./web",
				IndexHeaders: IndexHeaders{
					ResponseHeaders: map[string]string{
						"Cache-Control": "no-store",
						"Pragma":        "no-cache",
					},
				},
			},
		},
	}

	// SetLogFile("./server.log")

	// serverConfig := &ServerConfig{
	// 	Proxys: Proxys{
	// 		{Path: "/", Target: "http://lab.jowoiot.com:18308"},
	// 	},
	// }

	dataType, err := json.Marshal(serverConfig)
	if err != nil {
		return
	}

	dataString := string(dataType)
	addr := Start(dataString)

	// fmt.Printf("dataString: %s\n", dataString)
	// addr := Start(dataString)

	// addr := Start("{\"Proxys\":[{\"Path\":\"/\",\"Target\":\"http://lab.jowoiot.com:18308\"}]}")

	fmt.Printf("启动服务: %s\n", addr)

	go func() {
		time.Sleep(3 * time.Second)
		// Restart()
		logger.Println("Healthy", Healthy())
		logger.Println("IsRunning", IsRunning())
		logger.Println("ServerUrl", ServerUrl())
		// time.Sleep(10 * time.Second)
		// Stop()
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case sig := <-sigc:
			// log.Printf("Received signal: %s", sig.String())
			switch sig {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
				// log.Println("close server...")
				return
			case syscall.SIGHUP:
				// 这里可以添加配置文件重新加载逻辑
				// log.Println("Reloading server configuration...")
				// return
			}
		default:
			// 这里可以添加其他守护进程逻辑，例如健康检查、日志记录等
			time.Sleep(1 * time.Second)
		}
	}
}

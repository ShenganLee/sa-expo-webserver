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
	"testing"
)

func TestServer(t *testing.T) {
	// Run the server.
	// proxys := []Proxy{
	// 	// {"path": "/", "target": "https://studio.jowocloud.com"},
	// 	// {"/minio/", "http://192.168.3.167:18400"},
	// 	{Path: "/minio/", Target: "http://lab.jowoiot.com:18400"},
	// 	{Path: "/websocket/", Target: "http://lab.jowoiot.com:18400"},
	// 	{Path: "/jowoiot-adapter/", Target: "http://lab.jowoiot.com:18400"},
	// 	{Path: "/jowoiot-proxy/", Target: "http://lab.jowoiot.com:18400"},
	// 	{
	// 		Path:    "/studio_api/",
	// 		Target:  "http://lab.jowoiot.com:18400",
	// 		Exclude: []string{`^/studio_api/(((publish_site|uploads)/.*)|public_components_manifest.json)$`},
	// 	},
	// }

	// proxyConfig := []ProxyConfig{{
	// 	Targets: []string{"http://lab.jowoiot.com:18400"},
	// 	Rewrite: map[string]string{
	// 		"^/minio/*":           "/minio/$1",
	// 		"^/websocket/*":       "/websocket/$1",
	// 		"^/jowoiot-adapter/*": "/jowoiot-adapter/$1",
	// 		"^/jowoiot-proxy/*":   "/jowoiot-proxy/$1",
	// 	},
	// 	Include: []string{`^/(minio|websocket|jowoiot-adapter|jowoiot-proxy)/.*$`},
	// 	Exclude: []string{`^/studio_api/(((publish_site|uploads)/.*)|public_components_manifest.json)$`},
	// }}

	proxys := []Proxy{
		{Path: "/", Target: "http://47.237.20.177"},
	}

	dataType, err := json.Marshal(proxys)
	// dataType, err := json.Marshal(proxyConfig)
	if err != nil {
		return
	}

	dataString := string(dataType)
	// fmt.Printf("dataString: %s\n", dataString)
	// addr := Start("./web", dataString)
	addr := Start("", dataString)

	// addr := Start("", dataString)

	fmt.Printf("启动服务: %s\n", addr)

	// go func() {
	// var res string
	// fmt.Print("是否继续：y/n")
	// fmt.Scan(&res)

	// if res == "y" {
	// 	Stop(addr)
	// }
	// }()
}

// func TestServer1(t *testing.T) {
// 	// regex, err := regexp.Compile(`/publish_site/`)
// 	regex, err := regexp.Compile("^/studio_api/(((publish_site|uploads)/.*)|public_components_manifest.json)$")

// 	if err != nil {
// 		panic(err)
// 	}

// 	s1 := "/studio_api/publish_site/publish_site_info.json"
// 	f1 := regex.MatchString(s1)
// 	fmt.Printf("s1: %v\n", f1)

// 	s2 := "/studio_api/uploads/public_component_package/1/1715840926/_-_-index-js.59c431e6.chunk.js"
// 	f2 := regex.MatchString(s2)
// 	fmt.Printf("s2: %v\n", f2)

// 	s3 := "/studio_api/public_components_manifest.json"
// 	f3 := regex.MatchString(s3)
// 	fmt.Printf("s3: %v\n", f3)

// 	s4 := "/studio_api/publish_site1/aaa.txt"
// 	f4 := regex.MatchString(s4)
// 	fmt.Printf("s4: %v\n", f4)

// }

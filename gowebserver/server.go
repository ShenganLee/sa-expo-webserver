package gowebserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

var logOut = os.Stdout
var logger = log.New(os.Stdout, "go-web-server: ", log.LstdFlags)

type Config struct {
	IsRunning    bool
	ServerUrl    string
	ServerConfig *ServerConfig

	logFile string

	echo *echo.Echo
}

func (c *Config) Reset() {
	config.IsRunning = false
	config.ServerUrl = ""
	config.echo = nil
}

func newConfig() *Config {
	return &Config{
		IsRunning:    false,
		ServerUrl:    "",
		ServerConfig: nil,

		logFile: "",

		echo: nil,
	}
}

var config = newConfig()
var healthPath = "/go-web-server-health"

func Start(serverConfigStr string) string {
	logger.Println("Function: Start")

	if config.IsRunning {
		for {
			if config.IsRunning && !isBlabk(config.ServerUrl) {
				return config.ServerUrl
			}

			if !config.IsRunning {
				break
			}
		}
	}

	if isBlabk(serverConfigStr) {
		return ""
	}

	var serverConfig ServerConfig
	json.Unmarshal([]byte(serverConfigStr), &serverConfig)

	if serverConfig.Port == 0 {
		port := findAvailablePort()
		if port < 0 {
			return ""
		}
		serverConfig.Port = port
	}

	config.ServerConfig = &serverConfig

	generateServer()

	return config.ServerUrl
}

func Stop() {
	logger.Println("Function: Stop")

	echo := config.echo

	config = newConfig()

	if echo != nil {
		echo.Close()
	}
}

func Restart() {
	logger.Println("Function: Restart")

	if config.IsRunning {
		for {
			if config.IsRunning && !isBlabk(config.ServerUrl) {
				break
			}
		}

		if config.echo != nil {
			config.echo.Close()
		}

		config.Reset()

		generateServer()
	}
}

func IsRunning() bool {
	logger.Println("Function: IsRunning")

	return config.IsRunning
}

func ServerUrl() string {
	logger.Println("Function: ServerUrl")

	return config.ServerUrl
}

func Healthy() bool {
	logger.Println("Function: Healthy")

	if !config.IsRunning {
		return false
	}

	client := &http.Client{}
	req, eeqErr := http.NewRequest("GET", config.ServerUrl+healthPath, nil)
	if eeqErr != nil {
		return false
	}

	resp, respErr := client.Do(req)
	if respErr != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func SetLogFile(logFile string) {
	if isBlabk(logFile) {
		return
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	// defer file.Close()

	config.logFile = logFile
	logOut = file

	// 设置日志的标准输出为文件
	logger.SetOutput(file)
	logger.Println("Log entry written to file " + logFile)
	logger.Println("Function: SetLogFile")
}

func LogFileClose() {
	if logOut != os.Stdout {
		logOut.Close()
		logOut = os.Stdout
	}
	config.logFile = ""
}

func generateServer() {
	config.IsRunning = true
	config.echo = echo.New()

	// health
	config.echo.GET(healthPath, func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Proxys
	config.ServerConfig.Proxys.Add(config.echo)

	// Routes
	config.ServerConfig.Routers.Add(config.echo)

	ch := make(chan bool, 2)

	go func() {
		err := config.echo.Start(fmt.Sprintf(":%d", config.ServerConfig.Port))

		logger.Println("server start err", err)

		if ch != nil && err != nil {
			config.Reset()
			ch <- true
			close(ch)
		}

	}()

	for {
		time.Sleep(100 * time.Millisecond)
		if config.IsRunning {
			config.ServerUrl = fmt.Sprintf("http://127.0.0.1:%d", config.ServerConfig.Port)
			if ch != nil {
				ch <- true
				close(ch)
			}
			break
		}
	}

	for v := range ch {
		if v {
			ch = nil
			break
		}
	}

	logger.Println("server start end", config.IsRunning, config.ServerUrl)
}

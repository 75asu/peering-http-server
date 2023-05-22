package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"io"
)

var (
	configFilePath = "/config/hostports"
	pingInterval   = 30 * time.Second
)

func main() {
	// Create a channel to receive configuration update signals
	configUpdateCh := make(chan struct{})

	// Start a goroutine to periodically check for configuration updates
	go watchConfigUpdates(configUpdateCh)

	// Start the HTTP server to handle "/ping" requests
	go startServer()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Initial ping requests
	sendPingRequests()

	// Schedule periodic ping requests
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-configUpdateCh:
			// Configuration updated, send new ping requests
			sendPingRequests()
		case <-ticker.C:
			// Time to send periodic ping requests
			sendPingRequests()
		}
	}
}

func startServer() {
	// Expose path "/ping"
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		// Get the hostname of the server that handled the request
		hostname, _ := os.Hostname()

		// Include the response from the specific server in the log message
		response := fmt.Sprintf("pong from %s", hostname)
		fmt.Fprint(w, response)
		log.Print(response)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}


func sendPingRequests() {
	hostPorts, err := readConfig()
	if err != nil {
		log.Printf("Failed to read configuration: %v", err)
		return
	}

	hostPortList := strings.Split(hostPorts, "\n")

	// Send GET request to each host and port
	for _, hostPortLine := range hostPortList {
		hostPortLine = strings.TrimSpace(hostPortLine)

		if hostPortLine != "" {
			hostPort := strings.SplitN(hostPortLine, ",", 2)
			if len(hostPort) != 2 {
				log.Printf("Invalid host-port configuration: %s", hostPortLine)
				continue
			}

			host := hostPort[0]
			port := hostPort[1]

			// Send GET request to the host and port
			url := fmt.Sprintf("http://%s:%s/ping", host, port)

			// Check if the server is available
			response, err := http.Get(url)
			if err != nil {
				log.Printf("Server not available at %s", url)
				continue
			}

			// Read and discard the response body
			_, _ = io.Copy(ioutil.Discard, response.Body)
			response.Body.Close()

			// Include the server name in the log message
			log.Printf("pong from %s:%s", host, port)
		}
	}
}


func readConfig() (string, error) {
	hostPortData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read host-port configuration: %v", err)
	}
	hostPorts := string(hostPortData)

	return hostPorts, nil
}

func watchConfigUpdates(configUpdateCh chan<- struct{}) {
	// Start a goroutine to watch for changes in the config file
	go func() {
		var lastModified time.Time

		for {
			// Check the last modified time of the config file
			fileInfo, err := os.Stat(configFilePath)
			if err != nil {
				log.Printf("Failed to get file info: %v", err)
				continue
			}

			if fileInfo.ModTime().After(lastModified) {
				// The config file has been modified, send an update signal
				configUpdateCh <- struct{}{}
				lastModified = fileInfo.ModTime()
			}

			// Sleep for a short interval before checking again
			time.Sleep(1 * time.Second)
		}
	}()
}







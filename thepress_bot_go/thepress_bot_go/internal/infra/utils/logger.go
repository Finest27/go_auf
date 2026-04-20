package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type LogHub struct {
	mu      sync.Mutex
	clients map[chan string]bool
}

var Hub = &LogHub{
	clients: make(map[chan string]bool),
}

func (h *LogHub) Register() chan string {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch := make(chan string, 100)
	h.clients[ch] = true
	return ch
}

func (h *LogHub) Unregister(ch chan string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, ch)
	close(ch)
}

type wsEvent struct {
	Type      string      `json:"type"`
	Message   string      `json:"message,omitempty"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

func BroadcastEvent(eventType string, data interface{}) {
	event := wsEvent{
		Type:      eventType,
		Timestamp: time.Now().Format("15:04:05"),
		Data:      data,
	}
	
	bytes, err := json.Marshal(event)
	if err != nil {
		return
	}
	
	msg := string(bytes)
	
	Hub.mu.Lock()
	defer Hub.mu.Unlock()
	for ch := range Hub.clients {
		select {
		case ch <- msg:
		default:
		}
	}
}

func BroadcastLog(format string, v ...interface{}) {
	logMsg := fmt.Sprintf(format, v...)
	log.Println(logMsg)

	event := wsEvent{
		Type:      "log",
		Message:   logMsg,
		Timestamp: time.Now().Format("15:04:05"),
	}
	
	bytes, _ := json.Marshal(event)
	msg := string(bytes)

	Hub.mu.Lock()
	defer Hub.mu.Unlock()
	for ch := range Hub.clients {
		select {
		case ch <- msg:
		default:
		}
	}
}

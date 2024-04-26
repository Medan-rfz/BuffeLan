package sockethub

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	websock "buffelan/internal/network/websocket"
)

type WebsockHub struct {
	server  *websock.WebsockServer
	clients map[string]*websock.WebsockClient
}

func NewWebsockHub(listenPort uint16) (*WebsockHub, error) {
	server, err := websock.NewWebsockServer(websock.WebsockServerConfig{ListenPort: listenPort})
	if err != nil {
		return nil, err
	}

	hub := &WebsockHub{
		server:  server,
		clients: make(map[string]*websock.WebsockClient),
	}

	return hub, nil
}

func (h *WebsockHub) Serve(callback func(msg string)) {
	h.server.Serve(callback)
	defer h.server.Shutdown(context.Background())
}

func (h *WebsockHub) Close() {
	for _, v := range h.clients {
		v.Close()
	}
}

func (h *WebsockHub) AddClient(client *websock.WebsockClient) error {
	key := fmt.Sprintf("%s:%v", client.TargetHost, client.TargetPort)
	if _, ok := h.clients[key]; ok {
		return errors.New("client already connected")
	}

	h.clients[key] = client
	return nil
}

func (h *WebsockHub) CheckClientExists(client *websock.WebsockClient) bool {
	key := fmt.Sprintf("%s:%v", client.TargetHost, client.TargetPort)
	_, ok := h.clients[key]
	return ok
}

func (h *WebsockHub) SendMessage(msg string) {
	for _, c := range h.clients {
		err := c.Send([]byte(msg))
		if err != nil {
			log.Printf("client [%s] error: %v\n", c.TargetHost, err)
			key := fmt.Sprintf("%s:%v", c.TargetHost, c.TargetPort)
			delete(h.clients, key)
		}
	}
}

package websocket

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type WebsockClient struct {
	conn       *websocket.Conn
	TargetHost string
	TargetPort uint16
}

type WebsockClientConfig struct {
	TargetHost string
	TargetPort uint16
}

func NewWebsockClient(configs WebsockClientConfig) (*WebsockClient, error) {
	ws, err := websocket.Dial(
		fmt.Sprintf("ws://%s:%v/ws", configs.TargetHost, configs.TargetPort),
		"",
		fmt.Sprintf("http://%s/", configs.TargetHost))
	if err != nil {
		return nil, err
	}

	return &WebsockClient{
		conn:       ws,
		TargetHost: configs.TargetHost,
		TargetPort: configs.TargetPort,
	}, nil
}

func (c *WebsockClient) Close() {
	c.conn.Close()
}

func (c *WebsockClient) Send(message []byte) error {
	return websocket.Message.Send(c.conn, message)
}

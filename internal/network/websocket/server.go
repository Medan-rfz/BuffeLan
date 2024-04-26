package websocket

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

type WebsockServer struct {
	srv        *http.Server
	listenPort uint16
	listenCb   func(msg string)
}

type WebsockServerConfig struct {
	ListenPort uint16
}

func NewWebsockServer(configs WebsockServerConfig) (*WebsockServer, error) {
	return &WebsockServer{
		listenPort: configs.ListenPort,
	}, nil
}

func (s *WebsockServer) Serve(callback func(msg string)) error {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%v", s.listenPort))
	if err != nil {
		return err
	}
	defer conn.Close()

	mux := http.NewServeMux()
	mux.Handle("/ws", websocket.Handler(s.handleConnections))

	s.listenCb = callback
	s.srv = &http.Server{Handler: mux}

	log.Printf("Server starting... [0.0.0.0:%v]", s.listenPort)
	if err := s.srv.Serve(conn); err != nil {
		return err
	}

	return nil
}

func (s *WebsockServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *WebsockServer) handleConnections(ws *websocket.Conn) {
	defer ws.Close()

	for {
		var msg string

		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			log.Printf("Ошибка чтения сообщения: %v", err)
			break
		}

		s.listenCb(msg)

		// log.Printf("Получено сообщение: %s\n", msg)

		// // Отправляем обратно то же сообщение
		// err = websocket.Message.Send(ws, msg)
		// if err != nil {
		// 	log.Printf("Ошибка отправки сообщения: %v", err)
		// 	break
		// }
	}
}

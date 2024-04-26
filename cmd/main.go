package main

import (
	"context"
	"os"

	appstate "buffelan/internal/app_state"
	"buffelan/internal/beacon"
	"buffelan/internal/network/multicast"
	"buffelan/internal/network/sockethub"
	websock "buffelan/internal/network/websocket"
	"buffelan/internal/tray"

	"github.com/gen2brain/beeep"
	log "github.com/sirupsen/logrus"

	"golang.design/x/clipboard"
)

const (
	srvAddr         = "224.0.0.1:9999"
	listenerPort    = 7387
	maxDatagramSize = 8192
)

func main() {
	log.SetLevel(log.DebugLevel)

	ctx := context.Background()

	appState := &appstate.AppState{
		TxEnable: true,
		RxEnable: true,
	}

	config := tray.AppTrayConfig{
		CbExit: func() {
			os.Exit(0)
		},
	}

	multicastListener, err := multicast.NewMulticastListener(srvAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer multicastListener.Close()

	multicastSender, err := multicast.NewMulticastSender(srvAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer multicastSender.Close()

	hub, err := sockethub.NewWebsockHub(listenerPort)
	if err != nil {
		log.Fatalln(err)
	}
	defer hub.Close()

	tray := tray.NewAppTray(config, appState, multicastSender)

	go tray.Init()

	go func() {
		for msg := range multicastListener.Listen() {
			if msg.Data == beacon.MulticastBeaconMsg {
				log.Debugf("received beacon msg from: %v", msg.Src)

				config := websock.WebsockClientConfig{
					TargetHost: msg.Src.IP.String(),
					TargetPort: listenerPort,
				}
				client, err := websock.NewWebsockClient(config)
				if err != nil {
					log.Errorf("error connection to new client: %v:\n%v", msg.Src.IP.String(), err)
					continue
				}

				if hub.CheckClientExists(client) {
					continue
				}

				hub.AddClient(client)
				multicastSender.Send(beacon.MulticastBeaconMsg)
			}
		}
	}()

	prevData := ""

	go hub.Serve(func(msg string) {
		if !appState.RxEnable {
			return
		}

		prevData = msg
		clipboard.Write(clipboard.FmtText, []byte(msg))
		log.Debugf("added new value: %s\n", msg)
	})

	go func() {
		ch := clipboard.Watch(ctx, clipboard.FmtText)
		for data := range ch {
			if !appState.TxEnable {
				continue
			}

			msg := string(data)
			if prevData != msg {
				hub.SendMessage(msg)
				log.Debugf("new value sended all clients: %s\n", msg)
			}
		}
	}()

	multicastSender.Send(beacon.MulticastBeaconMsg)

	beeep.Notify("BuffeLan", "The application is working", "")
	log.Println("application ready")
	<-ctx.Done()
	log.Println("application shutdown")
}

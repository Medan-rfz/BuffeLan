package tray

import (
	appstate "buffelan/internal/app_state"
	"buffelan/internal/beacon"
	"buffelan/resources"

	"github.com/getlantern/systray"
)

type sender interface {
	Send(data string)
}

type AppTrayConfig struct {
	CbExit func()
	CbShot func()
}

type AppTray struct {
	config   AppTrayConfig
	sender   sender
	appState *appstate.AppState
}

func NewAppTray(config AppTrayConfig, appState *appstate.AppState, sender sender) *AppTray {
	return &AppTray{
		config:   config,
		sender:   sender,
		appState: appState,
	}
}

func (t *AppTray) Init() {
	systray.Run(t.onReady, t.onExit)
}

func (t *AppTray) onReady() {
	systray.SetIcon(resources.Icon)
	systray.SetTitle("BuffeLan")
	systray.SetTooltip("Shared copy buffer in LAN")

	mTxEn := systray.AddMenuItemCheckbox("Enable transmit", "Enable transmit share buffer", true)
	go t.txEnableButtonHandle(mTxEn)

	mRxEn := systray.AddMenuItemCheckbox("Enable receive", "Enable receive share buffer", true)
	go t.rxEnableButtonHandle(mRxEn)

	mFind := systray.AddMenuItem("Push find", "Send multicast message for force find all another devices")
	go t.fingButtonHandle(mFind)

	mQuit := systray.AddMenuItem("Quit", "Quit from the app")
	go t.exitButtonHandle(mQuit)
}

func (t *AppTray) onExit() {
}

func (t *AppTray) exitButtonHandle(item *systray.MenuItem) {
	for {
		<-item.ClickedCh
		if t.config.CbExit != nil {
			t.config.CbExit()
		}
	}
}

func (t *AppTray) fingButtonHandle(item *systray.MenuItem) {
	for {
		<-item.ClickedCh
		t.sender.Send(beacon.MulticastBeaconMsg)
	}
}

func (t *AppTray) txEnableButtonHandle(item *systray.MenuItem) {
	for {
		<-item.ClickedCh
		if item.Checked() {
			t.appState.TxEnable = false
			item.Uncheck()
		} else {
			t.appState.TxEnable = true
			item.Check()
		}
	}
}

func (t *AppTray) rxEnableButtonHandle(item *systray.MenuItem) {
	for {
		<-item.ClickedCh
		if item.Checked() {
			t.appState.RxEnable = false
			item.Uncheck()
		} else {
			t.appState.RxEnable = true
			item.Check()
		}
	}
}

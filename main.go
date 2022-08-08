package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type WinSize struct {
	Height uint16
	Width  uint16
	sync.RWMutex
}

type WinTty struct {
	session ssh.Session
	pty     ssh.Pty
	winch   <-chan ssh.Window
	size    *WinSize
}

func (w WinTty) Start() error {
	return nil
}

func newWinTty(s ssh.Session) (*WinTty, error) {
	var wtty WinTty

	wtty.session = s

	var isPty bool
	wtty.pty, wtty.winch, isPty = wtty.session.Pty()

	if !isPty {
		return nil, fmt.Errorf("not a pty")
	}

	var ws WinSize

	ws.Height = uint16(wtty.pty.Window.Height)
	ws.Width = uint16(wtty.pty.Window.Width)

	wtty.size = &ws

	return &wtty, nil
}

func (w WinTty) Read(p []byte) (n int, err error) {
	return w.session.Read(p)
}

func (w WinTty) Write(p []byte) (n int, err error) {
	return w.session.Write(p)
}

func (w WinTty) Stop() error {
	return nil
}

func (w WinTty) Drain() error {
	return nil
}

func (w WinTty) Close() error {
	return w.session.Close()
}

func (w WinTty) NotifyResize(cb func()) {
	// this does not work with more than one callback ...
	go func() {
		for win := range w.winch {
			w.size.Lock()
			w.size.Height = uint16(win.Height)
			w.size.Width = uint16(win.Width)
			w.size.Unlock()
			cb()
		}
	}()
}

func (w WinTty) WindowSize() (width int, height int, err error) {
	w.size.RLock()
	defer w.size.RUnlock()

	height = int(w.size.Height)
	width = int(w.size.Width)

	return width, height, nil
}

func handleWin(s ssh.Session, m *model) {
	var v view

	v.Lock()
	v.model = m

	p := &player{}
	p.setName(s.User())

	v.player = p

	m.addPlayer(p)

	wtty, err := newWinTty(s)
	if err != nil {
		log.Printf("No tty for user %s: %+v", s.User(), err)
		return
	}

	screen, err := tcell.NewTerminfoScreenFromTty(*wtty)
	if err != nil {
		panic(err)
	}

	screen.Init()

	v.Unlock()
	flex := v.flex()
	v.Lock()
	app := tview.NewApplication().SetScreen(screen).SetRoot(flex, true).EnableMouse(true)
	v.app = app
	v.Unlock()

	if err := app.Run(); err != nil {
		log.Printf("App for user %s crashed: %+v", s.User(), err)
	}

	s.Exit(0)
	m.delPlayer(s.User())
}

func main() {
	m := newModel()

	ssh.Handle(func(s ssh.Session) {
		handleWin(s, m)
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			players := m.getPlayers()
			for _, p := range players {
				fmt.Printf("Bye %s\n", p.getName())
				go m.delPlayer(p.getName())
			}
			fmt.Println("Shutting down")
			if len(players) > 0 {
				time.Sleep(2 * time.Second)
			}

			os.Exit(0)
		}
	}()

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil, ssh.HostKeyFile("id_ecdsa")))
}

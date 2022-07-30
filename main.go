package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type WinSize struct {
	Height uint16
	Width  uint16
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
			w.size.Height = uint16(win.Height)
			w.size.Width = uint16(win.Width)
			cb()
		}
	}()
}

func (w WinTty) WindowSize() (width int, height int, err error) {
	height = int(w.size.Height)
	width = int(w.size.Width)

	return width, height, nil
}

func handleWin(s ssh.Session) {

	fmt.Printf("new session: %+v\n", s)
	fmt.Printf("new user: %+v\n", s.User())
	fmt.Printf("new pubkey: %+v\n", s.PublicKey())

	wtty, err := newWinTty(s)
	if err != nil {
		return
	}

	screen, err := tcell.NewTerminfoScreenFromTty(*wtty)
	if err != nil {
		panic(err)
	}

	screen.Init()

	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := tview.NewApplication().SetScreen(screen).SetRoot(box, true).Run(); err != nil {
		panic(err)
	}

}

func main() {
	ssh.Handle(handleWin)

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return true // allow all keys, or use ssh.KeysEqual() to compare against known keys
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil, ssh.HostKeyFile("id_ecdsa"), publicKeyOption))
}

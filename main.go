package main

import (
	"github.com/murlokswarm/app"
	_ "github.com/murlokswarm/mac"
)

var (
	win app.Contexter
)

func main() {
	app.OnLaunch = func() {
		// Create the main window
		win = newMainWindow()
	}
	app.OnReopen = func() {
		if win != nil {
			return
		}
		win = newMainWindow()
	}
	app.Run()
}

func newMainWindow() app.Contexter {
	// Creates a window context.
	win := app.NewWindow(app.Window{
		Title:          "Pass",
		Width:          300,
		Height:         480,
		Vibrancy:       app.VibeUltraDark,
		TitlebarHidden: true,
		FixedSize:      true,
		OnClose: func() bool {
			win = nil
			return true
		},
	})
	// Create component...
	ps := &PasswordSearch{}
	// ...and mount to window
	win.Mount(ps)  
	// Return to context
	return win
}

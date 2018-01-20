package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
)

// AppMainMenu implements app.Componer interface.
type AppMainMenu struct {
	CustomTitle string
	Disabled    bool
}

// Render returns the HTML markup that describes the appearance of the
// component.
// In this case, the component will be mounted into a menu context.
// This restrict the markup to a compositon of menu and menuitem.
func (m *AppMainMenu) Render() string {
	return `
<menu>
    <menu label="app">
        <menuitem label="{{if .CustomTitle}}{{.CustomTitle}}{{else}}Custom item{{end}}" 
                  onclick="OnCustomMenuClick" 
                  icon="star.png"
                  separator="true"
                  disabled="{{.Disabled}}" />
        <menuitem label="Quit" shortcut="meta+q" selector="terminate:" />        
    </menu>
    <WindowMenu />
</menu>
    `
}

// OnCustomMenuClick is the handler called when an onclick event occurs in a menuitem.
func (m *AppMainMenu) OnCustomMenuClick() {
	log.Info("OnCustomMenuClick")
}

// WindowMenu implements app.Componer interface.
// It's another component which will be nested inside the AppMenu component.
type WindowMenu struct {
}

func (m *WindowMenu) Render() string {
	return `
<menu label="Accounts">
    <menuitem label="Generate" selector="performClose:" shortcut="meta+w" />
    <menuitem label="Close" selector="performClose:" shortcut="meta+w" />
</menu>
    `
}

func init() {
	// Allows the app to create a AppMainMenu and WindowMenu components when it finds its declaration
	// into a HTML markup.
	app.RegisterComponent(&AppMainMenu{})
	app.RegisterComponent(&WindowMenu{})
}

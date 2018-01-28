/*
Copyright (c) 2018 Carl Löndahl. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Pass Desktop nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
    "github.com/murlokswarm/app"
    _ "github.com/murlokswarm/log"
)

type AppMainMenu struct {
}

func (m *AppMainMenu) Render() string {
    return `
<menu>
    <menu label="app">
        <menuitem label="About Pass Desktop" 
                  onclick="ShowAboutView" 
                  separator="true" />
        <menuitem label="Search" 
                  shortcut="meta+f" 
                  onclick="ShowSearchView" />
        <menuitem label="Add Account" 
                  shortcut="meta+n"
                  onclick="ShowAddView" 
                  separator="true" />
        <menuitem label="Quit" shortcut="meta+q" selector="terminate:" />     
    </menu>
    <EditMenu />
</menu>`
}

func (m *AppMainMenu) ShowSearchView() {
    pass.CurrentView = ViewSearchDialog
    app.Render(win.Component())
}

func (m *AppMainMenu) ShowAddView() {
    // Clear internal struct, since we are starting fresh.
    pass.Account.Name = ""
    pass.Account.Username = ""
    pass.Account.Password = ""
    // Set add dialog.
    pass.CurrentView = ViewAddAccountDialog
    app.Render(win.Component())
}

func (m *AppMainMenu) ShowAboutView() {
    pass.CurrentView = ViewAboutDialog
    app.Render(win.Component())
}

type EditMenu struct {
}

func (m *EditMenu) Render() string {
    return `
<menu label="Edit">
    <menuitem label="Undo" selector="undo:" shortcut="meta+z" />
    <menuitem label="Redo" selector="redo:" shortcut="meta+shift+z" />
    <menuitem label="Cut" selector="cut:" shortcut="meta+x" />
    <menuitem label="Copy" selector="copy:" shortcut="meta+c" />
    <menuitem label="Paste" selector="pasteAsPlainText:" shortcut="meta+v" />
    <menuitem label="Select All" selector="selectAll:" shortcut="meta+a" />
</menu>`
}

func init() {
    app.RegisterComponent(&AppMainMenu{})
    app.RegisterComponent(&EditMenu{})
}

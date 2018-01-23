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
   * Neither the name of Google Inc. nor the names of its
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
        Height:         550,
        Vibrancy:       app.VibeUltraDark,
        //Vibrancy:       app.VibeMediumLight,
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

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
    "os"
    "log"
    "fmt"
    "time"
    "path/filepath"
    "net/http"
    "crypto/x509"
    "crypto/tls"
    "github.com/murlokswarm/app"
    _ "github.com/murlokswarm/mac"
)

var (
    win app.Contexter
    pass Application
)

const (
    ViewSearchDialog = 0
    ViewSearchClearedDialog = 1
    ViewAccountDialog = 2
    ViewConfirmDeleteDialog = 3
    ViewCreateAccountDialog = 4
    ViewUnlockDialog = 5
    ViewAboutDialog = 6
)

const (
    UseArgon2ForKeyDerivation = false
    DefaultGeneratedPasswordLength = 32
    ConfigFile = "/config/config.json"
    //ConfigFile = "/../Resources/config/config.json"
)

func ConfigureTLSClient() {
    // Setup entrypoint
    pass.EntryPoint = fmt.Sprintf("https://%s:%s/v1/secret", 
                                  pass.Config.Host, 
                                  pass.Config.Port)

    // Create a TLS context...
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM([]byte(pass.Config.CA))

    // ...and a client
    pass.Client = &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:      caCertPool,
            },
        },
        Timeout: time.Second * 10,
    }
}

func SetApplicationPath() {
    dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
    pass.FullPath = string(dir)
}

func main() {
    // Pass is locked by default.
    pass.Locked = true

    // Get current directory to read config and icons.
    SetApplicationPath()

    // Load config file.
    if _, err := os.Stat(pass.FullPath + ConfigFile); os.IsNotExist(err) {
        log.Println("No config file present.")
    } else {
        pass.Config = LoadConfiguration(pass.FullPath + ConfigFile)
        ConfigureTLSClient()
    }

    app.OnLaunch = func() {
        // Creates the AppMainMenu component.
        appMenu := &AppMainMenu{}

        // Mounts the AppMainMenu component into the application menu bar.
        if menuBar, ok := app.MenuBar(); ok {
            menuBar.Mount(appMenu)
        }

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
        Height:         548,
        Vibrancy:       app.VibeDark,
        TitlebarHidden: true,
        //FixedSize:      true,
        OnClose: func() bool {
            win = nil
            return true
        },
    })

    // Create component...
    ps := &PassView{}
   
    // ...and mount to window
    win.Mount(ps)  
    
    // Return to context
    return win
}

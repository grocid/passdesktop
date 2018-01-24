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
    "strconv"
    "path/filepath"
    "net/url"
    "net/http"
    "io/ioutil"
    "crypto/x509"
    "crypto/tls"
    "github.com/murlokswarm/app"
)

type PassView struct {
    Query          string 
    Account        string
    Token          string
    Username       string
    Password       string
    PasswordAgain  string
    Hostname       string
    Port           string
    Filename       string
}

type AccountInfo struct {
    Name           string
    Username       string
    Password       string
}

type Application struct {
    Client         *http.Client
    Config         Configuration
    DecryptedToken string
    Locked         bool
    Account        AccountInfo
    CurrentView    int
    SearchResult   []string
    EntryPoint     string
}

const ViewSearchDialog  = 0
const ViewAccountDialog = 1
const ViewConfirmDeleteDialog = 2
const ViewCreateAccountDialog = 3
const ViewUnlockDialog = 4

const DefaultGeneratedPasswordLength = 32
const ConfigFile = "config.json"
const CAFile = "ca.crt"

var pass Application

func ClearAccountInformation(h *PassView) {
    h.Account = ""
    h.Username = ""
    h.Password = ""
    return
}

func (h *PassView) Render() string {
    // Clear all account data
    ClearAccountInformation(h)

    // If no config file is present...
    if pass.Config == (Configuration{}) {
        return GetCreateConfigDialog()
    }

    // If locked ask for password...
    if pass.Locked || pass.DecryptedToken == "" {
        return GetPasswordInput()
    }

    // Get the view from CurrentView and display
    // accordingly
    switch pass.CurrentView {

        case ViewSearchDialog:
            if len(pass.SearchResult) > 0 {
                return GetListBody(pass.SearchResult)
            } else {
                return GetEmptySearchDialog()
            }

        case ViewAccountDialog:
            h.Account = pass.Account.Name
            h.Username = pass.Account.Username
            h.Password = pass.Account.Password
            return GetAccountBody(pass.Account.Name)

        case ViewConfirmDeleteDialog:
            return GetConfirmDeleteDialog()

        default:
            log.Fatal(pass.CurrentView)
            return ""
    }
}

func (h *PassView) Unlock(arg app.ChangeArg) {
    // Get the password
    password := arg.Value

    // Try to unlock the token with the password
    token, err := UnlockToken(pass.Config.Encrypted.Token, 
                              password, 
                              pass.Config.Encrypted.Nonce, 
                              pass.Config.Encrypted.Salt)

    // If we succed with message authentication, i.e.,
    // if password is correct...
    if err == nil {
        // Progress to unlocked display
        pass.Locked = false
        // Set token
        pass.DecryptedToken = token
        // Init the unlocked display with all entries
        pass.SearchResult = DoListRequest("")
    }

    // Preset the search dialog
    pass.CurrentView = ViewSearchDialog
    
    // Show if unlocking was successful
    app.Render(h)
}

func (h *PassView) PickFile(arg app.ChangeArg) {
    app.NewFilePicker(app.FilePicker {
            MultipleSelection: false,
            NoDir:             true,
            NoFile:            false,
            OnPick: func(filenames []string) {
                CopyFile(filenames[0], CAFile)
                h.Filename = filepath.Base(filenames[0])
                app.Render(h)
            },
        })

}

func (h *PassView) CreateConfig(arg app.ChangeArg) {
    if h.Query == "" {
        log.Println("No token.")
        return
    }

    if (len(h.Password) < 4) && (h.Password != h.PasswordAgain) {
        log.Println("Password too short or non-matching.")
        return
    }

    if h.Hostname == "" {
        log.Println("No hostname.")
        return
    }

    if h.Port == "" {
        log.Println("No port.")
        return
    }

    if _, err := strconv.Atoi(h.Port); err != nil {
        log.Println("Invalid port.")
        return
    }

    token, nonce, salt, err := LockToken(h.Token, h.Password)
    if err != nil {
        log.Println(err)
        return
    }

    if _, err := os.Stat(CAFile); os.IsNotExist(err) {
        log.Println(err)
        return
    }

    // Put data into config struct
    pass.Config.Encrypted.Token = token
    pass.Config.Encrypted.Nonce = nonce
    pass.Config.Encrypted.Salt = salt
    pass.Config.Host = h.Hostname
    pass.Config.Port = h.Port

    // Lock
    pass.Locked = true
    pass.DecryptedToken = ""

    // Wipe data from GUI
    *h = PassView{}

    // Setup the client
    ConfigureTLSClient()

    app.Render(h)
}

func (h *PassView) OnHref(URL *url.URL) {
    // Extract information from query and get account name from query.
    u := URL.Query()
    h.Query = u.Get("Account")

    // Go to account view and fetch information from server
    pass.CurrentView = ViewAccountDialog
    pass.Account = DoGetRequest(h.Query)

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) CancelTrashView(arg app.ChangeArg) {
    // Go back from trash view to account view
    pass.CurrentView = ViewAccountDialog

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) CancelAccountView(arg app.ChangeArg) {
    // Go back from account view to search view
    pass.CurrentView = ViewSearchDialog

    // Go back from account view
    pass.SearchResult = DoListRequest("")

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) OkAccountView(arg app.ChangeArg) {
    // Go back from account view to search view and fetch all accounts
    pass.CurrentView = ViewSearchDialog
    pass.SearchResult = DoListRequest("")

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) RerandomizePasswordAccountView(arg app.ChangeArg) {
    // Generate a new password
    h.Password, _ = RandomPassword(DefaultGeneratedPasswordLength)

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) DeleteAccountView(arg app.ChangeArg) {
    // Go to confirm delete view
    pass.CurrentView = ViewConfirmDeleteDialog

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) Search(arg app.ChangeArg) {
    // Stay in search view and fetch accounts based on query
    pass.CurrentView = ViewSearchDialog
    pass.SearchResult = DoListRequest(arg.Value)

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func ConfigureTLSClient() {
    // Setup entrypoint
    pass.EntryPoint = fmt.Sprintf("https://%s:%s/v1/secret", 
                                  pass.Config.Host, 
                                  pass.Config.Port)

    // Init client for performing REST queries to Vault
    caCert, err := ioutil.ReadFile(CAFile)

    // Unless something went wrong with reading the certificate...
    if err != nil {
        log.Fatal(err)
    }

    // Create a TLS context...
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // ...and a client
    pass.Client = &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:      caCertPool,
            },
        },
    }
}

func init() {
    // Pass is locked by default
    pass.Locked = true

    // Load config file
    if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
        log.Println("No config file present.")
    } else {
        pass.Config = LoadConfiguration(ConfigFile)
        ConfigureTLSClient()
    }

    // Register UI component
    app.RegisterComponent(&PassView{})
}

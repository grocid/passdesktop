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
    "strconv"
    "path/filepath"
    "net/url"
    "net/http"
    "io/ioutil"
    "crypto/x509"
    "crypto/tls"
    "github.com/murlokswarm/app"
)

var decryptedToken string
var client *http.Client
var entryPoint string
var unlocked bool
var config Config

const ConfigFile = "config.json"
const CAFile = "ca.crt"

type PasswordSearch struct {
    Query          string 
    Account        string
    Username       string
    Password       string
    List           []string
    UnlockPassword string
    View           int
    Meta           string
    Filename       string
}

type CurrentAccount struct {
    Account        string
    Username       string
    Password       string
}

func (h *PasswordSearch) Render() string {
    // If no config file is present...
    if config == (Config{}) {
        return GetCreateConfigDialog()
    }
    // If locked ask for password...
    if !unlocked {
        //if h.UnlockPassword != "nil" TODO
        return GetPasswordInput()
    }
    // If the client has been unlocked, show the screen with
    // a search bar and a list with entries
    if h.View == 1 {
        // If an account has been specified, i.e., if an entry was clicked...
        return GetAccountBody(h)
    } else if h.View == 0 {
        // ...or show the filtered list...
        if len(h.List) > 0 {
            return GetListBody(h)
        } else {
            return GetEmptySearchDialog()
        }
    } else {
        return GetConfirmDeleteDialog()
    }

}

// Does a search and filter
func GetList(s string) []string {
    return DoListRequest(s)
}

// Requests a specific entry
func GetEntry(s string) (string, string) {
    return DoGetRequest(s)
}

func (h *PasswordSearch) OnHref(URL *url.URL) {
    // Extract information from query
    u := URL.Query()
    // Set account and username/password in internal
    // data structure for showing
    h.Account = u.Get("Account")
    h.Query = u.Get("Account")
    h.View = 1
    h.Password, h.Username = DoGetRequest(h.Account)
    // Show it!
    app.Render(h)
}

func (h *PasswordSearch) Unlock(arg app.ChangeArg) {
    // Get the password
    password := arg.Value
    // Try to unlock the token with the password
    token, err := UnlockToken(config.Encrypted.Token, 
        password, config.Encrypted.Nonce, config.Encrypted.Salt)
    // If we succed with message authentication, i.e.,
    // if password is correct...
    if err == nil {
        // Progress to unlocked display
        unlocked = true
        log.Println("Unlocked")
        // Set token
        decryptedToken = token
        // Init the unlocked display with all entries
        h.List = GetList("")
    }
    h.View = 0
    // Show
    app.Render(h)
}

func (h *PasswordSearch) PickFile(arg app.ChangeArg) {
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

func (h *PasswordSearch) CreateConfig(arg app.ChangeArg) {
    if h.Query == "" {
        log.Println("No token.")
        return
    }
    // Passwords
    if (len(h.Password) < 4) && (h.Password != h.UnlockPassword) {
        log.Println("Password too short or non-matching.")
        return
    }
    // Hostname
    if h.Account == "" {
        log.Println("No hostname.")
        return
    }
    // Port
    if h.Meta == "" {
        log.Println("No port.")
        return
    }
    if _, err := strconv.Atoi(h.Meta); err != nil {
        log.Println("Invalid port.")
        return
    }
    token, nonce, salt, err := LockToken(h.Query, h.Password)
    if err != nil {
        log.Println(err)
        return
    }
    if _, err := os.Stat(CAFile); os.IsNotExist(err) {
        log.Println(err)
        return
    }
    config.Encrypted.Token = token
    config.Encrypted.Nonce = nonce
    config.Encrypted.Salt = salt
    config.Host = h.Account
    config.Port = h.Meta
    //unlocked = true
    log.Println(h.Password)
    log.Println(token)
    log.Println(nonce)
    log.Println(salt)
    //decryptedToken = h.Query
    h.Query = ""
    h.Password = ""
    h.UnlockPassword = ""
    SetupConfig()
    app.Render(h)
}

func (h *PasswordSearch) CancelTrashView(arg app.ChangeArg) {
    h.View = 1
    app.Render(h)
}

func (h *PasswordSearch) CancelAccountView(arg app.ChangeArg) {
    // Go back from account view
    h.List = GetList("")
    h.View = 0
    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PasswordSearch) OkAccountView(arg app.ChangeArg) {
    // Go back from account view
    h.List = GetList("")
    h.View = 0
    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PasswordSearch) RerandomizePasswordAccountView(arg app.ChangeArg) {
    h.Password, _ = RandomPassword(64)
    h.View = 1
    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PasswordSearch) DeleteAccountView(arg app.ChangeArg) {
    // Go back from account view
    h.List = GetList("")
    // Tells the app to update the rendering of the component.
    h.View = 2
    app.Render(h)
}

func (h *PasswordSearch) Search(arg app.ChangeArg) {
    h.Query = arg.Value
    // Create a list of entries
    h.List = GetList(h.Query)
    // Tells the app to update the rendering of the component.
    h.View = 0
    app.Render(h)
}

func SetupConfig() {
    entryPoint = "https://" + config.Host + ":" + config.Port + "/v1/secret"
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
    client = &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:      caCertPool,
            },
        },
    }
}

func init() {

    // Load config file
    if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
        log.Println("No config file present.")
    } else {
        config = LoadConfiguration(ConfigFile)
        SetupConfig()
    }

    // Register UI component
    app.RegisterComponent(&PasswordSearch{})
}

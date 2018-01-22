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
    "log"
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


type PasswordSearch struct {
    Query          string 
    Account        string
    Username       string
    Password       string
    List           []string
    UnlockPassword string
    View           int
}

type CurrentAccount struct {
    Account        string
    Username       string
    Password       string
}

func CreateAccountBody() string {
    return ``
}

func (h *PasswordSearch) Render() string {
    // If the client has been unlocked, show the screen with
    // a search bar and a list with entries 
    if unlocked {
        var body string
        var head string
        // If an account has been specified, i.e., if an entry was clicked...
        if h.View == 1 {
            head = GetSearchInput()
            body = GetAccountBody()
        } else if h.View == 0 {
            // ...or show the filtered list...
            if len(h.List) > 0 {
                head = GetSearchInput()
                body = GetListBody()
            } else {
                  return GetEmptySearchDialog()
            }
        } else {
            return GetConfirmDeleteDialog()
        }
        return head + body
    } else {
        //if h.UnlockPassword != "nil" TODO
        return GetPasswordInput()
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
        // Set token
        decryptedToken = token
        // Init the unlocked display with all entries
        h.List = GetList("")
    }
    h.View = 0
    // Show
    app.Render(h)
}

func (h *PasswordSearch) CancelTrashView(arg app.ChangeArg) {
    h.View = 1
    app.Render(h)
}

func (h *PasswordSearch) CancelAccountView(arg app.ChangeArg) {
    log.Println(h.Password)
    log.Println(h.Username)
    // Go back from account view
    h.List = GetList("")
    h.View = 0
    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PasswordSearch) OkAccountView(arg app.ChangeArg) {
    log.Println(h.Password)
    log.Println(h.Username)
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

func (h *PasswordSearch) OnContextMenu() {
    ///ctxmenu := app.NewContextMenu()
    //ctxmenu.Mount(&AppMainMenu{})
}

func init() {
    // Init client for performing REST queries to Vault
    caCert, err := ioutil.ReadFile("ca.crt")
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
    // Load config file
    config = LoadConfiguration("config.json")
    entryPoint = "https://" + config.Host + ":" + config.Port + "/v1/secret"
    // Register UI component
    app.RegisterComponent(&PasswordSearch{})
}

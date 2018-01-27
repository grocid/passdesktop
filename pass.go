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
    "log"
    "net/http"
    "strconv"
    "path/filepath"
    "net/url"
    "github.com/murlokswarm/app"
)

type (



    Application struct {
        Client         *http.Client
        Config         Configuration
        DecryptedToken string
        EncryptionKey  []byte
        Locked         bool

        Account        AccountInfo
        CurrentView    int
        SearchResult   []Entry // Encrypted part
        EntryPoint     string
        FullPath       string
    }

    PassView struct {
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

    AccountInfo struct {
        Name           string
        Username       string
        Password       string
        Encrypted      string
    }
)

func (h *PassView) Render() string {


    // If no config file is present...
    if pass.Config == (Configuration{}) {
        return GetCreateConfigDialog()
    }

    // If locked ask for password...
    if pass.Locked || pass.DecryptedToken == "" {
        log.Println("Locked")
        return GetPasswordInput()
    }
    
    // Clear all account data from UI
    AccountClearInformation(h)

    // This is a MUX for views.
    // Get the view from CurrentView and display
    // accordingly.
    switch pass.CurrentView {

        // Show search dialog.
        case ViewSearchDialog:
            if len(pass.SearchResult) > 0 {
                // Show results from search.
                return GetListBody(pass.SearchResult)
            } else {
                // Whenever no match is found, show
                // search glass and a button to add
                // that particular account. .

                return GetEmptySearchDialog()

            }

        // Show the full list, with not queries filtered.
        case ViewSearchClearedDialog:
            h.Query = ""
            return GetListBody(pass.SearchResult)

        // When a list item is clicked, bring up account view.
        case ViewAccountDialog:
            // Pass information from internal struct
            // to the UI components.
            h.Account = pass.Account.Name
            h.Username = pass.Account.Username
            h.Password = pass.Account.Password

            return GetAccountBody(pass.Account.Name)

        // When an the add item is clicked during an empty 
        // search result, (we are now in GetEmptySearchDialog) 
        // we generate an account based on that query. Since
        // a presumably interesting name was searched for
        // we set the name on forehand with this guess.
        case ViewCreateAccountDialog:
            h.Account = h.Query
            return GetAddDialog(h.Query)

        // When pressing add account in menu, the name
        // is not determined. We need the view with an
        // editable name.
        case ViewAddAccountDialog:
            h.Query = ""
            return GetAddDialog(h.Query)

        // A dialog to confirm deletion of account.
        case ViewConfirmDeleteDialog:
            return GetConfirmDeleteDialog()

        // Just some cool about dialog.
        case ViewAboutDialog:
            return GetAboutDialog()

        default:
            log.Println("No such window")
            log.Fatal(pass.CurrentView)
            return ""
    }
}

func (h *PassView) Unlock(arg app.ChangeArg) {
    // Get the password.
    password := arg.Value

    // Try to unlock the token with the password
    token, key, err := UnlockToken(pass.Config.Encrypted.Token, 
                                   password,
                                   pass.Config.Encrypted.Salt)

    // If we succed with message authentication, i.e.,
    // if password is correct...
    if err == nil {
        // Progress to unlocked display
        pass.Locked = false
        // Set token
        pass.DecryptedToken = token
        pass.EncryptionKey = key
        // Init the unlocked display with all entries
        pass.SearchResult = DoListRequest("")
    }

    log.Println(err)

    // Preset the search dialog
    pass.CurrentView = ViewSearchDialog
    
    // Show if unlocking was successful
    app.Render(h)
}

func (h *PassView) Search(arg app.ChangeArg) {
    // Stay in search view and fetch accounts based on query
    pass.CurrentView = ViewSearchDialog
    pass.SearchResult = DoListRequest(arg.Value)

    // We need to keep track of the query, since if result is empty
    // and we want to create account, this is where information is
    // fetched from.
    h.Query = arg.Value

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) OnHref(URL *url.URL) {
    // Extract information from query and get account name from query.
    u := URL.Query()
    h.Query = u.Get("Account")
    encryptedName := u.Get("Encrypted")

    // Go to account view and fetch information from server
    pass.CurrentView = ViewAccountDialog
    pass.Account = DoGetRequest(Entry{h.Query, encryptedName})

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) AccountCreate(arg app.ChangeArg) {
    // Initialize the query as name and
    // with empty crendentials
    h.Account = h.Query
    pass.Account.Name = h.Query
    pass.Account.Username = ""
    pass.Account.Password = ""
    // We need to set encrypted to "", since we are creating a new
    // account -- otherwise it will map to another one.
    pass.Account.Encrypted = ""

    pass.CurrentView = ViewCreateAccountDialog

    app.Render(h)
}

func (h *PassView) AccountAddOk(arg app.ChangeArg) {
    // Update internal struct.
    pass.Account.Name = h.Account
    pass.Account.Encrypted = "" // We are creating new entry.
    pass.Account.Username = h.Username
    pass.Account.Password = h.Password
    log.Println(pass.Account)
    if pass.Account.Name != "" {
        AccountUpdate(h)
    } else {
        log.Println("Empty name!")
    }
}

func (h *PassView) AccountOk(arg app.ChangeArg) {
    AccountUpdate(h)
}

func (h *PassView) AccountTrashOk(arg app.ChangeArg) {
    // Send delete request to REST.
    DoDeleteRequest(pass.Account) 

    // Since we deleted the account, we remove its name
    // from the search bar
    h.Query = ""

    // Go back from account view to search view
    pass.CurrentView = ViewSearchDialog

    // Go back from account view
    pass.SearchResult = DoListRequest("")

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) AccountTrashCancel(arg app.ChangeArg) {
    // Go back from trash view to account view
    pass.CurrentView = ViewAccountDialog

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) AccountCancel(arg app.ChangeArg) {
    // Go back from account view to search view.
    pass.CurrentView = ViewSearchDialog

    // Go back from account view
    pass.SearchResult = DoListRequest("")

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) AccountRandomizePassword(arg app.ChangeArg) {
    // Generate a new password.
    pass.Account.Password, _ = RandomPassword(DefaultGeneratedPasswordLength)

    // Update name and username from UI.
    pass.Account.Name = h.Account
    pass.Account.Username = h.Username

    // Update UI with new password
    h.Password = pass.Account.Password

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) AccountDelete(arg app.ChangeArg) {
    // Go to confirm delete view.
    pass.CurrentView = ViewConfirmDeleteDialog

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func AccountClearInformation(h *PassView) {
    h.Account = ""
    h.Username = ""
    h.Password = ""

    return
}

func AccountUpdate(h *PassView) {
    // Update internal struct holding information
    // from UI.
    pass.Account.Username = h.Username
    pass.Account.Password = h.Password
    // Account.Username must remain unchanged or we will
    // create a new entry

    // Update it on server.
    DoPutRequest(pass.Account) 

    // Go back from account view to search view and fetch all accounts
    pass.CurrentView = ViewSearchDialog
    pass.SearchResult = DoListRequest("")

    app.Render(h)
}

func (h *PassView) PickFile(arg app.ChangeArg) {
    app.NewFilePicker(app.FilePicker {
            MultipleSelection: false,
            NoDir:             true,
            NoFile:            false,
            OnPick: func(filenames []string) {
                h.Filename = filepath.Base(filenames[0])
                app.Render(h)
            },
        })
}

func (h *PassView) CreateConfig(arg app.ChangeArg) {
    if h.Token == "" {
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

    token, salt, err := LockToken(h.Token, h.Password)
    if err != nil {
        log.Println(err)
        return
    }

    // Put data into config struct
    pass.Config.Encrypted.Token = token
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

func init() {
    // Register UI component
    app.RegisterComponent(&PassView{})
}

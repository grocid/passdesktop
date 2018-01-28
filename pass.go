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
    "log"
    "net/url"
    "path/filepath"
    "strconv"
)

// To remember the query we did before we entered an account.
// When pressing back, we want to return to that query.
var previousQuery string

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

    // This is a MUX for views
    // Get the view from CurrentView and display
    // accordingly.
    switch pass.CurrentView {

    // Show search dialog.
    case ViewSearchDialog:
        // Clear all account data from UI
        AccountClearInformation(h)

        filteredList := VaultListSecrets(h.Query)
        if len(filteredList) > 0 {
            // Show results from search.
            return GetListBody(filteredList)
        } else {
            // Whenever no match is found, show
            // search glass and a button to add
            // that particular account.
            return GetEmptySearchDialog()
        }

    // Show the full list, with not queries filtered.
    // This occurs when shortcut is invoked.
    case ViewSearchClearedDialog:
        // Clear all account data from UI
        AccountClearInformation(h)

        // Clear the query
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
    case ViewPresetAddAccountDialog:
        AccountClearInformation(h)
        h.Account = h.Query
        pass.CurrentView = ViewStayAddAccountDialog
        return GetAddDialog(h.Query)

    // When pressing add account in menu, the name
    // is not determined. We need the view with an
    // editable name.
    case ViewAddAccountDialog:
        h.Query = ""
        AccountClearInformation(h)
        pass.CurrentView = ViewStayAddAccountDialog
        return GetAddDialog(h.Query)

    // Essentially used only when the add account
    // dialog navigates back to itself, i.e., when
    // randomizing password.
    case ViewStayAddAccountDialog:
        return GetAddDialog(h.Query)

    // A dialog to confirm deletion of account.
    case ViewConfirmDeleteDialog:
        return GetConfirmDeleteDialog()

    case ViewSecureFileDialog:
        return GetSecureFileDialog()

    // Just some cool about-this-program dialog.
    case ViewAboutDialog:
        // Clear all account data from UI
        AccountClearInformation(h)

        return GetAboutDialog()

    default:
        // Clear all account data from UI
        AccountClearInformation(h)

        log.Println("No such window")
        log.Fatal(pass.CurrentView)
        return ""
    }
}

func (h *PassView) OnHref(URL *url.URL) {
    // Extract information from query and get account name and
    // its encrypted counterpart from query (this is need since
    // if we were to encrypt again, we would get a different
    // encrypted name).
    u := URL.Query()
    h.Query = u.Get("Account")
    encryptedName := u.Get("Encrypted")

    // Go to account view and fetch information from server
    pass.CurrentView = ViewAccountDialog
    pass.Account = VaultReadSecret(Entry{h.Query, encryptedName})

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) Search(arg app.ChangeArg) {
    // Stay in search view...
    pass.CurrentView = ViewSearchDialog

    // ...and fetch accounts based on query.
    //VaultListSecrets("")

    // We need to keep track of the query, since if result is empty
    // and we want to create account, this is where information is
    // fetched from.
    h.Query = arg.Value
    previousQuery = h.Query

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) AccountCreate(arg app.ChangeArg) {
    // Initialize the query as name and with empty crendentials.
    log.Println("AccountCreate")
    h.Account = h.Query
    pass.Account.Name = h.Query
    pass.Account.Username = ""
    pass.Account.Password = ""

    // We need to set encrypted to "", since we are creating a new
    // account -- otherwise it will map to another one.
    pass.Account.Encrypted = ""

    // Set the view.
    pass.CurrentView = ViewPresetAddAccountDialog

    app.Render(h)
}

func (h *PassView) AccountAddOk(arg app.ChangeArg) {
    // Update internal struct.
    pass.Account.Name = h.Account
    pass.Account.Encrypted = "" // We are creating new entry.
    pass.Account.Username = h.Username
    pass.Account.Password = h.Password

    // We will not accept empty names.
    if pass.Account.Name != "" {
        AccountUpdate(h)
    } else {
        log.Println("Empty name!")
    }

    app.Render(h)
}

func (h *PassView) AccountOk(arg app.ChangeArg) {
    h.Query = previousQuery
    AccountUpdate(h)

}

func (h *PassView) AccountTrashOk(arg app.ChangeArg) {
    // Send delete request to REST.
    VaultDeleteSecret(pass.Account)

    // Since we deleted the account, we remove its name
    // from the search bar.
    h.Query = previousQuery

    // Go back from account view to search view.
    pass.CurrentView = ViewSearchDialog

    // Get full list.
    VaultListSecrets(previousQuery)

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) AccountTrashCancel(arg app.ChangeArg) {
    // Go back from trash view to account view.
    pass.CurrentView = ViewAccountDialog

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) AccountCancel(arg app.ChangeArg) {
    // Go back from account view to search view.
    pass.CurrentView = ViewSearchDialog

    // Get full list.
    h.Query = previousQuery
    VaultListSecrets(previousQuery)

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *PassView) AccountRandomizePassword(arg app.ChangeArg) {

    // Generate a new password.
    pass.Account.Password, _ = Entropy(DefaultGeneratedPasswordLength)

    // Update UI with new password
    h.Password = pass.Account.Password

    // Update name and username from UI.
    pass.Account.Name = h.Account
    pass.Account.Username = h.Username

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

    // Update it on server.
    VaultWriteSecret(pass.Account)

    // Go back from account view to search view and fetch all accounts
    pass.CurrentView = ViewSearchDialog

    VaultListSecrets("")

    app.Render(h)
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

        log.Println("Unlocked")

        // Set token
        pass.DecryptedToken = token
        pass.EncryptionKey = key

        // Init the unlocked display with all entries
        VaultListSecrets("")

        // Preset the search dialog
        pass.CurrentView = ViewSecureFileDialog
    } else {
        log.Println(err)
    }

    // Show if unlocking was successful
    app.Render(h)
}

func (h *PassView) PickFile(arg app.ChangeArg) {
    app.NewFilePicker(app.FilePicker{
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

    // Require a minimum length on the password and
    // that the user actually knew what he/she typed in.
    if len(h.Password) < MinimumPasswordLength && h.Password != h.PasswordAgain {
        log.Println("Password too short or non-matching.")
        return
    }

    // An empty hostname will not do. We will not perform any
    // verification of the hostname.
    if h.Hostname == "" {
        log.Println("No hostname.")
        return
    }

    // Port needs to be non-empty...
    if h.Port == "" {
        log.Println("No port.")
        return
    }

    // ...as well as an integer for all I know.
    if _, err := strconv.Atoi(h.Port); err != nil {
        log.Println("Invalid port.")
        return
    }

    // Try to encrypt the token with the password.
    token, salt, err := LockToken(h.Token, h.Password)
    if err != nil {
        log.Println(err)
        return
    }

    // Put data into config struct.
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

package main

import (
    "log"
    "net/url"
    "net/http"
    "io/ioutil"
    "crypto/x509"
    "crypto/tls"
    "strings"
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
    List           string
    UnlockPassword string
}

func (h *PasswordSearch) Render() string {
    // If the client has been unlocked, show the screen with
    // a search bar and a list with entries 
    if unlocked {
        var body string
        head := `<div class="WindowLayout" oncontextmenu="OnContextMenu">    
                    <div class="SearchLayout">
                        <input type="text"
                               value="{{html .Query}}"
                               placeholder="Account"
                               autofocus="true"
                               onchange="Search"
                               autocomplete="off" 
                               autocorrect="off" 
                               autocapitalize="off" 
                               spellcheck="false"
                               selectable="on" 
                               class="selectable"/>`
        // If an account has been specified, i.e., if an entry was clicked...
        if h.Account != "" {
            body = `<div  style="overflow-y:scroll; max-height:375px; margin-top:20px" clickable="on" class="clickable">
                        <h1 style="margin-top: 0; margin-left: 15px">{{.Account}}</h1>
                        <h2 style="margin-left: 15px">Username</h2>
                        <div selectable="on" class="selectable" style="color: PaleGreen; margin-top: 0; margin-left: 15px">{{.Username}}</div>
                        <h2 style="margin-left: 15px">Password</h2>
                        <div selectable="on" class="selectable" style="color: DeepPink; margin-top: 0; margin-left: 15px; -webkit-user-select: auto;
   pointer-events: auto;">{{.Password}}</div>
                    </div>`
        } else {
            // ...or show the filtered list...
            body = `<div  style="overflow-y:scroll; max-height:375px; margin-top:20px" clickable="on" class="clickable">
                        {{.List}}
                    </div>`
        }
        tail := `   </div>
                </div>`
        return head + body + tail
    } else {
        //if h.UnlockPassword != "nil" TODO

        return `<div class="WindowLayout">    
                    <div class="PasswordEntryLayout">
                        <input type="password"
                               value="{{html .UnlockPassword}}"
                               placeholder="Password"
                               autofocus="true"
                               onchange="Unlock"
                               autocomplete="off" 
                               autocorrect="off" 
                               autocapitalize="off" 
                               spellcheck="false"
                               selectable="on" 
                               class="selectable"/>
                    </div>
                    <div class="image"></div>​
                </div>`
    }
}

func FormatListItem(v string) string {
    return `<h3 clickable="on" class="clickable" style="margin-top: 0; margin-left: 15px">
                <a href="PasswordSearch?Account=` + v + `">` + v + `</a>
            </h3>`
}

// Does a search and filter
func GetList(s string) string {
    accounts := DoList(s)
    // Create list with entries
    list := Map(accounts, func(v string) string {
        return FormatListItem(v)
    })
    // Return as a string
    return strings.Join(list, "")
}

// Requests a specific entry
func GetEntry(s string) (string, string) {
    return DoRequest(s)
}


func (h *PasswordSearch) OnHref(URL *url.URL) {
    // Extract information from query
    u := URL.Query()
    // Set account and username/password in internal
    // data structure for showing
    h.Account = u.Get("Account")
    h.Password, h.Username = DoRequest(h.Account)
    // Show it!
    app.Render(h)
    // Clear the response
    h.Account = ""
    h.Password = ""
    h.Username = ""
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
    // Show
    app.Render(h)
}

func (h *PasswordSearch) Search(arg app.ChangeArg) {
  	h.Query = arg.Value
    // Create a list of entries
    h.List = GetList(h.Query)
    // Tells the app to update the rendering of the component.
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

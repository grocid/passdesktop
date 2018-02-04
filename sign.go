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
    "encoding/base64"
    "encoding/json"
    "github.com/murlokswarm/app"
    "golang.org/x/crypto/ed25519"
    "io/ioutil"
    "log"
    "net/url"
    "pass/rest"
)

type Sign struct {
    Title        string
    PublicBase64 string
    Keys         KeyPair
    Data         rest.DecodedEntry
}

type KeyPair struct {
    Priv []byte `json:"priv"`
    Pub  []byte `json:"pub"`
}

var keyPair KeyPair

const (
    SignatureFileExtension = ".signature"
)

func (h *Sign) Render() string {
    var gen string

    if len(h.Data.File) == 0 {
        gen = `<button class="button add" onclick="Generate"/>`
    } else {
        gen = ""
    }

    return `
<div class="WindowLayout">
    <div class="SearchLayout">
        <input type="text"
               value="{{html .Title}}"
               placeholder="Account"
               onchange="DoSearchQuery"
               autocomplete="off"
               autocorrect="off"
               autocapitalize="off"
               spellcheck="false"
               selectable="on"
               class="editable searchfield"/>
        <div class="animated">
            <div style="text-align: center;
                        margin-left: auto;
                        margin-right: auto;
                        margin-top: -webkit-calc(20vh - 20px);">
            ` + GetFingerprint(keyPair.Pub, 255, 90, 60) + `
                <h1>{{.Title}}</h1>

                <h2>Public key</h2>
                <input type="text" 
                     value="{{.PublicBase64}}" 
                     readonly="readonly"
                     class="editable username"/>
            </div>
          </div>

          <div class="bottom-toolbar">
              <div>
                  <button class="button ok" onclick="OK"/>` +
        gen + `<button class="button sign" onclick="Sign"/>
                  <button class="button delete" onclick="Delete"/>
              </div>
          </div>
     </div>
</div>`

}

func (h *Sign) OnHref(URL *url.URL) {
    // Extract information from query and get account name and
    // its encrypted counterpart from query (this is need since
    // if we were to encrypt again, we would get a different
    // encrypted name).
    u := URL.Query()
    h.Title = u.Get("Name")
    restResponse, err := restClient.VaultReadSecret(
        &rest.Name{
            Text:      h.Title,
            Encrypted: u.Get("Encrypted"),
        })

    if err != nil {
        log.Println(err)
        return
    }

    h.Data = *restResponse
    json.Unmarshal(h.Data.File, &keyPair)
    h.PublicBase64 = base64.StdEncoding.EncodeToString(keyPair.Pub)

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *Sign) OK() {
    // Make sure we do not save already saved information.
    // Now, we just need to go back.
    h.Cancel()
}

func (h *Sign) Generate() {
    // By passing nil as argument, we default to
    // crypto/rand, which is desirable.
    pub, priv, err := ed25519.GenerateKey(nil)

    if err != nil {
        log.Println(err)
        return
    }

    keyPair := KeyPair{
        Priv: priv,
        Pub:  pub,
    }

    // We will store the key pair as JSON in the
    // the space where we would store file data.
    jsonKeyPair, err := json.Marshal(&keyPair)
    h.Data.File = jsonKeyPair

    // Write it to remote.
    err = restClient.VaultWriteSecret(&h.Data)

    if err != nil {
        log.Println(err)
        return
    }

    app.Render(h)
}

func (h *Sign) Cancel() {
    keyPair = KeyPair{}
    NavigateBack("")
}

func (h *Sign) Sign() {
    // Open filepicker window to get filename.
    app.NewFilePicker(app.FilePicker{
        MultipleSelection: false,
        NoDir:             true,
        NoFile:            false,
        OnPick: func(filenames []string) {
            // Get contents of file.
            fileData, err := ioutil.ReadFile(filenames[0])

            // If there was an error, probably due
            // to permissions, dump error message
            // to log.
            if err != nil {
                log.Println(err)
                return
            }

            // Generate signature from private key.
            signature := ed25519.Sign(
                ed25519.PrivateKey(keyPair.Priv),
                fileData)

            // Write to signature file, which is
            // the file we are signing + the
            // .signature file extension.
            err = ioutil.WriteFile(
                filenames[0]+SignatureFileExtension, 
                signature, 0644)

            if err != nil {
                log.Println(err)
                return
            }

            app.Render(h)
        },
    })
}

func (h *Sign) DoSearchQuery(arg app.ChangeArg) {
    keyPair = KeyPair{}
    NavigateBack(arg.Value)
}

func (h *Sign) Delete() {
    keyPair = KeyPair{}
    d := h.Data.Name

    if d != nil {
        restClient.VaultDeleteSecret(&h.Data)
    }

    h.Cancel()
}

func init() {
    app.RegisterComponent(&Sign{})
}

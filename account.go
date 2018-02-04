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
    //"pass/dialog"
    "pass/lock"
    "pass/rest"
)

type Account struct {
    //Deletable
    Title     string
    ImageName string
    Query     string
    Data      rest.DecodedEntry
}

func (h *Account) DefaultView() string {
    return `
<div class="WindowLayout">
    <div class="SearchLayout">
        <input type="text"
               value="{{html .Title}}"
               placeholder="Account"
               onchange="DoSearchQuery"
               onsubmit="DoSearchQuery"
               onload="Query"
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
                <img src="iconpack/{{.ImageName}}.png"
                      style="max-width: 128px; "/>
                <p><input name="Name"
                       type="text"
                       value="{{.Title}}"
                       placeholder="Account"
                       onchange="Title"
                       autocomplete="off"
                       autocorrect="off"
                       autocapitalize="off"
                       spellcheck="false"
                       selectable="on"
                       class="editable name"/></p>
            </div>
            <input name="username"
                   type="text"
                   value="{{html .Data.Username}}"
                   placeholder="Username"
                   onchange="Data.Username"
                   autocomplete="off"
                   autocorrect="off"
                   autocapitalize="off"
                   spellcheck="false"
                   selectable="on"
                   class="editable username"/><br/>
            <input name="password"
                   type="text"
                   value="{{html .Data.Password}}"
                   placeholder="Password"
                   onchange="Data.Password"
                   autocomplete="off"
                   autocorrect="off"
                   autocapitalize="off"
                   spellcheck="false"
                   selectable="on"
                   class="editable password"/>
          </div>
          <div class="bottom-toolbar">
              <div>
                  <button class="button ok" onclick="OK"/>
                  <button class="button cancel" onclick="Cancel"/>
                  <button class="button rerand" onclick="RandomizePassword"/>
                  <button class="button delete" onclick="Delete"/>                 
              </div>
          </div>
     </div>
</div>
`
}
func (h *Account) Render() string {
    return h.DefaultView()
}

func (h *Account) OnHref(URL *url.URL) {
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

    // Acquire the image name. If it exists in preloaded map,
    // use it as is, but if it is not, we subsitute.
    h.ImageName = GetImageName(h.Title)

    //log.Println(h.Render())

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *Account) OK() {
    // We do not want empty names.
    if h.Title == "" {
        return
    }

    d := h.Data.Name
    if d != nil && h.Title != (*d).Text {
        // need to remove old and submit new
    } else {
        if d == nil {
            h.Data.Name = &rest.Name{
                Text: h.Title,
            }
        }
        // Modify the decoded entry so that it matches
        //the contents of the UI.
        err := restClient.VaultWriteSecret(&h.Data)

        if err != nil {
            log.Println(err)
            return
        }
    }
    // Now, we just need to go back.
    h.Cancel()
}

func (h *Account) Cancel() {
    NavigateBack("")
}

func (h *Account) RandomizePassword() {
    // Generate a new password.
    h.Data.Password = lock.EntropyAlphabet(32)

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *Account) Delete() {
    d := h.Data.Name
    if d != nil {
        restClient.VaultDeleteSecret(&h.Data)
    }
    h.Cancel()
}

func (h *Account) DoSearchQuery(arg app.ChangeArg) {
    log.Println(arg.Value, h.Query)
    NavigateBack(arg.Value)
}

func init() {
    app.RegisterComponent(&Account{})
}

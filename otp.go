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
    "net/url"
    "pass/otp"
    "pass/rest"
    "log"
)

type OTP struct {
    Title string
    Query string
    OTP   string
    Data  rest.DecodedEntry
}

func (*OTP) Render() string {
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
            <img src="iconpack/otp.png" style="max-width: 128px;"/>
            <h1>{{.Title}}</h1>
            </div>
          <h2>{{.OTP}}</h2>
          </div>
          <div class="bottom-toolbar">
              <div>
                  <button class="button ok" onclick="Cancel"/>
                  <button class="button refresh" onclick="RefreshOTP"/>
                  <button class="button delete" onclick="Delete"/>
              </div>
          </div>
     </div>
</div>`
}

func (h *OTP) OnHref(URL *url.URL) {
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

    // Read contents from data segment.
    h.Data = *restResponse
    h.OTP = otp.ComputeOTPCode(h.Data.Password)

    app.Render(h)

}

func (h *OTP) Cancel() {
    NavigateBack("")
}

func (h *OTP) Delete() {
    d := h.Data.Name
    if d != nil {
        restClient.VaultDeleteSecret(&h.Data)
    }
    h.Cancel()
}

func (h *OTP) RefreshOTP(arg app.ChangeArg) {
    // Compute OTP based on retrieved data segment.
    h.OTP = otp.ComputeOTPCode(h.Data.Password)
    app.Render(h)
}

func (h *OTP) DoSearchQuery(arg app.ChangeArg) {
    NavigateBack(arg.Value)
}

func init() {
    app.RegisterComponent(&OTP{})
}

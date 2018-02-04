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
    "encoding/hex"
    "github.com/murlokswarm/app"
    "log"
    "pass/lock"
    "pass/rest"
    "pass/util"
)

type UnlockScreen struct{}

var restClient rest.Client

func (h *UnlockScreen) OnDismount() {
    log.Println("UnlockScreen dismounted")
}

func (h *UnlockScreen) Render() string {
    return `
<div class="WindowLayout">
    <div class="animated">
        <div class="PasswordEntryLayout">
            <input type="password"
                   placeholder="Password"
                   autofocus="true"
                   onchange="Unlock"
                   autocomplete="off" 
                   autocorrect="off" 
                   autocapitalize="off" 
                   spellcheck="false"
                   selectable="on" 
                   class="editable password"/>
        </div>
        <div class="symbol lock"/><i class="fas fa-font"></i><i class="fas fa-font"></i><i class="fas fa-font"></i><i class="fas fa-font"></i><i class="fas fa-font"></i><i class="fas fa-font"></i>
    </div>
</div>`
}

func (h *UnlockScreen) Unlock(arg app.ChangeArg) {
    // Get the password.
    password := arg.Value
    salt, _ := hex.DecodeString(config.Encrypted.Salt)

    // Verify password against encrypted token + mac
    lock := lock.New(password, salt)
    _, err := lock.UnlockToken(config.Encrypted.Token)

    if err == nil {
        log.Println("Unlocked.")

        pass.Locked = false
        restClient = rest.New(&lock)
        restClient.Init(config.Host, config.Port, config.CA)
        restClient.Unlock(config.Encrypted.Token)

        config = util.Configuration{}

        log.Println("Fetching data.")
        r, _ := restClient.VaultListSecrets()
        log.Println("OK")

        ps := &Search{Result: *r}
        win.Mount(ps)

        return

    } else {
        log.Println(err)
    }

    app.Render(h)
}

func init() {
    // Register UI component
    app.RegisterComponent(&UnlockScreen{})
}

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

import "os"

func GetCreateConfigDialog() string {
    return `<div class="clickable" >
                <div class="animated">
                    <div style="text-align: center; 
                                margin-left: auto; 
                                margin-right: auto;
                                padding-top: 30px">
                        <h1>Welcome</h1>
                    </div>
                    <p>It is the first time you are using Pass Desktop. Please enter credentials.</p>
                        <h2 style="margin-bottom: 15px">Secrets</h2>
                        <input type="text"
                               placeholder="Token"
                               autofocus="true"
                               onchange="Query"
                               autocomplete="off" 
                               autocorrect="off" 
                               autocapitalize="off" 
                               spellcheck="false"
                               selectable="on" 
                               class="selectable"/>
                        <input type="password"
                               placeholder="Password"
                               onchange="Password"
                               autocomplete="off" 
                               autocorrect="off" 
                               autocapitalize="off" 
                               spellcheck="false"
                               selectable="on" 
                               class="selectable"/>
                        <input type="password"
                               placeholder="Password"
                               onchange="UnlockPassword"
                               autocomplete="off" 
                               autocorrect="off" 
                               autocapitalize="off" 
                               spellcheck="false"
                               selectable="on" 
                               class="selectable"/>
                        <h2 style="margin-top: 15px; margin-bottom: 15px">Server</h2>
                        <input type="text"
                               placeholder="myserver.com"
                               onchange="Account"
                               autocomplete="off" 
                               autocorrect="off" 
                               autocapitalize="off" 
                               spellcheck="false"
                               selectable="on" 
                               class="selectable"/>
                        <input type="text"
                               placeholder="8001"
                               onchange="Meta"
                               autocomplete="off" 
                               autocorrect="off" 
                               autocapitalize="off" 
                               spellcheck="false"
                               selectable="on" 
                               class="selectable"/>
                        <div id="field">
                            <input type="text"
                                   value="{{html .Filename}}"
                                   placeholder="Certificate Authority"
                                   autofocus="true"
                                   onchange="Query"
                                   autocomplete="off" 
                                   autocorrect="off" 
                                   autocapitalize="off" 
                                   spellcheck="false"
                                   style="max-width: 220px"/>
                              <button class="button add" 
                                      style="background-size: 14px; 
                                             padding: 10px;" 
                                      onclick="PickFile"/>
                        </div>
                    <div class="bottom-toolbar">
                        <button class="button ok" onclick="CreateConfig"/>
                    </div>
                </div>
            </div>`
}

func GetConfirmDeleteDialog() string {
    return `<div class="WindowLayout">    
                <div class="SearchLayout">
                    <input type="text"
                           value="{{html .Query}}"
                           placeholder="Account"
                           onchange="Search"
                           autocomplete="off" 
                           autocorrect="off" 
                           autocapitalize="off" 
                           spellcheck="false"
                           selectable="on" 
                           class="editable searchfield"/>
                </div>
                <div class="animated">
                    <div class="symbol trash"/>
                    <div class="bottom-toolbar">
                        <button class="button ok" onclick="OkAccountView"/>
                        <button class="button cancel" onclick="CancelTrashView"/>
                    </div>
                </div>
            </div>`
}

func GetEmptySearchDialog() string {
    return `<div class="WindowLayout">    
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
                           class="editable searchfield"/>
                </div>
                <div class="animated">
                    <div class="symbol search"/>
                    <div class="bottom-toolbar">
                        <button class="button add" onclick="OkAccountView"/>
                    </div>
                </div>
            </div>`
}


func GetAddDialog() string {
    return `<div class="WindowLayout">    
                <div class="SearchLayout">
                    <input type="text"
                           value="{{html .Query}}"
                           placeholder="Account"
                           onchange="Search"
                           autocomplete="off" 
                           autocorrect="off" 
                           autocapitalize="off" 
                           spellcheck="false"
                           selectable="on" 
                           class="editable searchfield"/>
                    <div  style="overflow-y:scroll; 
                         max-height:335px;
                         height:335px;"
                         clickable="on" 
                         class="clickable">
                        <div class="animated">
                            <h1>{{.Account}}</h1>
                            <h2>Username</h2>
                            <input name="username"
                                   type="text"
                                   value="{{html .Username}}"
                                   placeholder="Username"
                                   onchange="Username"
                                   autocomplete="off" 
                                   autocorrect="off" 
                                   autocapitalize="off" 
                                   spellcheck="false"
                                   selectable="on" 
                                   class="editable"/>
                            <h2>Password</h2>
                            <input name="password"
                                   type="text"
                                   value="{{html .Password}}"
                                   placeholder="Password"
                                   onchange="Password"
                                   autocomplete="off" 
                                   autocorrect="off" 
                                   autocapitalize="off" 
                                   spellcheck="false"
                                   selectable="on" 
                                   class="editable"/>
                        </div>
                        <div class="bottom-toolbar">
                            <div>
                                <button class="button ok" onlick="OkAccountView"/>
                                <button class="button cancel" onclick="CancelAccountView"/>
                                <button class="button rerand" onclick="RerandomizePasswordAccountView"/>
                                <button class="button delete" onclick="DeleteAccountView"/>
                            </div>
                        </div>
                    </div>
                 </div>
            </div>`
}

// Show account details
func GetAccountBody(h *PasswordSearch) string {
    // Some ugly solution since the fallback on image not found does not work
    var image string
    if _, err := os.Stat("/Users/carl/Projekt/Go/Pass/resources/iconpack/" + h.Account + ".png"); os.IsNotExist(err) {
        image = `<img src="/Users/carl/Projekt/Go/Pass/resources/iconpack/default.png" 
                                 style="max-width: 128px; "/>`
    } else {
        image = `<img src="/Users/carl/Projekt/Go/Pass/resources/iconpack/` + h.Account + `.png" 
                                 style="max-width: 128px; "/>`
    }
    return `<div class="WindowLayout">    
                <div class="SearchLayout">
                    <input type="text"
                           value="{{html .Query}}"
                           placeholder="Account"
                           onchange="Search"
                           autocomplete="off" 
                           autocorrect="off" 
                           autocapitalize="off" 
                           spellcheck="false"
                           selectable="on" 
                           class="editable searchfield"/>
                      <div clickable="on" 
                           class="clickable">
                          <div class="animated">
                            <div style="text-align: center; 
                                        margin-left: auto; 
                                        margin-right: auto;
                                        padding-top: 30px">` + image + `
                                 <h1>{{.Account}}</h1>
                            </div>
                            <h2>Username</h2>
                            <input name="username"
                                   type="text"
                                   value="{{html .Username}}"
                                   placeholder="Username"
                                   onchange="Username"
                                   autocomplete="off" 
                                   autocorrect="off" 
                                   autocapitalize="off" 
                                   spellcheck="false"
                                   selectable="on" 
                                   class="editable"/>
                            <h2>Password</h2>
                            <input name="password"
                                   type="text"
                                   value="{{html .Password}}"
                                   placeholder="Password"
                                   onchange="Password"
                                   autocomplete="off" 
                                   autocorrect="off" 
                                   autocapitalize="off" 
                                   spellcheck="false"
                                   selectable="on" 
                                   class="editable"/>
                        </div>
                        <div class="bottom-toolbar">
                            <div style="">
                                <button class="button ok" onlick="OkAccountView"/>
                                <button class="button cancel" onclick="CancelAccountView"/>
                                <button class="button rerand" onclick="RerandomizePasswordAccountView"/>
                                <button class="button delete" onclick="DeleteAccountView"/>
                            </div>
                        </div>
                    </div>
                 </div>
            </div>`
}

// List view
func GetListBody(h *PasswordSearch) string {
    var siteList string
    for _, element := range h.List {
        if _, err := os.Stat("/Users/carl/Projekt/Go/Pass/resources/iconpack/" + element + ".png"); os.IsNotExist(err) {
            siteList = siteList + `<a href="PasswordSearch?Account=` + element + `">
                                        <li>
                                            <img src="/Users/carl/Projekt/Go/Pass/resources/iconpack/default.png"/>
                                            <h3>` + element + `</h3>
                                        </li>
                                    </a>`
        } else {
            siteList = siteList +`<a href="PasswordSearch?Account=` + element + `">
                                        <li>
                                            <img src="/Users/carl/Projekt/Go/Pass/resources/iconpack/` + element + `.png"/>
                                            <h3>` + element + `</h3>
                                        </li>
                                    </a>`
        }
    }


  // <img src="https://{{$element}}/favicon.ico"/>
    return `<div class="WindowLayout">    
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
                           class="editable searchfield"/>
                      <div  style="overflow-y:scroll; 
                                   max-height:450px; 
                                   margin-top:20px" 
                            clickable="on" 
                            class="clickable">
                            <div class="animated">
                                <ul>` + siteList + `
                                </ul>
                          </div>
                      </div>
                  </div>
              </div>`
}


// Password-input dialog
func GetPasswordInput() string {
    return `<div class="WindowLayout">
                <div class="animated">
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
                               class="editable password"/>
                    </div>
                    <div class="symbol lock"/>
                </div>
            </div>`
}

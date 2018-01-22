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

// This is the search-window header
func GetSearchInput() string {
    return  `<div class="WindowLayout">    
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
}

func GetConfirmDeleteDialog() string {
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
                           class="selectable"/>
                </div>
                <div class="animated">
                    <div class="symbol trash"/>
                    <div class="bottom-toolbar">
                        <button class="button ok" onlick="OkAccountView"/>
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
                           class="selectable"/>
                </div>
                <div class="animated">
                    <div class="symbol search"/>
                    <div class="bottom-toolbar">
                        <button class="button add" onlick="OkAccountView"/>
                    </div>
                </div>
            </div>`
}


func GetAddDialog() string {
    return `<div  style="overflow-y:scroll; 
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

// Show account details
func GetAccountBody() string {
    return `<div  style="overflow-y:scroll; 
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
func GetListBody() string {
    return `<div  style="overflow-y:scroll; 
                         max-height:375px; 
                         margin-top:20px" 
                  clickable="on" 
                  class="clickable">
                  <div class="animated">
                {{ range $index, $element := .List }}
                    <div>
                        <h3 clickable="on" 
                            class="clickable">
                            <a href="PasswordSearch?Account={{$element}}">{{$element}}</a>
                        </h3>
                    </div>
                {{ end }}
                </div>
            </div>
        </div>
    </div>`
}



func GetTail() string {
    return `</div></div>`
}

// Password-input dialog
func GetPasswordInput() string {
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
                <div class="symbol lock"/>
            </div>`
}

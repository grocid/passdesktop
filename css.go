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
    "os"
    "fmt"
)

// TODO: fix this to relative path
const ImagePathSuffix = "/../Resources/iconpack/"

func GetCreateConfigDialog() string {
    return `
<div class="WindowLayout" >
    <div class="animated">
        <div style="padding-top: 30px">
            <h1>Welcome</h1>
        </div>
        <p>
            “Well, all information looks like noise until you break the code.”
        </p>
        <p>
            ― Neal Stephenson, Snow Crash
        </p>
        <p>
            <input type="text"
                   placeholder="Token"
                   autofocus="true"
                   onchange="Token"
                   autocomplete="off" 
                   autocorrect="off" 
                   autocapitalize="off" 
                   spellcheck="false"
                   selectable="on" 
                   class="editable username"/>
        </p><p>
            <input type="password"
                   placeholder="Password"
                   onchange="Password"
                   autocomplete="off" 
                   autocorrect="off" 
                   autocapitalize="off" 
                   spellcheck="false"
                   selectable="on" 
                   class="editable password"/>
            <input type="password"
                   placeholder="Password"
                   onchange="PasswordAgain"
                   autocomplete="off" 
                   autocorrect="off" 
                   autocapitalize="off" 
                   spellcheck="false"
                   selectable="on" 
                   class="editable password"/>
        </p><p>
            <input type="text"
                   placeholder="myserver.com"
                   onchange="Hostname"
                   autocomplete="off" 
                   autocorrect="off" 
                   autocapitalize="off" 
                   spellcheck="false"
                   selectable="on" 
                   class="editable server"/>
            <input type="text"
                   placeholder="8001"
                   onchange="Port"
                   autocomplete="off" 
                   autocorrect="off" 
                   autocapitalize="off" 
                   spellcheck="false"
                   selectable="on" 
                   class="editable port"/>
        </p>
        <div class="bottom-toolbar">
            <button class="button ok" onclick="PickFileAndCreateConfig"/>
        </div>
    </div>
</div>`
}

// Password-input dialog
func GetPasswordInput() string {
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
        <div class="symbol lock"/>
    </div>
</div>`
}

func GetConfirmDeleteDialog() string {
    return `
<div class="WindowLayout">    
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
            <button class="button ok" onclick="OkTrashView"/>
            <button class="button cancel" onclick="CancelTrashView"/>
        </div>
    </div>
</div>`
}

func GetEmptySearchDialog() string {
    return `
<div class="WindowLayout">    
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
            <button class="button add" onclick="CreateAccountView"/>
        </div>
    </div>
</div>`
}


func GetAddDialog() string {
    return `
<div class="WindowLayout">    
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
             class="scrollable">
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
                    <button class="button ok" onclick="OkAccountView"/>
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
func GetAccountBody(account string) string {
    imagePath := pass.FullPath + ImagePathSuffix

    // Some ugly solution since the fallback on image not found does not work...
    image := account

    if _, err := os.Stat(imagePath + account + ".png"); os.IsNotExist(err) {
        image = "default"
    }

    return `
<div class="WindowLayout">    
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
                            padding-top: 30px">
                     <img src="` + imagePath + image + `.png" 
                          style="max-width: 128px; "/>
                     <h1>{{.Account}}</h1>
                </div>
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
                       class="editable username"/><br/>
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
                       class="editable password"/>
            </div>
            <div class="bottom-toolbar">
                <div>
                    <button class="button ok" onclick="OkAccountView"/>
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
func GetListBody(searchResults [] string) string {
    var accountListFormatted string
    
    imagePath := pass.FullPath + ImagePathSuffix

    // Iterate through the search results.
    for _, element := range searchResults {
        image := element

        // Revert to default icon if account icon does not exist.
        if _, err := os.Stat(imagePath + element + ".png"); os.IsNotExist(err) {
            image = "default"
        }

        // Format listitem.
        item := fmt.Sprintf(`<a href="PassView?Account=%s">
                                <li>
                                    <img src="%s%s.png"/>
                                    <div class="SearchListItemCaption"><span>%s</span></div>
                                </li>
                             </a>`, element, imagePath, image, element)

        // Concatenate list.
        accountListFormatted = accountListFormatted + item
    }


  // <img src="https://{{$element}}/favicon.ico"/>
    return `
<div class="WindowLayout">    
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
         <div clickable="on" 
              class="scrollable">
            <div class="animated">
                <ul>` + accountListFormatted + `
                </ul>
            </div>
        </div>
    </div>
</div>`
}

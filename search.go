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
    "fmt"
    "github.com/murlokswarm/app"
    "pass/rest"
    "strings"
)

type Search struct {
    Query  string
    Result []rest.Name
}

func (h *Search) Render() string {
    // Ouput
    filteredNameList := `
<div class="WindowLayout">
    <div class="SearchLayout">
        <input type="text"
               value="{{html .Query}}"
               placeholder="Account"
               onchange="DoSearchQuery"
               autofocus="true"
               autocomplete="off"
               autocorrect="off"
               autocapitalize="off"
               spellcheck="false"
               selectable="on"
               class="editable searchfield"/>
         <div clickable="on"
              class="scrollable">
            <div class="animated">
                <ul>`

    // Since we need to concatenate the results to a string, it is
    // cheapest (both in terms of memory and computations) to perform
    // filtering at this stage, rather than earlier filtering of the list.
    for _, name := range h.Result {
        // Due to optimization, we have encoded some data in the name.
        // We extract this data.
        title, label, image := rest.DecodeName(name.Text)

        // Get the icon if it exists, otherwise, substitue.
        filename := GetImageName(image)

        // Match the query against the current item to decide if we
        // should display it or not.
        if h.Query == "" || strings.Contains(strings.ToLower(title), strings.ToLower(h.Query)) {
            // Append to output.
            filteredNameList = filteredNameList + fmt.Sprintf(`
                    <a href="%s?Name=%s;Encrypted=%s">
                        <li>
                            <img src="iconpack/%s.png"/>
                            <div class="SearchListItemCaption">%s</div>
                            <div class="SearchListItemLabel">%s</div>
                        </li>
                    </a>`,
                label, title, name.Encrypted, filename, title, label)
        }
    }

    // Just to make the code look a bit cleaner.
    filteredNameList = filteredNameList + `
                </ul>
            </div>
        </div>
    </div>
</div>`

    return filteredNameList
}

func (h *Search) Prefetch(query string) {
    // Fetch from Vault.
    r, _ := restClient.VaultListSecrets()
    h.Result = *r

    // Update query field.
    h.Query = query
}

func (h *Search) DoSearchQuery(arg app.ChangeArg) {
    // ...and fetch accounts based on query.
    h.Prefetch(arg.Value)
    app.Render(h)
}

func NavigateBack(query string) {
    // This is used by other windows, to get less code repetition.
    s := Search{}
    s.Prefetch(query)
    win.Mount(&s)
}

func init() {
    app.RegisterComponent(&Search{})
}

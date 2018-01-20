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
    "net/http"
    "io/ioutil"
    "log"
    "encoding/json"
    "strings"
)

type VaultResponseGet struct {
    Data struct {
        Password  string `json:"password"`
        Username string  `json:"username"`
    } `json:"data"`
}

type VaultResponseList struct {
    Data struct {
        Keys  []string `json:"keys"`
    } `json:"data"`
}

func DoRequest(s string) (string, string) {
    // Do a GET for the specified entry...
    req, err := http.NewRequest("GET", entryPoint + "/" + s, nil)
    req.Header.Add("X-Vault-Token", decryptedToken)
    resp, err := client.Do(req)
    // This should not happen, unless entry was deleted in the meantime...
    // TODO: we should handle this more gracefully...
    if err != nil {
      log.Fatal(err)
    }
    // Read the body...
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    // ...and parse the JSON
    r := VaultResponseGet{}
    json.Unmarshal([]byte(body), &r)
    // and give back to UI
    return r.Data.Password, r.Data.Username 
}

func DoList(s string) []string {
    // Do a LIST to get all entries...
    req, err := http.NewRequest("LIST", entryPoint, nil)
    req.Header.Add("X-Vault-Token", decryptedToken)
    resp, err := client.Do(req)
    // Again, this should not happen...
    if err != nil {
      log.Fatal(err)
    }
    // Read in the data...
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    // and parse JSON
    r := VaultResponseList{}
    json.Unmarshal([]byte(body), &r)
    // Create a variable with the accounts...
    var accounts [] string
    // Here is the filtering part...
    if s != "" {
        // Filter the entries
        accounts = Filter(r.Data.Keys, func(v string) bool {
            return strings.Contains(v, s)
        })
    } else {
        // If filter was empty, we treat it as *
        accounts = r.Data.Keys
    }
    // Return to UI
    return accounts
}

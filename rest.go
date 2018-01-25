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
    "bytes"
)

type (
    UserData struct {
        Password string `json:"password"`
        Username string  `json:"username"`
    }

    VaultStruct struct {
        Data UserData `json:"data"`
    }

    VaultResponseList struct {
        Data struct {
            Keys  []string `json:"keys"`
        } `json:"data"`
    }
)

const VaultTokenHeader = "X-Vault-Token"

func DoRequest(operation string, s string) (*http.Response, error) {
    // Create the request based on operation input.
    req, err := http.NewRequest(operation, pass.EntryPoint + s, nil)

    if err != nil {
        return nil, err
    }

    // Add header and do a GET for the specified entry...
    req.Header.Add(VaultTokenHeader, pass.DecryptedToken)
    resp, err := pass.Client.Do(req)

    // This should not happen, unless entry was deleted in the meantime...
    if err != nil {
        return nil, err
    }

    return resp, nil
}

func DoGetRequest(s string) AccountInfo {
    // Retrieve data for a specific account.
    resp, err := DoRequest(http.MethodGet, "/" + s)

    if err != nil {
        log.Fatal(err)
    }

    // Read the body...
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        log.Fatal(err)
    }

    // ...and parse the JSON
    r := VaultStruct{}
    json.Unmarshal([]byte(body), &r)
    
    // ...generate a AccountInfo struct...
    account := AccountInfo{}
    
    // ...with the proper information...
    account.Name = s
    account.Username = r.Data.Username
    account.Password = r.Data.Password
    
    // ...and return to caller.
    return account
}


func DoPutRequest(data AccountInfo) error {
    // Create payload
    payload := &UserData {
        Username: data.Username,
        Password: data.Password,
    }

    // Encode data as JSON.
    jsonPayload, err := json.Marshal(payload)
    encodedPayload := bytes.NewBuffer(jsonPayload)

    // Create the actual request.
    req, err := http.NewRequest(http.MethodPut, 
                                pass.EntryPoint + "/" + data.Name, 
                                encodedPayload)
    req.Header.Add(VaultTokenHeader, pass.DecryptedToken)

    if err != nil {
        return err
    }

    // Do a PUT with the associated data.
    _, err = pass.Client.Do(req)

    if err != nil {
        return err
    }

    return nil
}

func DoDeleteRequest(data AccountInfo) error {
    _, err := DoRequest(http.MethodDelete, "/" + data.Name)

    return err
}

func DoListRequest(s string) []string {
    // Do a LIST to get all entries.
    resp, err := DoRequest("LIST", "")

    if err != nil {
        log.Fatal(err)
    }

    // Read in the data...
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        log.Fatal(err)
    }

    // ...and parse JSON.
    r := VaultResponseList{}
    json.Unmarshal([]byte(body), &r)

    // Create a variable with the accounts.
    var accounts [] string

    // Filtering of data.
    if s != "" {
        // Filter the entries.
        accounts = Filter(r.Data.Keys, func(v string) bool {
            return strings.Contains(v, s)
        })
    } else {
        // If filter was empty, we treat it as wildcard.
        accounts = r.Data.Keys
    }
    
    // Return to UI.
    return accounts
}

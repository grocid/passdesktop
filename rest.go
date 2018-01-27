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


func Enc(data string) (string, error) {
    encData, err := EncryptAndEncode(data, pass.EncryptionKey)
    return encData, err
}

func Dec(data string) (string, error) {
    encData, err := DecodeAndDecrypt(data, pass.EncryptionKey)
    return encData, err
}

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

func DoGetRequest(data Entry) AccountInfo {
    // Retrieve data for a specific account.
    resp, err := DoRequest(http.MethodGet, "/" + data.Encrypted)

    if err != nil {
        log.Println("1")
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
    account.Name = data.Name
    account.Encrypted = data.Encrypted
    account.Username, err = Dec(r.Data.Username)
    account.Password, err = Dec(r.Data.Password)
    
    // ...and return to caller.
    return account
}


func DoPutRequest(data AccountInfo) error {
    // Create payload
    encUsername, err := Enc(data.Username)
    encPassword, err := Enc(data.Password)
    payload := &UserData {
        Username: encUsername,
        Password: encPassword,
    }
    log.Println(data)


    // Encode data as JSON.
    jsonPayload, err := json.Marshal(payload)
    encodedPayload := bytes.NewBuffer(jsonPayload)

    if data.Encrypted == "" {
        data.Encrypted, err = Enc(data.Name)
    }

    if err != nil {
        log.Println("1")
        return err
    }
    log.Println(data)

    // Create the actual request.
    req, err := http.NewRequest(http.MethodPut, 
                                pass.EntryPoint + "/" + data.Encrypted,
                                encodedPayload)
    req.Header.Add(VaultTokenHeader, pass.DecryptedToken)

    if err != nil {
        log.Println("1")
        return err
    }

    // Do a PUT with the associated data.
    _, err = pass.Client.Do(req)

    if err != nil {
        log.Println("1")
        return err
    }

    return nil
}

func DoDeleteRequest(data AccountInfo) error {
    if data.Encrypted == "" {
        log.Fatal("No encrypted data stored")
    }

    _, err := DoRequest(http.MethodDelete, "/" + data.Encrypted)

    return err
}

func DoListRequest(s string) []Entry {
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

    // Create a variable with the accounts, this time
    // decrypted so we can filter them based on our
    // search query.
    accounts := Map(r.Data.Keys, func(v string) Entry {
                        decrypted, err := Dec(v)
                        if err != nil {
                            log.Println(v, err)
                            return Entry{}
                        }
                        return Entry{decrypted, v}
                    })

    // Filter out erronous entries, which may have failed
    // due to format or message-authentication error.
    accounts = Filter(accounts, func(v Entry) bool {
        return v.Name != ""

    })

    // Filtering of data.
    if s != "" {
        // Filter the entries.
        accounts = Filter(accounts, func(v Entry) bool {
            return strings.Contains(v.Name, s)

        })
    }

    // Return to UI.
    return accounts
}

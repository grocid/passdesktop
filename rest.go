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
    "bytes"
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
)

type (
    UserData struct {
        Password string `json:"password"`
        Username string `json:"username"`
        File     []byte `json:"file"`
        Padding  string `json:"padding"`
    }

    MyRequestTag struct {
        Tag string `json:"tag"`
    }

    MyRequestEncrypted struct {
        Encrypted []byte `json:"encrypted"`
    }

    MyResponse struct {
        Data struct {
            Tag       string   `json:"tag"`
            Encrypted []byte   `json:"encrypted"`
            Keys      []string `json:"keys"`
        } `json:"data"`
        Errors []string `json:"errors"`
    }
)

const (
    VaultTokenHeader  = "X-Vault-Token"
    TagPath           = "/updated"
    MinimumDataLength = 3 * 32
)

func EncHex(data string) (string, error) {
    encData, err := EncryptAndEncodeHex(data, pass.EncryptionKey)
    return encData, err
}

func DecHex(data string) (string, error) {
    encData, err := HexDecodeAndDecrypt(data, pass.EncryptionKey)
    return encData, err
}

func EncBase64(data string) ([]byte, error) {
    encData, err := EncryptAndEncodeBase64(data, pass.EncryptionKey)
    return encData, err
}

func DecBase64(data []byte) (string, error) {
    encData, err := Base64DecodeAndDecrypt(data, pass.EncryptionKey)
    return encData, err
}

func UpdateTag() error {
    // Let the client know that we did an tag update and therefore
    // do not need to check it.
    pass.LocalUpdate = true

    // Update the tag to indicate (for other clients) that something
    // has changed, i.e., we have done a PUT or a DELETE. To do so,
    // we generate random string, which w.h.p does not collide with
    // the existing one.
    tag, err := Entropy(32)
    storedTag := MyRequestTag{
        Tag: tag,
    }

    // Convert to the payload to JSON.
    jsonStoredTag, err := json.Marshal(storedTag)

    // Put the new tag in place by doing a PUT on the tag path.
    _, err = Request(http.MethodPut, TagPath,
        bytes.NewBuffer([]byte(jsonStoredTag)))

    return err
}

func IsTagUpdated() bool {
    // If we did a PUT or DELETE from this client, we already know
    // it must be updated.
    if pass.LocalUpdate {
        return true
    }

    // Obtain the value of the tag.
    vaultResponse, err := Request(http.MethodGet, TagPath, nil)

    if err != nil {
        return false
    }

    // Check whether old tag and obtained tag match or not.
    tagUpdated := pass.CachedTag != vaultResponse.Data.Tag

    // Next time we call it, it will not differ if there was no
    // additional change between the calls.
    pass.CachedTag = vaultResponse.Data.Tag

    return tagUpdated
}

func Request(operation string, s string, data *bytes.Buffer) (MyResponse, error) {
    var err error
    var req *http.Request

    log.Println(operation, pass.EntryPoint+s)

    // These two cases need to be handled separarely, i.e., the buffer must
    // explicitly be set to nil, we cannot pass a pointer with nil, or
    // program will throw a SIGSEGV.
    if data != nil {
        req, err = http.NewRequest(operation, pass.EntryPoint+s, data)
    } else {
        req, err = http.NewRequest(operation, pass.EntryPoint+s, nil)
    }

    if err != nil {
        return MyResponse{}, err
    }

    // Add header and do a GET for the specified entry...
    req.Header.Add(VaultTokenHeader, pass.DecryptedToken)
    resp, err := pass.Client.Do(req)

    // This should not happen, unless entry was deleted in the meantime...
    if err != nil {
        return MyResponse{}, err
    }

    // Read the body...
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    response := MyResponse{}
    json.Unmarshal([]byte(body), &response)

    return response, nil
}

func VaultReadSecret(data Entry) AccountInfo {
    // Retrieve data for a specific account.
    vaultResponse, err := Request(http.MethodGet, "/"+data.Encrypted, nil)

    if err != nil {
        log.Fatal(err)
    }

    // Decrypt
    decryptedData, err := DecBase64(vaultResponse.Data.Encrypted)

    // ...generate a AccountInfo struct...
    account := AccountInfo{}
    json.Unmarshal([]byte(decryptedData), &account)

    // ...with the proper information...
    account.Name = data.Name
    account.Encrypted = data.Encrypted

    // ...and return to caller.
    return account
}

func VaultWriteSecret(data AccountInfo) error {
    var padding string

    log.Println("WRITE", data)

    // Get some padding data, so that, ciphertext length
    // does not reveal information about password lenth.
    contentLength := len(data.Username) + len(data.Username) + len(data.File)

    if contentLength < MinimumDataLength {
        padding, _ = Entropy(MinimumDataLength - contentLength)
    }

    userData := &UserData{
        Username: data.Username,
        Password: data.Password,
        File:     data.File,
        Padding:  padding,
    }

    // Encode data as JSON...
    jsonUserData, err := json.Marshal(userData)

    // ...and encrypt.
    encryptedUserData, err := EncBase64(string(jsonUserData))

    if data.Encrypted == "" {
        data.Encrypted, err = EncHex(data.Name)
    }

    if err != nil {
        return err
    }

    vaultRequestEncrypted := MyRequestEncrypted{
        Encrypted: encryptedUserData,
    }
    jsonVaultRequestEncrypted, err := json.Marshal(vaultRequestEncrypted)

    // Create the actual request.
    _, err = Request(http.MethodPut, "/"+data.Encrypted,
        bytes.NewBuffer(jsonVaultRequestEncrypted))

    if err != nil {
        return err
    }

    log.Println("OK")

    return UpdateTag()
}

func VaultDeleteSecret(data AccountInfo) error {
    if data.Encrypted == "" {
        log.Fatal("No encrypted data stored")
    }

    _, err := Request(http.MethodDelete, "/"+data.Encrypted, nil)

    if err != nil {
        return err
    }

    if err != nil {
        return err
    }

    return UpdateTag()
}

func VaultListSecrets(matchingString string) []Entry {
    if IsTagUpdated() {
        // Do a LIST to get all entries.
        vaultResponse, err := Request("LIST", "", nil)

        if err != nil {
            return []Entry{}
        }

        pass.SearchResult = Map(vaultResponse.Data.Keys,
            func(v string) Entry {
                decrypted, err := DecHex(v)
                if err != nil {
                    return Entry{}
                }
                return Entry{decrypted, v}
            })

        // Filter out erronous entries, which may have failed
        // due to format or message-authentication error.
        pass.SearchResult = Filter(pass.SearchResult,
            func(v Entry) bool {
                return v.Name != ""
                //return v != Entry{}
            })
        // sort it!

    } else {
        log.Println("Not tag change: using cached results")
    }

    // If the empty matchingString should match all results, while
    // the non-empty must be substring of the results.
    if matchingString != "" {
        // Filter the entries.
        filteredSearchResult := Filter(pass.SearchResult,
            func(v Entry) bool {
                return strings.Contains(v.Name, matchingString)
            })
        return filteredSearchResult
    } else {
        // Return non-filtered result.
        return pass.SearchResult
    }
}

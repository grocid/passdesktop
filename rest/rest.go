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

package rest

import (
    "bytes"
    "crypto/tls"
    "crypto/x509"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "pass/lock"
    "time"
    "errors"
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

    Name struct {
        Text      string
        Encrypted string
    }

    DecodedEntry struct {
        Name     *Name
        Username string
        Password string
        File     []byte
    }
)

const (
    VaultTokenHeader  = "X-Vault-Token"
    TagPath           = "/updated"
    MinimumDataLength = 3 * 32

    MethodList = "LIST"
)

type Client struct {
    LocalUpdate  bool
    CachedTag    string
    SearchResult []Name

    DecryptedToken string
    EncryptionKey  []byte
    Lock           *lock.Lock

    Client     *http.Client
    EntryPoint string
}

func New(lock *lock.Lock) Client {

    r := Client{
        LocalUpdate:    true,
        Client:         nil,
        Lock:           lock,
        DecryptedToken: "",
    }

    return r
}

func (r *Client) Unlock(encryptedToken string) error {
    t, err := r.Lock.UnlockToken(encryptedToken)

    if err != nil {
        return err
    }

    r.DecryptedToken = t
    log.Println(t)

    return nil
}

func (r *Client) Init(hostname string, port int, CA string) {
    // Setup entrypoint
    r.EntryPoint = fmt.Sprintf("https://%s:%v/v1/secret", hostname, port)
    r.CachedTag = "-"

    // Create a TLS context...
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM([]byte(CA))

    // ...and a client
    r.Client = &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs: caCertPool,
            },
        },
        Timeout: time.Second * 10,
    }
}

func (r *Client) EncHex(data string) (string, error) {
    encData, err := r.Lock.EncryptAndEncodeHex(data)
    return encData, err
}

func (r *Client) DecHex(data string) (string, error) {
    encData, err := r.Lock.HexDecodeAndDecrypt(data)
    return encData, err
}

func (r *Client) EncBase64(data string) ([]byte, error) {
    encData, err := r.Lock.EncryptAndEncodeBase64(data)
    return encData, err
}

func (r *Client) DecBase64(data []byte) (string, error) {
    encData, err := r.Lock.Base64DecodeAndDecrypt(data)
    return encData, err
}

func (r *Client) Request(operation string, s string, data *bytes.Buffer) (MyResponse, error) {
    var err error
    var req *http.Request

    // These two cases need to be handled separarely, i.e., the buffer must
    // explicitly be set to nil, we cannot pass a pointer with nil, or
    // program will throw a SIGSEGV.
    if data != nil {
        req, err = http.NewRequest(operation, r.EntryPoint+s, data)
    } else {
        req, err = http.NewRequest(operation, r.EntryPoint+s, nil)
    }

    if err != nil {
        return MyResponse{}, err
    }

    // Add header and do a GET for the specified entry...
    req.Header.Add(VaultTokenHeader, r.DecryptedToken)
    resp, err := r.Client.Do(req)

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

func (r *Client) UpdateTag() error {
    // Let the client know that we did an tag update and therefore
    // do not need to check it.
    r.LocalUpdate = true

    // Update the tag to indicate (for other clients) that something
    // has changed, i.e., we have done a PUT or a DELETE. To do so,
    // we generate random string, which w.h.p does not collide with
    // the existing one.
    storedTag := MyRequestTag{
        Tag: string(lock.Entropy(32)),
    }

    // Convert to the payload to JSON.
    jsonStoredTag, err := json.Marshal(storedTag)

    // Put the new tag in place by doing a PUT on the tag path.
    _, err = r.Request(http.MethodPut, TagPath,
        bytes.NewBuffer([]byte(jsonStoredTag)))

    return err
}

func (r *Client) IsTagUpdated() bool {
    // If we did a PUT or DELETE from this client, we already know
    // it must be updated.
    if r.LocalUpdate {
        return true
    }

    // Obtain the value of the tag.
    vaultResponse, err := r.Request(http.MethodGet, TagPath, nil)

    if err != nil {
        return false
    }

    // Check whether old tag and obtained tag match or not.
    tagUpdated := r.CachedTag != vaultResponse.Data.Tag

    // Next time we call it, it will not differ if there was no
    // additional change between the calls.
    r.CachedTag = vaultResponse.Data.Tag

    return tagUpdated
}

func (r *Client) VaultReadSecret(data *Name) (*DecodedEntry, error) {
    log.Println("READ")

    // Retrieve data for a specific account.
    vaultResponse, err := r.Request(http.MethodGet, "/"+(*data).Encrypted, nil)

    if err != nil {
        return nil, err
    }

    // Decrypt
    decryptedData, err := r.DecBase64(vaultResponse.Data.Encrypted)

    // ...generate a DecodedEntry struct...
    decodedEntry := DecodedEntry{}
    err = json.Unmarshal([]byte(decryptedData), &decodedEntry)

    if err != nil {
        return nil, err
    }

    // ...with the proper information...
    decodedEntry.Name = data

    // ...and return to caller.
    return &decodedEntry, nil
}

func (r *Client) VaultWriteSecret(data *DecodedEntry) error {
    var padding string

    log.Println("WRITE", data)

    // Get some padding data, so that, ciphertext length
    // does not reveal information about password lenth.
    contentLength := len(data.Username) + len(data.Username) + len(data.File)

    if contentLength < MinimumDataLength {
        padding = string(lock.Entropy(MinimumDataLength - contentLength))
    }

    userData := &UserData{
        Username: (*data).Username,
        Password: (*data).Password,
        File:     (*data).File,
        Padding:  padding,
    }

    // Encode data as JSON...
    jsonUserData, _ := json.Marshal(userData)

    // ...and encrypt.
    encryptedUserData, _ := r.EncBase64(string(jsonUserData))

    if (*data).Name.Encrypted == "" {
        (*data).Name.Encrypted, _ = r.EncHex((*data).Name.Text)
    }

    vaultRequestEncrypted := MyRequestEncrypted{
        Encrypted: encryptedUserData,
    }
    jsonVaultRequestEncrypted, _ := json.Marshal(vaultRequestEncrypted)

    // Create the actual request.
    _, err := r.Request(http.MethodPut, "/"+(*data).Name.Encrypted,
        bytes.NewBuffer(jsonVaultRequestEncrypted))

    if err != nil {
        return err
    }

    log.Println("OK")

    return r.UpdateTag()
}

func (r *Client) VaultDeleteSecret(data *DecodedEntry) error {
    if (*data).Name.Encrypted == "" {
        log.Fatal("No encrypted data stored")
    }

    _, err := r.Request(http.MethodDelete, "/"+(*data).Name.Encrypted, nil)

    if err != nil {
        return err
    }

    return r.UpdateTag()
}

//VaultRenameSecret

func (r *Client) VaultListSecrets() (*[]Name, error) {
    log.Println("LIST")

    if r.IsTagUpdated() {
        // Do a LIST to get all entries.
        vaultResponse, err := r.Request(MethodList, "", nil)

        if err != nil {
            return nil, err
        }

        if len(vaultResponse.Errors) > 0 {
            return nil, errors.New(vaultResponse.Errors[0])
        }

        r.SearchResult = make([]Name, 0)

        for _, key := range vaultResponse.Data.Keys {
            decrypted, err := r.DecHex(key)

            if err == nil {
                r.SearchResult = append(r.SearchResult,
                    Name{
                        Text:      decrypted,
                        Encrypted: key,
                    })
            }
        }

    } else {
        log.Println("No tag change: using cached results")
    }
    return &r.SearchResult, nil
}

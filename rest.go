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

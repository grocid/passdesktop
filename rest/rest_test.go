package rest

import (
    "encoding/hex"
    "fmt"
    "pass/lock"
    "testing"
)

const (
    token    = ""
    salt     = ""
    password = ""
    CA       = ""
    server   = ""
    port     = 8200
)

func TestListVaultSecrets(t *testing.T) {
    s, _ := hex.DecodeString(salt)
    mylock := lock.New(password, s)
    r := New(&mylock)
    r.Unlock(token)
    r.Init(server, port, CA)
    m, _ := r.VaultListSecrets()
    fmt.Println(*m)
}

func TestListAndGetFirstSecret(t *testing.T) {
    s, _ := hex.DecodeString(salt)
    mylock := lock.New(password, s)
    r := New(&mylock)
    r.Unlock(token)
    r.Init(server, port, CA)
    m, _ := r.VaultListSecrets()
    if r.IsTagUpdated() {
        t.Errorf("Tag update not properly working...")
    }
    fmt.Println(r.VaultReadSecret(&(*m)[0]))
}

func TestListAndGetFirstSecretAndWriteItAgain(t *testing.T) {
    s, _ := hex.DecodeString(salt)
    mylock := lock.New(password, s)
    r := New(&mylock)
    r.Unlock(token)
    r.Init(server, port, CA)
    m, _ := r.VaultListSecrets()
    f := r.VaultReadSecret(&(*m)[0])
    (*f).Username = "zzz"
    r.VaultWriteSecret(f)
    f = r.VaultReadSecret(&(*m)[0])
    if (*f).Username != "zzz" {
        t.Errorf("Write error")
    }
    (*f).Username = "mmm"
    r.VaultWriteSecret(f)
}

func CreateEmptySecretAndDeleteIt(t *testing.T) {
    s, _ := hex.DecodeString(salt)
    mylock := lock.New(password, s)
    r := New(&mylock)
    r.Unlock(token)
    r.Init(server, port, CA)

}

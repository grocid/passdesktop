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

package lock

import (
    "crypto/rand"
    "golang.org/x/crypto/chacha20poly1305"
    "testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
    key := make([]byte, 32)
    rand.Read(key)

    plaintext := []byte("this is a test vector")
    ciphertext, err := Chacha20Poly1305Encrypt(plaintext, key)

    if err != nil {
        t.Errorf("Encryption error")
    }

    decryptedCiphertext, err := Chacha20Poly1305Decrypt(ciphertext, key)

    if err != nil {
        t.Errorf("Decryption error")
    }

    if string(decryptedCiphertext) != string(plaintext) {
        t.Errorf("Decryption was incorrect, got: %s, want: %s.",
            plaintext,
            decryptedCiphertext)
    }
}

func TestRepeatedNonce(t *testing.T) {
    key := make([]byte, 32)
    rand.Read(key)

    plaintext := []byte("this is a test vector")

    nonces := make(map[string]bool)

    for i := 0; i < 10000; i++ {
        ciphertext, _ := Chacha20Poly1305Encrypt(plaintext, key)

        nonce := string(ciphertext[:chacha20poly1305.NonceSize])

        if nonces[nonce] == true {
            t.Errorf("Colliding nonces, should be nearly impossible.")
        }

        nonces[nonce] = true
    }
}

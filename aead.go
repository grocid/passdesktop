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
    "io"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
     "golang.org/x/crypto/chacha20poly1305"
  )

func Chacha20Poly1305Encrypt(plaintext []byte, key []byte) ([]byte, error) {
    // Create a key data structure
    chacha20aead, err := chacha20poly1305.New(key)

    if err != nil {
        return []byte{}, err
    }

    nonce := make([]byte, chacha20poly1305.NonceSize)

    // ...using a CSPRNG...
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return []byte{}, err
    }

    ciphertext := chacha20aead.Seal(nil, nonce, plaintext, nil)

    return append(nonce, ciphertext...), nil
}

func Chacha20Poly1305Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
    // Create a key data structure
    chacha20aead, err := chacha20poly1305.New(key)

    if err != nil || len(ciphertext) < chacha20poly1305.NonceSize {
        return []byte{}, err
    }

    plaintext, err := chacha20aead.Open(nil, ciphertext[:12], ciphertext[12:], nil)

    if err != nil {
        return []byte{}, err
    }

    return plaintext, nil
}

func EncryptGCM(plaintext []byte, key []byte) ([]byte, []byte, error) {
    // Create a key data structure
    block, err := aes.NewCipher(key)

    if err != nil {
        return []byte{}, []byte{}, err
    }

    // Create a random nonce with 8 * 12 bits of entropy...
    nonce := make([]byte, NonceLength)

    // ...using a CSPRNG...
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return []byte{}, []byte{}, err
    }

    // ...and try to initialize cipher...
    aesgcm, err := cipher.NewGCM(block)

    // ...if we fail, we return an error...
    if err != nil {
        return []byte{}, []byte{}, err
    }

    // Encrypt token and return ciphertext to caller
    return aesgcm.Seal(nil, nonce, plaintext, nil), nonce, nil

}
  
func DecryptGCM(ciphertext []byte, key []byte, nonce []byte) ([]byte, error) {
    // Create a key data structure
    block, err := aes.NewCipher(key)

    if err != nil {
        return []byte{}, err
    }

    // Initialize cipher...
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return []byte{}, err
    }

    // ...and try to decrypt...
    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)

    // ...no errors, hopefully...
    if err != nil {
        return []byte{}, err
    }

    // ...otherwise return dat shit to caller
    return plaintext, nil
}

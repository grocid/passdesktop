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
    "io"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "golang.org/x/crypto/pbkdf2"
    "encoding/hex"
  )

func LockToken(plaintext string, password string) (string, string, string, error) {
    // Generate random salt
    salt := make([]byte, 32)
    rand.Read(salt)
    // Generate key using PBKDF2
    key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)
    // Encrypt with AES-GCM-256
    ciphertext, nonce, err := EncryptGCM([]byte(plaintext), key)
    // Handle error gracefully
    if err != nil {
        return "", "", "", err
    }
    // Encode to hex strings
    encodedCipertext := hex.EncodeToString(ciphertext)
    encodedSalt := hex.EncodeToString(salt)
    encodedNonce := hex.EncodeToString(nonce)
    // and return to caller
    return encodedCipertext, encodedNonce, encodedSalt, nil
}

func UnlockToken(ciphertext string, password string, nonce string, salt string) (string, error) {
    // Decode the provided hex strings
    decodedSalt, _ := hex.DecodeString(salt)
    decodedNonce, _ := hex.DecodeString(nonce)
    decodedCipertext, _ := hex.DecodeString(ciphertext)
    // Generate the key from the salt and password...
    key := pbkdf2.Key([]byte(password), decodedSalt, 4096, 32, sha256.New)
    // Use nonce and key to decrypt ciphertext
    token, err := DecryptGCM(decodedCipertext, key, decodedNonce)
    // Take care of errors, i.e., if message authentication failed...
    if err != nil {
        return "", err
    }
    // ...otherwise, return to UI
    return string(token), nil
}

func EncryptGCM(plaintext []byte, key []byte) ([]byte, []byte, error) {
    // Create a key data structure
    block, err := aes.NewCipher(key)
    if err != nil {
        return []byte{}, []byte{}, err
    }
    // Create a random nonce with 8 * 12 bits of entropy...
    nonce := make([]byte, 12)
    // ...using a CSPRNG...
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return []byte{}, []byte{}, err
    }
    // ...and try to encrypt...
    aesgcm, err := cipher.NewGCM(block)
    // ...if we fail, we return an error...
    if err != nil {
        return []byte{}, []byte{}, err
    }
    // Return ciphertext to caller
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
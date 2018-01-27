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
    "math/big"
    "crypto/rand"
    "crypto/sha256"
    "golang.org/x/crypto/pbkdf2"
    "golang.org/x/crypto/argon2"
    "encoding/hex"
)

// The alphabet from which characters are drawn. The entropy is per
// symbol for this alphabet is roughly ~5.95 bits, so the default
// setting (length = 32) generates passwords with ~190 bits.
var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const (
    NonceLength = 12 // Only applicable to AES-GCM
    SaltLength = 32
    PBKDF2Iterations = 4096
    Argon2Time = 4
    Argon2Memory = 32 * 1024
    Argon2Threads = 4
    Argon2KeyLen = 32
)

func RandomPassword(length int) (string, error) {
    var result string

    // Create a random password of specified length.
    for {
        if len(result) >= length {
          return result, nil
        }
        // Read from from a CSPRNG.
        num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
        if err != nil {
            return "", err
        }
        // Add from alphabet.
        n := num.Int64()
        result += string(alphabet[n])
    }
}

func DeriveKey(password []byte, salt []byte) []byte {
    // Setting to specify key derivation algorithm.
    if UseArgon2ForKeyDerivation {
        return argon2.Key(password, salt, 
                          Argon2Time,
                          Argon2Memory,
                          Argon2Threads,
                          Argon2KeyLen)
    } else {
        // Fallback.
        return pbkdf2.Key(password, salt, 
                          PBKDF2Iterations, SaltLength, 
                          sha256.New)
    }
}

func LockToken(plaintext string, password string) (string, string, string, error) {
    // Generate random salt.
    salt := make([]byte, SaltLength)
    rand.Read(salt)

    // Derive key from password + salt.
    key := DeriveKey([]byte(password), salt)

    // Encrypt with AES-GCM-256.
    ciphertext, nonce, err := EncryptGCM([]byte(plaintext), key)

    // Handle error gracefully.
    if err != nil {
        return "", "", "", err
    }

    // Encode to hex strings...
    encodedCipertext := hex.EncodeToString(ciphertext)
    encodedSalt := hex.EncodeToString(salt)
    encodedNonce := hex.EncodeToString(nonce)

    // ...and return to caller.
    return encodedCipertext, encodedNonce, encodedSalt, nil
}

func UnlockToken(ciphertext string, password string, nonce string, salt string) (string, error) {
    // Decode the provided hex strings
    decodedSalt, _ := hex.DecodeString(salt)
    decodedNonce, _ := hex.DecodeString(nonce)
    decodedCipertext, _ := hex.DecodeString(ciphertext)

    // Derive key from password + salt.
    key := DeriveKey([]byte(password), decodedSalt)

    // Use nonce and key to decrypt ciphertext
    token, err := DecryptGCM(decodedCipertext, key, decodedNonce)

    // Take care of errors, i.e., if message authentication failed...
    if err != nil {
        return "", err
    }

    // ...otherwise, return to UI
    return string(token), nil
}


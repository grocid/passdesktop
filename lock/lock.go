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
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "golang.org/x/crypto/argon2"
    "golang.org/x/crypto/pbkdf2"
    "math/big"
)

const (
    SaltLength = 32

    PBKDF2Iterations = 4096
    Argon2Time       = 4
    Argon2Memory     = 32 * 1024
    Argon2Threads    = 4
    Argon2KeyLen     = 32

    UseArgon2ForKeyDerivation      = true
    DefaultGeneratedPasswordLength = 32
)

// The alphabet from which characters are drawn. The entropy is per
// symbol for this alphabet is roughly ~5.95 bits, so the default
// setting (length = 32) generates passwords with ~190 bits.
var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type Lock struct {
    Key  []byte
    Salt []byte
}

func New(password string, salt []byte) Lock {
    key := DeriveKey([]byte(password), salt)
    l := Lock{
        Key:  key,
        Salt: salt}
    return l
}

func Entropy(length int) []byte {
    r := make([]byte, length)
    rand.Read(r)

    return r
}

func EntropyAlphabet(length int) string {
    var result string

    // Create a random password of specified length.
    for {
        if len(result) >= length {
            return result
        }
        // Read from from a CSPRNG.
        num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))

        if err != nil {
            return ""
        }
        // Add from alphabet.
        n := num.Int64()
        result += string(alphabet[n])
    }
}

func DeriveKey(password []byte, salt []byte) []byte {
    // Setting to specify key derivation algorithm.
    if UseArgon2ForKeyDerivation {
        // Return derived key from Argon2.
        key := argon2.Key(password, salt,
            Argon2Time,
            Argon2Memory,
            Argon2Threads,
            Argon2KeyLen)
        return key
    } else {
        // Fallback to PBKDF2.
        key := pbkdf2.Key(password, salt,
            PBKDF2Iterations, SaltLength,
            sha256.New)
        return key
    }
}

func (l *Lock) EncryptAndEncodeBase64(plaintext string) ([]byte, error) {
    // Encrypt with AEAD.
    ciphertext, err := Chacha20Poly1305Encrypt([]byte(plaintext), l.Key)

    // Return base64 encoded
    return []byte(base64.StdEncoding.EncodeToString(ciphertext)), err
}

func (l *Lock) Base64DecodeAndDecrypt(base64Ciphertext []byte) (string, error) {
    // Decode the provided base64 string, use strings for simplicity
    ciphertext, err := base64.StdEncoding.DecodeString(string(base64Ciphertext))

    // Use key to decrypt ciphertext
    token, err := Chacha20Poly1305Decrypt(ciphertext, l.Key)
    return string(token), err
}

func (l *Lock) EncryptAndEncodeHex(plaintext string) (string, error) {
    // Encrypt with AEAD.
    ciphertext, err := Chacha20Poly1305Encrypt([]byte(plaintext), l.Key)

    // Return hexencoded
    return string(hex.EncodeToString(ciphertext)), err
}

func (l *Lock) HexDecodeAndDecrypt(hexCiphertext string) (string, error) {
    // Decode the provided hex strings
    ciphertext, _ := hex.DecodeString(hexCiphertext)

    // Use key to decrypt ciphertext
    token, err := Chacha20Poly1305Decrypt(ciphertext, l.Key)

    return string(token), err
}

func (l *Lock) LockToken(token string) (string, string, error) {
    // Generate random salt.
    l.Salt = Entropy(SaltLength)

    hexToken, err := l.EncryptAndEncodeHex(token)

    if err != nil {
        return "", "", err
    }

    hexSalt := hex.EncodeToString(l.Salt)

    // ...and return to caller.
    return hexToken, hexSalt, nil
}

func (l *Lock) UnlockToken(hexToken string) (string, error) {

    token, err := l.HexDecodeAndDecrypt(hexToken)

    // Take care of errors, i.e., if message authentication failed...
    if err != nil {
        return "", err
    }

    // ...otherwise, return to UI
    return token, nil
}

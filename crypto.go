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
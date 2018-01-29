package main

import (
    "crypto/hmac"
    "crypto/sha1"
    "encoding/base32"
    "encoding/binary"
    "fmt"
    "time"
)

func ComputeOTPCode(secret string) string {

    key, err := base32.StdEncoding.DecodeString(secret)
    if err != nil {
        return "n/a"
    }

    hash := hmac.New(sha1.New, key)
    err = binary.Write(hash, binary.BigEndian,
        int64(time.Now().Unix()/30))

    if err != nil {
        return "n/a"
    }

    h := hash.Sum(nil)
    offset := h[19] & 0x0f
    truncated := binary.BigEndian.Uint32(h[offset:offset+4]) & 0x7fffffff
    return fmt.Sprintf("%06d", int(truncated%1000000))
}

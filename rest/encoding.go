package rest

import "strings"

const (
    UserCredentialsSuffix = ",0"
    OTPSuffix             = ",1"
    FileSuffix            = ",2"
    SignSuffix            = ",3"
    AccountLabel          = "Account"
    OTPLabel              = "OTP"
    FileLabel             = "File"
    SignLabel             = "Sign"
)

func DecodeName(name string) (string, string, string) {

    if strings.HasSuffix(name, FileSuffix) {
        return name[:len(name)-2], FileLabel, "file"
    }

    if strings.HasSuffix(name, OTPSuffix) {
        return name[:len(name)-2], OTPLabel, "otp"
    }

    if strings.HasSuffix(name, SignSuffix) {
        return name[:len(name)-2], SignLabel, "sign"
    }

    return name, AccountLabel, name
}

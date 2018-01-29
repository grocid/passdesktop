package main

import (
	"crypto/md5" // used only for fingerprint generation!
	"fmt"
	"os"
	"strings"
)

const (
	ImagePathSuffix = "/../Resources/iconpack/"

	OTPSuffix  = "__otp"
	FileSuffix = "__file"
)

func GetFingerprint(data []byte) string {
	h := md5.New()
	hashDigest := h.Sum(data)

	grid := `<div class="grid-container">`
	for i := 0; i < 16; i++ {
		item := fmt.Sprintf(`<div class="grid-item" 
                                  style="background-color: rgba(255, 255, 255, 0.%v)">
                            </div>`, 10*int(hashDigest[i])/255)
		grid = grid + item
	}

	return grid + `</div>`
}

func GetImage(name string) string {
	// Get the path
	imagePath := pass.FullPath + ImagePathSuffix

	// Some ugly solution since the fallback on image not found does not work...
	image := name
	if _, err := os.Stat(imagePath + name + ".png"); os.IsNotExist(err) {
		image = "default"
	}
	return imagePath + image + ".png"
}

func GetTypeFromName(name string) int {
	if strings.HasSuffix(name, FileSuffix) {
		return TypeFile
	}
	if strings.HasSuffix(name, OTPSuffix) {
		return TypeOTP
	}
	return TypeUserCredentials
}

func ImageFromType(name string, entryType int) string {
	switch entryType {
	case TypeFile:
		return GetImage("file")
	case TypeOTP:
		return GetImage("otp")
	default:
		return GetImage(name)
	}
}

func GetDescriptionFromType(entryType int) string {
	switch entryType {
	case TypeFile:
		return "File"
	case TypeOTP:
		return "OTP"
	default:
		return ""
	}
}

func RemoveTypeFromName(name string) string {
	ns := len(name)
	if strings.HasSuffix(name, OTPSuffix) {
		return name[:ns-len(OTPSuffix)]
	}
	if strings.HasSuffix(name, FileSuffix) {
		return name[:ns-len(FileSuffix)]
	}
	return name
}

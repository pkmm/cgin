package util

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// MD5 return string hex upper case of md5 result.
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

package util

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 return string hex upper case of md5 result.
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

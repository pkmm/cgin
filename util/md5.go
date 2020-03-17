package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%X", hash.Sum(nil))
}

func Md5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

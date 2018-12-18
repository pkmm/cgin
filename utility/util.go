package utility

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	snowflake.Epoch = time.Now().Unix()
}

// 实用函数

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Md5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func Signature(params map[string]string) string {
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var buffer bytes.Buffer
	for i, key := range keys {
		if i == 0 {
			buffer.WriteString(fmt.Sprintf("%s=%s", key, params[key]))
		} else {
			buffer.WriteString(fmt.Sprintf("&%s=%s", key, params[key]))
		}
	}
	return strings.ToUpper(Md5String(buffer.String()))
}

func GenerateSignatureAndId(params map[string]string) (string, string) {
	node, _ := snowflake.NewNode(1)
	uniqueId := node.Generate().String()
	ans := Signature(params)
	return Md5String(ans + uniqueId + "salt(can random it)"), uniqueId
}

func IpAddressOfLocal() string {
	netInfos, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range netInfos {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func SourceCodePath() string {
	path, _ := os.Getwd()
	return path
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ9876543210")
	N := len(letters)
	ans := make([]rune, length)
	for i := 0; i < length; i++ {
		ans = append(ans, letters[rand.Intn(N)])
	}
	return string(ans)
}

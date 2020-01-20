package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
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

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ9876543210")
	N := len(letters)
	ans := make([]rune, 0)
	for i := 0; i < length; i++ {
		ans = append(ans, letters[rand.Intn(N)])
	}
	return string(ans)
}

func GenerateToken(key uint64) string {
	now := time.Now().Unix()
	str := RandomString(100)
	newKey := fmt.Sprintf("%s_%s_%d_%d", str, Md5String(str), now, key)

	return Md5String(newKey)
}

/// 路径函数
func GetSourceCodePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func GetExecPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

// 成员相同的结构体进行拷贝
func BeanDeepCopy(src, des interface{}) {
	retstring, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(retstring), &des)
	if err != nil {
		panic(err)
	}
}

// like PHP data('Y-m-d H:i:s')
func DateTime() string {
	now := time.Now()
	return fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func Date() string {
	now := time.Now()
	return fmt.Sprintf("%02d-%02d-%02d", now.Year(), now.Month(), now.Day())
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false

}

package utility

import (
	"crypto/md5"
	"fmt"
	"net"
	"encoding/hex"
	"os"
	"sort"
	"github.com/bwmarrin/snowflake"
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
	var ans string
	for _, key := range keys {
		ans += key + "=" + params[key]
	}
	return Md5String(ans)
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
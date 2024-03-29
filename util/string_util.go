package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"time"
)

var snowflakeNode *snowflake.Node

func init() {
	snowflake.Epoch = time.Now().Unix()
	snowflakeNode, _ = snowflake.NewNode(1)
}

func GUID() string {
	return snowflakeNode.Generate().String()
}

func GetInt64() int64 {
	return snowflakeNode.Generate().Int64()
}

// Signature 实用函数
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
	return MD5(buffer.String())
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

func GenerateToken(key uint64) string {
	now := time.Now().Unix()
	str := RandomString(100)
	newKey := fmt.Sprintf("%s_%s_%d_%d", str, MD5(str), now, key)

	return MD5(newKey)
}

// GetCurrentCodePath 获取当前代码的路径
func GetCurrentCodePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func GetExecPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

// GetExecutableCodePath 获取当前代码的路径
func GetExecutableCodePath() string {
	dir, _ := os.Executable()
	return filepath.Dir(dir)
}

// BeanDeepCopy 成员相同的结构体进行拷贝
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

// DateTime like PHP data('Y-m-d H:i:s')
func DateTime() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func Date() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func TimeString() string {
	return time.Now().Format("15:04:05")
}

// time func end

package util

import (
	"cgin/conf"
	"os"
)

func Getwd() string {
	if conf.IsProd() {
		wd, _ := os.Getwd()
		return wd
	} else {
		return GetCurrentCodePath()
	}
}

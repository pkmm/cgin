package main

import (
	"cgin/util"
	"testing"
)

// 测试结构体拷贝函数
func Test_structCopy(t *testing.T) {
	type A struct {
		Name string
		Age int
		Attrs []int
	}

	type B struct {
		Name string
		Age int
		Attrs []int
	}

	a := &A{
		Name: "Alice",
		Age: 24,
		Attrs: []int{1,2,45},
	}
	b := &B{}

	util.StructDeepCopy(a, b)
	if b.Name == a.Name && a.Age == b.Age {
		t.Log("结构体拷贝测试通过", b)
	} else {
		t.Error("结构体拷贝测试失败")
	}
}

func Test_RandomString(t *testing.T) {
	str := util.RandomString(21)
	if len(str) == 21 {
		t.Log("生成随机字符串测试通过", str)
	} else {
		t.Error("生成随机字符串测试失败", str)
	}
}
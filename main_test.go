package main

import (
	"cgin/util"
	"testing"
)

// 测试结构体拷贝函数
func TestStructCopy(t *testing.T) {
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
		t.Log("结构体测试通过", b)
	} else {
		t.Error("测试失败，结构体拷贝存在问题")
	}
}

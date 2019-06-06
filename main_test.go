package main

import (
	"cgin/service"
	"cgin/util"
	"testing"
)

// 测试结构体拷贝函数
func Test_structCopy(t *testing.T) {
	type A struct {
		Name  string
		Age   int
		Attrs []int
	}

	type B struct {
		Name  string
		Age   int
		Attrs []int
	}

	a := &A{
		Name:  "Alice",
		Age:   24,
		Attrs: []int{1, 2, 45},
	}
	b := &B{}

	util.BeanDeepCopy(a, b)
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

func Test_RestStudentSyncStatus(t *testing.T) {
	if err := service.StudentService.RestSyncStatus(); err != nil {
		t.Error("测试重置状态失败", err.Error())
	} else {
		t.Log("测试重置状态成功")
	}

}

func Test_GetNeedSyncScoreStudents(t *testing.T) {
	students, err := service.StudentService.GetStudentNeedSyncScore(0, 1)
	if err != nil || len(students) != 1 {
		if err != nil {
			t.Error(err.Error())
		} else {
			t.Log("没有学生需要同步， 测试通过")
		}
	} else {
		t.Log("测试获取学生成功")
	}
}

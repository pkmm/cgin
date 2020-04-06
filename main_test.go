package main

import (
	"cgin/conf"
	"cgin/model"
	"cgin/service"
	"cgin/util"
	"fmt"
	"io/ioutil"
	"runtime"
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
	if err := model.ResetStudentSyncScoreStatus(); err != nil {
		t.Error("测试重置状态失败", err.Error())
	} else {
		t.Log("测试重置状态成功")
	}

}

func Test_GetNeedSyncScoreStudents(t *testing.T) {
	err, students, total := new(model.Student).GetStudentsNeedSyncScore(0, 100)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("students number: ", len(*students) == 100)
	t.Log("total: ", total)
}

func TestSystemInfo(t *testing.T) {
	t.Log("cpu number:", runtime.NumCPU())
	t.Log("运行环境:", runtime.GOOS)
	t.Log("go root:", runtime.GOROOT())
}

func TestSomething(t *testing.T) {
	cookie := conf.WeiBoCookie()
	file, err := ioutil.ReadFile("static/xiaocc.jpeg")
	if err != nil {
		t.Fatal(err)
	}
	t4 := service.NewWeiBoStorage(cookie).UploadImage(file)
	fmt.Printf("%#v", t4)
	//var user model.User
	//fmt.Printf("%#v", &user == nil)

	//i, err := strconv.ParseInt("", 10, 64)
	//t.Log(err)
	//t.Log(i == 0)
	//
	//type name struct {
	//	H int
	//}
	//
	//a := make([]*name, 10)
	//a[0] = &name{23}
	//b := a
	//b[0] = &name{34}
	//fmt.Printf("a == b %#v\n", a[0])
	//fmt.Printf("a == b %#v", b[0])
	//
	//
	//// test gorm firstOrInit
	//var user model.User
	//conf.DB.FirstOrCreate(&user, model.User{OpenId:"33"}).Assign(model.User{RoleId: 4})
	//fmt.Printf("user model %#v", user)
}

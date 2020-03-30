package main

import (
	"cgin/model"
	"testing"
)

func TestUserModel(t *testing.T) {
	user := model.User{Id: 1}
	err, roles := user.GetRoles()
	t.Log(err)
	t.Log(roles)

	err, stu := user.GetStudent()
	t.Log(stu)
}

func TestResourceModel(t *testing.T) {
	//res := model.Resource{ID: 1}
	//err, pers := res.GetPermissions()
	//t.Log(err)
	//t.Log(pers)
	//
	//
	//t.Log("+++++++++++++++++")
	//err, rt := model.GetPermissionByResourceIdentityAndMethod("users", "post")
	//t.Log(rt)
	//urlPath := "/api/v1/students/334"
	//pattern := regexp.MustCompile(`api/v.*?/(.*?)/`)
	//resourceIdentity := pattern.FindStringSubmatch(urlPath)
	//t.Log(resourceIdentity)
	//user := service.AuthService.LoginFromMiniProgram("xxx=ccc")
	//t.Log(user)
	t.Log(model.AdminRoleId())
}
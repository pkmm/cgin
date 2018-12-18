package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xlstudio/wxbizdatacrypt"
	"io/ioutil"
	"net/http"
	"pkmm_gin/manager/userManager"
	"pkmm_gin/model"
	"pkmm_gin/utility"
)

// 微信小程序授权之后
// 解析用户的数据 返回accessToken userId 用于后续的请求的认证
func WxLogin(ctx *gin.Context) {
	var req model.WxLoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}

	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appId=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		model.GetConfig().WeChatSmallProgram.AppId,
		model.GetConfig().WeChatSmallProgram.Secret,
		req.Code,
	)
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer resp.Body.Close()
	var code2session model.Code2SessionResp
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	_ = json.Unmarshal(body, &code2session)

	// 解析出加密的数据
	pc := wxbizdatacrypt.WxBizDataCrypt{
		AppID:      model.GetConfig().WeChatSmallProgram.AppId,
		SessionKey: code2session.SessionKey,
	}

	result, err := pc.Decrypt(req.Data, req.Iv, true)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}

	type Tmp struct {
		OpenId    string `json:"OpenId"`
		NickName  string `json:"nickName"`
		Language  string `json:"language"`
		Province  string `json:"province"`
		AvatarUrl string `json:"avatarUrl"`
		City      string `json:"city"`
		Country   string `json:"country"`
		Gender    int    `json:"gender"`
	}
	var t Tmp
	_ = json.Unmarshal([]byte(result.(string)), &t)

	var wuser model.WeChatUser
	var user model.User
	if wuser = model.GetWeChatUserByOpenId(t.OpenId); wuser.Id == 0 {
		// 创建user
		user = model.User{
			Username: "微信登陆-系统创建",
			Password: utility.Md5String(utility.RandomString(10)),
			Salt:     utility.Md5String(utility.RandomString(10)),
		}
		user = model.CreateUser(user)

		// 新建一个微信用户
		wuser.OpenId = t.OpenId
		wuser.Nickname = t.NickName
		wuser.AvatarUrl = t.AvatarUrl
		wuser.City = t.City
		wuser.Province = t.Province
		wuser.Country = t.Country
		wuser.Gender = t.Gender
		wuser.Language = t.Language
		wuser.UserId = user.Id

		wuser = model.CreateWeChatUser(wuser)
	}
	if user.Id == 0 {
		user, _ = model.GetUserByOpenId(wuser.OpenId)
	}

	token := userManager.GenerateUserToken(&user, code2session.OpenId)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"user":    user,
		"session": model.Session{AccessToken: token, UserId: user.Id},
		"w_user":  wuser,
	})
}

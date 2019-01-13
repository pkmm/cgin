package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pkmm_gin/model"
	"pkmm_gin/service"
	"pkmm_gin/util"
)

func loginAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); err != nil {
		result.Msg = "preses login requests failed."
		result.Code = util.InvalidRequstParamter
		return
	}

	iv := arg["iv"].(string)
	code := arg["code"].(string)
	encryptedData := arg["encrypted_data"].(string)
	fmt.Println(iv, code, encryptedData)

	wechatUserInfo, err := util.DecodeWchatUserInfo(iv, code, encryptedData)

	if err != nil {
		result.Code = util.InvalidRequstParamter
		result.Msg = err.Error()
		return
	}

	user := service.User.GetUserByOpenId(wechatUserInfo.OpenId)
	if user == nil {
		user = &model.User{
			OpenId:   wechatUserInfo.OpenId,
			Nickname: wechatUserInfo.NickName,
		}
		if err = service.User.UpdateUser(user); err != nil {
			result.Code = util.InternalError
			result.Msg = "server internal error"
			return
		}
	}

	sess := service.UserSessionService.GetUserSession(user.ID)

	data := map[string]interface{}{}
	data["user"] = user
	data["session"] = sess

	result.Data = data

}

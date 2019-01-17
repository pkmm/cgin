package controller

import (
	"github.com/gin-gonic/gin"
	"pkmm_gin/errno"
	"pkmm_gin/model"
	"pkmm_gin/service"
	"pkmm_gin/util"
	"pkmm_gin/util/zf"
)

func loginAction(c *gin.Context) {
	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); err != nil {
		service.SendResponse(c, errno.ErrBind, nil)
		return
	}

	iv := arg["iv"].(string)
	code := arg["code"].(string)
	encryptedData := arg["encrypted_data"].(string)

	wechatUserInfo, err := util.DecodeWchatUserInfo(iv, code, encryptedData)

	if err != nil {
		service.SendResponse(c, errno.ErrBind.UpdateErrnoWithMsg(err.Error()), nil)
		return
	}

	user := service.User.GetUserByOpenId(wechatUserInfo.OpenId)
	if user == nil {
		user = &model.User{
			OpenId:   wechatUserInfo.OpenId,
			Nickname: wechatUserInfo.NickName,
		}
		if err = service.User.UpdateUser(user); err != nil {
			service.SendResponse(c, errno.InternalServerError, nil)
			return
		}
	}

	sess := service.UserSessionService.GetUserSession(user.ID)

	data := map[string]interface{}{}
	data["user"] = user
	data["token"] = sess

	service.SendResponse(c, errno.Success, data)
}

func getScoresAction(c *gin.Context) {
	defer func() {
		service.SendResponse(c, errno.InternalServerError, nil)
	}()

	data, exist := c.Get("user")
	if exist == false {
		service.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	user := data.(*model.User)

	scores := service.ScoreService.GetOwnScores(user.ID)
	service.SendResponse(c, errno.Success, scores)
}

func setAccountAction(c *gin.Context) {
	defer func() {
		service.SendResponse(c, errno.InternalServerError, nil)
	}()

	args := map[string]interface{}{}
	if err := c.BindJSON(&args); err != nil {
		service.SendResponse(c, errno.ErrBind, nil)
		return
	}
	num := args["num"].(string)
	pwd := args["pwd"].(string)

	checker, err := zf.NewCrawl(num, pwd)
	if err != nil {
		service.SendResponse(c, errno.ErrCheckZfAccountFailed.UpdateErrnoWithMsg(err.Error()), nil)
		return
	}

	if errMsg := checker.CheckAccount(); errMsg == "" {
		service.SendResponse(c, errno.Success.UpdateErrnoWithMsg("通过验证"), nil)
	} else  {
		service.SendResponse(c, errno.ErrCheckZfAccountFailed.UpdateErrnoWithMsg(errMsg), nil)
	}
}


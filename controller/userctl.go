package controller

import (
	"cgin/conf"
	"cgin/errno"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type userController struct{}

var UserController = &userController{}

//func (u *userController) LoginAction(c *gin.Context) {
//	arg := map[string]interface{}{}
//	if err := c.ShouldBindWith(&arg, binding.JSON); err != nil {
//		service.SendResponse(c, errno.InvalidParameters, nil)
//		return
//	}
//	var code string
//	var openid string
//	var ok bool
//	if code, ok = arg["code"].(string); !ok {
//		service.SendResponse(c, errno.InvalidParameters.AppendErrorMsg("code must have."), nil)
//		return
//	}
//	if openid, ok = arg["openid"].(string); !ok {
//		service.SendResponse(c, errno.InvalidParameters.AppendErrorMsg("encrypted data must have."), nil)
//		return
//	}
//
//	sign := util.Md5String("xiaocc_ai_liu_yan_lin" + conf.AppConfig.String("miniprogram_app_id") + openid)
//	if sign != code {
//		service.SendResponse(c, errno.InvalidParameters, nil)
//		return
//	}
//
//	user := service.User.GetUserByOpenId(openid)
//	if user == nil { // 没有找到用户 注册一个
//		user = &model.User{
//			OpenId:  openid,
//			CanSync: 1,
//		}
//		if err := service.User.UpdateUser(user); err != nil {
//			service.SendResponse(c, errno.UserNotFoundException, nil)
//			return
//		}
//	}
//
//	token, err := service.JWTSrv.GenerateToken(user)
//	if err != nil {
//		service.SendResponse(c, errno.GenerateJwtTokenFailed, nil)
//		return
//	}
//
//	data := map[string]interface{}{}
//	data["user"] = user
//	data["token"] = token
//
//	service.SendResponse(c, errno.Success, data)
//}
//
//func (u *userController) GetScoresAction(c *gin.Context) {
//	val, ok := c.Get("uid")
//	if ok == false {
//		service.SendResponse(c, errno.UserNotAuth, nil)
//		return
//	}
//	uid, _ := val.(uint64)
//	if uid == 0 {
//		service.SendResponse(c, errno.UserNotFoundException, nil)
//		return
//	}
//
//	scores := service.ScoreService.GetOwnScores(uint64(uid))
//
//	if len(scores) == 0 { // 提取为一个方法 todo
//		user := service.User.GetUser(uint64(uid))
//		worker, _ := zcmu.NewCrawl(user.Num, user.Pwd)
//		myScores, _ := worker.GetScores()
//		modelScores := make([]*model.Score, 0)
//		for _, s := range myScores {
//			score := &model.Score{
//				Xn:     s.Xn,
//				Xq:     s.Xq,
//				Kcmc:   s.Kcmc,
//				Cj:     s.Cj,
//				Jd:     s.Jd,
//				Cxcj:   s.Cxcj,
//				Bkcj:   s.Bkcj,
//				Xf:     s.Xf,
//				Type:   s.Type,
//				UserId: user.ID,
//			}
//			modelScores = append(modelScores, score)
//		}
//		service.SendResponse(c, errno.Success, modelScores)
//		go func() {
//			service.ScoreService.BatchCreate(modelScores)
//		}()
//		return
//	}
//	service.SendResponse(c, errno.Success, scores)
//}
//
//func (u *userController) SetAccountAction(c *gin.Context) {
//	args := map[string]interface{}{}
//	if err := c.ShouldBindWith(&args, binding.JSON); err != nil {
//		service.SendResponse(c, errno.InvalidParameters, nil)
//		return
//	}
//	var (
//		num, pwd string
//		ok       bool
//	)
//
//	if num, ok = args["num"].(string); !ok {
//		service.SendResponse(c, errno.InvalidParameters.AppendErrorMsg("num must have."), nil)
//		return
//	}
//	if pwd, ok = args["pwd"].(string); !ok {
//		service.SendResponse(c, errno.InvalidParameters.AppendErrorMsg("pwd must have."), nil)
//		return
//	}
//
//	checker, err := zcmu.NewCrawl(num, pwd)
//	if err != nil {
//		service.SendResponse(c, errno.CheckZfAccountFailedException.ReplaceErrorMsgWith(err.Error()), nil)
//		return
//	}
//
//	if errMsg := checker.CheckAccount(); errMsg == "" {
//
//		// 验证通过后更新token的值
//		var (
//			val interface{}
//			ok  bool
//		)
//		if val, ok = c.Get("uid"); !ok {
//			service.SendResponse(c, errno.UserNotAuth, nil)
//			return
//		}
//		uid, _ := val.(uint64)
//		user := service.User.GetUser(uid)
//		if user == nil {
//			service.SendResponse(c, errno.UserNotAuth, nil)
//			return
//		}
//		// 更新学生的num pwd
//		user.Num = num
//		user.Pwd = pwd
//		service.User.UpdateUser(user) // todo
//		token, err := service.JWTSrv.GenerateToken(user)
//		if err != nil {
//			service.SendResponse(c, errno.UserNotAuth, nil)
//			return
//		}
//		service.SendResponse(c, errno.Success.ReplaceErrorMsgWith("通过验证"), &map[string]interface{}{
//			"token": token,
//			"user":  user,
//		})
//	} else {
//		service.SendResponse(c, errno.CheckZfAccountFailedException.ReplaceErrorMsgWith(errMsg), nil)
//	}
//}
//
//func (u *userController) CheckTokenAction(c *gin.Context) {
//	val, ok := c.Get("uid")
//	if !ok {
//		service.SendResponse(c, errno.UserNotAuth, nil)
//		return
//	}
//	uid, _ := val.(uint64)
//	if uid == 0 {
//		service.SendResponse(c, errno.UserNotAuth, nil)
//		return
//	}
//	service.SendResponse(c, errno.Success, nil)
//}

func (u *userController) SendTemplateMsg(c *gin.Context) {
	params := map[string]interface{}{}
	if err := c.ShouldBindWith(&params, binding.JSON); err != nil {
		service.SendResponse(c, errno.InvalidParameters, nil)
		return
	}

	formId, ok := params["form_id"].(string)
	if !ok {
		service.SendResponseWithInvalidParameters(c, "form_id must supply")
		return
	}
	openId, ok := params["open_id"].(string)
	if !ok {
		service.SendResponseWithInvalidParameters(c, "open_id must supply")
		return
	}
	templateKeyData := &util.TemplateMsgData{}
	templateKeyData.Keyword1.Value = "11"
	templateKeyData.Keyword2.Value = "22"
	msg := &util.TemplateMsg{
		FormId:     formId,
		ToUser:     openId,
		TemplateId: conf.AppConfig.String("template_id"),
		Page:       conf.AppConfig.String("template_msg_open_page"),
		Data:       templateKeyData,
	}
	ret := util.SendUserTemplateMsg(msg)
	service.SendResponse(c, errno.Success, ret)
}

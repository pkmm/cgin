package v1

import (
	"cgin/conf"
	"cgin/controller/co"
	"cgin/controller/contextHelper"
	"cgin/controller/respobj"
	"cgin/errno"
	"cgin/model"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
)

// 微信小程序控制器，后台配置小程序各个页面，小程序获取配置接口

type miniProgramController struct{}

var MiniProgramController = &miniProgramController{}

// 发送模板消息
// @Tags Mini program
// @Summary 发送微信小程序订阅消息
// @Router /mini_program/send_template_msg [GET]
// @Param open_id query string true "用户的open id"
// @Success 200 object service.Response
// @Security ApiKeyAuth
func (m *miniProgramController) SendTemplateMsg(c *gin.Context) {
	helper := contextHelper.New(c)
	openId := helper.GetString("open_id")
	templateKeyData := service.Keys{}
	templateKeyData.Key1.Value = "背单词签到"
	templateKeyData.Key2.Value = "手机App签到"
	msg := &service.TemplateMsg{
		ToUser:     openId,
		TemplateId: conf.AppConfig.String("miniprogram.template_id"),
		Page:       conf.AppConfig.String("miniprogram.open_page"),
		Data:       templateKeyData,
	}
	ret := service.SendUserTemplateMsg(msg)
	helper.Response(ret)
}

// @Summary 配置菜单项
// @Tags Mini program
// @Security ApiKeyAuth
// @Router /mini_program/menus [post]
// @Success 200 object service.Response
// @Param menus body co.Menus true "配置小程序首页的菜单项"
func (m *miniProgramController) CreateMenus(c *gin.Context) {
	helper := contextHelper.New(c)
	var menus = &co.Menus{}
	helper.NeedAuthOrPanic()
	if err := c.ShouldBindJSON(&menus); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}

	// 有一个创建失败就会全部创建失败
	// TODO: 可以支持部分创建成功
	for _, menu := range menus.Menus {
		dbModel := model.Menu{}
		util.BeanDeepCopy(menu, &dbModel)
		err := dbModel.CreateMenu()
		if err != nil {
			panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
		}
	}
	_, dbMenus := model.GetActiveMenus()
	helper.Response(dbMenus)
}

// IndexPreference @summary 首页的配置
// @Tags Mini program
// @router /mini_program/index_preferences [get]
// @Security ApiKeyAuth
// @Success 200 object service.Response
func (m *miniProgramController) IndexPreference(c *gin.Context) {
	helper := contextHelper.New(c)
	// 菜单
	_, menus := model.GetActiveMenus()
	// 首页配置 slogan等
	_, indexConfig := new(model.IndexConfig).GetLatest()

	data := gin.H{
		"menus":        menus,
		"index_config": indexConfig,
	}
	helper.Response(data)
}

// CreateIndexConfig @Summary 首页slogan image等的配置信息
// @Tags Mini program
// @Router /mini_program/index_config [post]
// @Param config body co.IndexConfig true "小程序首页配置"
// @Success 200 object service.Response
// @Security ApiKeyAuth
// @Produce json
// @Accept json
func (m *miniProgramController) CreateIndexConfig(c *gin.Context) {
	helper := contextHelper.New(c)
	config := &co.IndexConfig{}
	if err := c.BindJSON(config); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	ic := model.IndexConfig{
		Slogan:     config.Slogan,
		Motto:      config.Motto,
		ImageUrl:   config.ImageUrl,
		ImageStyle: config.ImageStyle,
	}
	_, savedConfig := ic.Save()
	helper.Response(savedConfig)
}

// GetNotifications @Security ApiKeyAuth
// @Tags Mini program
// @Summary 获取notifications 分页查询
// @Router /mini_program/notifications [get]
// @Param pagingInfo query co.PageLimitOffset true "分页参数"
// @Success 200 object service.Response
// @Produce json
func (m *miniProgramController) GetNotifications(c *gin.Context) {
	helper := contextHelper.New(c)
	size := helper.GetInt("size")
	page := helper.GetInt("page")
	err, notifications, total := new(model.Notification).GetList(page, size)
	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
	helper.Response(gin.H{"notifications": notifications, "total": total})
}

// UpdateOrCreateNotification @Security ApiKeyAuth
// @Tags Mini program
// @Summary 更新创建一个notification
// @Router /mini_program/notifications [put]
// @Success 200 object service.Response
// @Param notification body co.Notification true "one notification"
// @Produce json
func (m *miniProgramController) UpdateOrCreateNotification(c *gin.Context) {
	helper := contextHelper.New(c)
	helper.GetAuthUserId()
	notification := &co.Notification{}
	if err := c.ShouldBindJSON(notification); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	dbModel := model.Notification{}
	util.BeanDeepCopy(notification, &dbModel)
	if notification.Id == 0 {
		// 没有传id认为是创建一个notification
		_, saved := dbModel.CreateNotification()
		helper.Response(saved)
	} else {
		// 更新已经存在的一个notification
		_, updated := dbModel.UpdateNotification(notification.Id)
		helper.Response(updated)
	}
}

// GetSponsors @Summary 查看赞助我的人
// @Tags Mini program
// @Security ApiKeyAuth
// @Router /mini_program/sponsors [get]
// @Param pagingInfo query co.PageLimitOffset true "分页参数"
// @Success 200 object service.Response
// @Produce json
func (m *miniProgramController) GetSponsors(c *gin.Context) {
	helper := contextHelper.New(c)
	err, data, total := new(model.Sponsor).GetList(helper.GetInt("page"), helper.GetInt("size"))
	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
	sponsors := data.([]*model.Sponsor)
	var result = make([]*respobj.Sponsor, len(sponsors))

	// 处理用户名
	processUsername := func(user *model.User) string {
		if user == nil || len(user.Username) == 0 {
			return "不愿意透露姓名的: Alice"
		}
		return user.Username
	}

	for ind, s := range sponsors {
		o := &respobj.Sponsor{
			Id:        s.Id,
			Money:     s.Money,
			CreatedAt: s.CreatedAt,
			Username:  processUsername(s.User),
		}
		if s.User != nil {
			o.OpenId = s.User.OpenId
		}
		result[ind] = o
	}
	helper.Response(gin.H{
		"sponsors": result,
		"total":    total,
	})
}

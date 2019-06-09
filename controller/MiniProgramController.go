package controller

import (
	"cgin/conf"
	"cgin/controller/co"
	"cgin/errno"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 微信小程序控制器，后台配置小程序各个页面，小程序获取配置接口

type miniProgramController struct {
	BaseController
}

var MiniProgramController = &miniProgramController{}

// 发送模板消息
func (m *miniProgramController) SendTemplateMsg(c *gin.Context) {
	m.ProcessParams(c)

	formId, ok := m.Params["form_id"].(string)
	if !ok {
		panic(errno.NormalException.ReplaceErrorMsgWith("require form id."))
	}
	openId, ok := m.Params["open_id"].(string)
	if !ok {
		panic(errno.NormalException.ReplaceErrorMsgWith("require open id."))
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
	m.Response(c, ret)
}

// 配置小程序首页的菜单项
func (m *miniProgramController) DisposeMenu(c *gin.Context) {
	m.GetAuthUserId(c)
	var menus []co.Menu
	if err := c.BindJSON(&menus); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}

	// 有一个创建失败就会全部创建失败 TODO: 可以支持部分创建成功
	var savedMenus []interface{}
	for _, menu := range menus {
		savedMenu, errI := service.MiniProgramService.DisposeMenu(menu.Desp, menu.Title, menu.Icon, menu.ActionType, menu.ActionValue)
		if errI != nil {
			panic(errno.NormalException.ReplaceErrorMsgWith(errI.Error()))
		}
		savedMenus = append(savedMenus, savedMenu)
	}
	m.Response(c, savedMenus)
}

// 首页的配置
func (m *miniProgramController) GetIndexPreference(c *gin.Context) {
	// 菜单
	menus := service.MiniProgramService.GetAllActiveMenus()
	// 首页配置 slogan等
	indexConfig := service.MiniProgramService.GetIndexConfig()

	data := gin.H{
		"menus":        menus,
		"index_config": indexConfig,
	}
	m.Response(c, data)
}

// 首页solgan image等的配置信息
func (m *miniProgramController) SetIndexConfig(c *gin.Context) {
	m.GetAuthUserId(c)
	config := &co.IndexConfig{}
	if err := c.BindJSON(config); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	savedConfig := service.MiniProgramService.SetIndexConfig(config.Slogan, config.ImageUrl, config.ImageStyle)
	m.Response(c, savedConfig)
}

// 获取notification 默认是显示最新的10条
func (m *miniProgramController) GetNotification(c *gin.Context) {
	limit, err := strconv.ParseUint(c.DefaultQuery("count", "10"), 10, 64)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	notifications := service.MiniProgramService.GetNotifications(limit)
	m.Response(c, gin.H{"notifications": notifications})
}

// 更新或者创建一个notification
func (m *miniProgramController) UpdateOrCreateNotification(c *gin.Context) {
	m.GetAuthUserId(c)
	notification := &co.Notification{}
	if err := c.BindJSON(notification); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	if notification.Id == 0 {
		// 没有传id认为是创建一个notification
		createdNotification := service.MiniProgramService.SaveNotification(notification.Content, notification.StartAt, notification.EndAt)
		m.Response(c, createdNotification)
		return
	}
	// 更新已经存在的一个notification
	savedNotification := service.MiniProgramService.UpdateNotification(notification.Id, notification.Content, notification.StartAt, notification.EndAt)
	m.Response(c, savedNotification)
}

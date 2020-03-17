package v1

import (
	"cgin/conf"
	"cgin/controller/co"
	"cgin/controller/context_helper"
	"cgin/controller/respobj"
	"cgin/errno"
	"cgin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 微信小程序控制器，后台配置小程序各个页面，小程序获取配置接口

type miniProgramController struct {}

var MiniProgramController = &miniProgramController{}

// 发送模板消息
// @Summary 发送微信小程序订阅消息
// @Router /mini_program/send_template_msg [GET]
// @Param open_id query string true "用户的open id"
// @Success 200 {object} service.Response
func (m *miniProgramController) SendTemplateMsg(c *gin.Context) {
	helper := context_helper.New(c)
	//formId := helper.GetString("form_id")
	openId := helper.GetString("open_id")
	templateKeyData := service.TemplateMsgData{}
	templateKeyData.Key1.Value = "背单词签到"
	templateKeyData.Key2.Value = "手机App签到"
	fmt.Printf("%#v\n", templateKeyData)
	msg := &service.TemplateMsg{
		//FormId:     formId,
		ToUser:     openId,
		TemplateId: conf.AppConfig.String("miniprogram.template_id"),
		Page:       conf.AppConfig.String("miniprogram.open_page"),
		Data:       templateKeyData,
	}
	ret := service.SendUserTemplateMsg(msg)
	helper.Response(ret)
}

// 配置小程序首页的菜单项
// @Summary 配置菜单项
// @Security ApiKeyAuth
// @Router /mini_program/config_menu [post]
// @Success 200 {object} service.Response
// @Param menus body co.Menus true "dispose menus"
func (m *miniProgramController) DisposeMenu(c *gin.Context) {
	helper := context_helper.New(c)
	var menus = &co.Menus{}
	helper.NeedAuthOrPanic()
	if err := c.ShouldBindJSON(&menus); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}

	// 有一个创建失败就会全部创建失败 TODO: 可以支持部分创建成功
	var savedMenus []interface{}
	for _, menu := range menus.Menus {
		savedMenu, errI := service.MiniProgramService.DisposeMenu(menu.Desp, menu.Title, menu.Icon, menu.ActionType, menu.ActionValue)
		if errI != nil {
			panic(errno.NormalException.ReplaceErrorMsgWith(errI.Error()))
		}
		savedMenus = append(savedMenus, savedMenu)
	}
	helper.Response(savedMenus)
}

// 首页的配置
func (m *miniProgramController) GetIndexPreference(c *gin.Context) {
	helper := context_helper.New(c)
	// 菜单
	menus := service.MiniProgramService.GetAllActiveMenus()
	// 首页配置 slogan等
	indexConfig := service.MiniProgramService.GetIndexConfig()

	data := gin.H{
		"menus":        menus,
		"index_config": indexConfig,
	}
	helper.Response(data)
}

// 首页slogan image等的配置信息
func (m *miniProgramController) SetIndexConfig(c *gin.Context) {
	helper := context_helper.New(c)
	config := &co.IndexConfig{}
	if err := c.BindJSON(config); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	savedConfig := service.MiniProgramService.SetIndexConfig(config.Slogan, config.ImageUrl, config.ImageStyle)
	helper.Response(savedConfig)
}

// 获取notification 默认是显示最新的10条
// @Security ApiKeyAuth
// @Summary 获取notification
// @Router /mini_program/get_notifications [get]
// @Success 200 {object} service.Response
// @Produce json
func (m *miniProgramController) GetNotification(c *gin.Context) {
	helper := context_helper.New(c)
	limit, err := strconv.ParseUint(c.DefaultQuery("count", "10"), 10, 64)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	notifications := service.MiniProgramService.GetNotifications(limit)
	helper.Response(gin.H{"notifications": notifications})
}

// 更新或者创建一个notification
// @Security ApiKeyAuth
// @Summary 更新创建一个notification
// @Router /mini_program/change_notification [post]
// @Success 200 {object} service.Response
// @Param notification body co.Notification true "one notification"
// @Produce json
func (m *miniProgramController) UpdateOrCreateNotification(c *gin.Context) {
	helper := context_helper.New(c)
	helper.GetAuthUserId()
	notification := &co.Notification{}
	if err := c.ShouldBindJSON(notification); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	if notification.Id == 0 {
		// 没有传id认为是创建一个notification
		createdNotification := service.MiniProgramService.SaveNotification(notification.Content, notification.StartAt, notification.EndAt)
		helper.Response(createdNotification)
		return
	}
	// 更新已经存在的一个notification
	savedNotification := service.MiniProgramService.UpdateNotification(notification.Id, notification.Content, notification.StartAt, notification.EndAt)
	helper.Response(savedNotification)
}

// 赞助的人
// @Summary 查看赞助我的人
// @Router /mini_program/get_sponsors [get]
// @Success 200 {object} service.Response
// @Produce json
func (m *miniProgramController) GetSponsors(c *gin.Context) {
	helper := context_helper.New(c)
	sponsors := service.MiniProgramService.GetSponsors()
	var result = make([]*respobj.Sponsor, len(sponsors))
	for ind, s := range sponsors {
		o := &respobj.Sponsor{
			Id:        s.Id,
			Money:     s.Money,
			CreatedAt: s.CreatedAt,
		}
		if s.User != nil {
			o.OpenId = s.User.OpenId
		}
		result[ind] = o
	}
	helper.Response(gin.H{
		"sponsors": result,
	})
}

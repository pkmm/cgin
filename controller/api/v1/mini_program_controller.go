package v1

import (
	"cgin/conf"
	"cgin/controller/co"
	"cgin/controller/context_helper"
	"cgin/controller/respobj"
	"cgin/errno"
	"cgin/model"
	"cgin/model/modelInterface"
	"cgin/service"
	"cgin/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 微信小程序控制器，后台配置小程序各个页面，小程序获取配置接口

type miniProgramController struct{}

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

// @summary 首页的配置
// @router /mini_program/get_index_preference [get]
// @Security ApiKeyAuth
// @Success 200 {object} service.Response
func (m *miniProgramController) GetIndexPreference(c *gin.Context) {
	helper := context_helper.New(c)
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

// @Summary 首页slogan image等的配置信息
// @Router /mini_program/set_index_config [post]
// @Param config body co.IndexConfig true "小程序首页配置"
// @Success 200 {object} service.Response
// @Security ApiKeyAuth
// @Produce json
// @Accept json
func (m *miniProgramController) SetIndexConfig(c *gin.Context) {
	helper := context_helper.New(c)
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

// @Security ApiKeyAuth
// @Summary 获取notifications 分页查询
// @Router /mini_program/get_notifications [get]
// @Param pagingInfo query co.PageLimitOffset true "分页参数"
// @Success 200 {object} service.Response
// @Produce json
func (m *miniProgramController) GetNotifications(c *gin.Context) {
	helper := context_helper.New(c)
	size := helper.GetInt("size")
	page := helper.GetInt("page")
	err, notifications, total := new(model.Notification).GetList(modelInterface.PageSizeInfo{
		Page:     page,
		PageSize: size,
	})
	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
	helper.Response(gin.H{"notifications": notifications, "total": total})
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

// @Summary 查看赞助我的人
// @Security ApiKeyAuth
// @Router /mini_program/get_sponsors [get]
// @Param pagingInfo query co.PageLimitOffset true "分页参数"
// @Success 200 {object} service.Response
// @Produce json
func (m *miniProgramController) GetSponsors(c *gin.Context) {
	helper := context_helper.New(c)
	err, data, total := new(model.Sponsor).GetList(modelInterface.PageSizeInfo{
		Page:     helper.GetInt("page"),
		PageSize: helper.GetInt("size"),
	})
	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
	sponsors := data.([]*model.Sponsor)
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
		"total":    total,
	})
}

package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"hxyxx/lottery/services"
)

type AdminController struct{
	Ctx iris.Context
	ServiceUser services.UserService
	ServiceGift services.GiftService
	ServiceCode services.CodeService
	ServiceResult services.ResultService
	ServiceUserDay services.UserdayService
	ServiceBlackip services.BlackipService
}
func (c *AdminController) Get()mvc.Result{
	return mvc.View{//返回模板对象
		Name: "admin/index.html", //后台首页
		Data: iris.Map{
			"Title":"管理后台",
			"Channel":"",//频道
		},
		Layout: "admin/layout.html",//模板
	}
}
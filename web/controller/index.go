package controller

import (
	"fmt"
	"github.com/kataras/iris"
	"hxyxx/lottery/comm"
	"hxyxx/lottery/models"
	"hxyxx/lottery/services"
)

type IndexController struct{
	Ctx iris.Context
	ServiceUser services.UserService
	ServiceGift services.GiftService
	ServiceCode services.CodeService
	ServiceResult services.ResultService
	ServiceUserDay services.UserdayService
	ServiceBlackip services.BlackipService
}

//http://localhost:8080/
func (c *IndexController) Get() string {
	c.Ctx.Header("Content-Type","text/html")
	return "welcome to Go抽奖系统，<a href='/public/index.html'>开始抽奖</a>"
}
//获取奖品信息
func (c *IndexController) GetGifts() map[string]interface{} {
	rs := make(map[string]interface{},0)
	rs["code"] = 0
	rs["msg"] = ""
	datalist := c.ServiceGift.GetAll(false)
	list := make([]models.LtGift,0)
	for _,data := range datalist{
		if data.SysStatus ==0 {
			list = append(list,data)
		}
	}
	rs["gifts"] = list
	return rs
}
//最新的获奖列表
func (c *IndexController) GetNewPrize() map[string]interface{}{
	rs := make(map[string]interface{},0)
	rs["code"] = 0
	rs["msg"] = ""
	//todo
	return rs
}
//登陆
func (c *IndexController) GetLogin(){
	uid := comm.Random(100000)
	loginuser := &models.ObjLoginuser{
		Uid: uid,
		Username: fmt.Sprintf("admin-%d",uid),
		Now: comm.NowUnix(),
		Ip: comm.ClientIP(c.Ctx.Request()),
	}
	comm.SetLoginuser(c.Ctx.ResponseWriter(),loginuser)
	comm.Redirect(c.Ctx.ResponseWriter(),
		"/public/index.html?from=login")
}
//退出
func (c *IndexController) GetLogout(){
	comm.SetLoginuser(c.Ctx.ResponseWriter(),nil)
	comm.Redirect(c.Ctx.ResponseWriter(),
		"/public/index.html?from=logout")
}
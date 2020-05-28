package controller

import (
	"fmt"
	"hxyxx/lottery/comm"
	"hxyxx/lottery/conf"
	"hxyxx/lottery/models"
	"hxyxx/lottery/web/utils"
)

//抽奖
func (c *IndexController) GetLucky() map[string]interface{}{
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	loginuser := comm.GetLoginUser(c.Ctx.Request())
	if loginuser == nil || loginuser.Uid<1{
		rs["code"] = 101
		rs["msg"] = "请先登录，再来抽奖"
	}
	//用户抽奖分布式锁，避免某个用户连续点击
	ok := utils.LockLucky(loginuser.Uid)
	if ok{
		defer utils.UnlockLucky(loginuser.Uid)
	}else {
		rs["code"] = 102
		rs["msg"] = "请先登录，再来抽奖"
		return rs
	}
	//验证用户今日参与次数
	ok = c.checkUserday(loginuser.Uid)
	if !ok{
		rs["code"] = 103
		rs["msg"] = "今日抽奖次数已用完，明天再来吧"
		return rs
	}
	//验证ip今日参与次数
	ip := comm.ClientIP(c.Ctx.Request())
	//每次请求次数+1
	ipDayNum :=utils.IncrIpLuckyNum(ip)
	if ipDayNum >conf.IpLimitMax{
		rs["code"] = 104
		rs["msg"] = "相同ip参与次数太多，明天再来参加吧"
		return rs
	}
	limitBlack := false
	if ipDayNum >conf.IpPrizeMax{
		limitBlack = true
	}
	//验证ip黑名单
	var blackipInfo *models.LtBlackip
	if !limitBlack{
		ok,blackipInfo = c.checkBlakcip(ip)
		if !ok{
			fmt.Println("黑名单中的IP",ip,limitBlack)
		}
	}
	//验证用户黑名单
	var userInfo *models.LtUser
	if !limitBlack{
		ok,userInfo =c.checkBlackUser(loginuser.Uid )
		if !ok{
			fmt.Println("黑名单中的用户",loginuser.Uid,limitBlack)
			limitBlack = true
		}
	}

	//获得抽奖编码
	prizeCode := comm.Random(10000)
	//匹配奖品是否中奖
	prizeGift := c.prize(prizeCode,limitBlack)
	if prizeGift == nil || prizeGift.PrizeNum<0 ||(prizeGift.PrizeNum > 0 && prizeGift.LeftNum<=0){
		rs["code"] = 205
		rs["msg"] = "很遗憾，没有中奖，请下次再试"
		return rs
	}
	//有限制的奖品发放
	if prizeGift.LeftNum > 0{
		ok = utils.PrizeGift(prizeGift.Id,prizeGift.LeftNum)
		if !ok{
			rs["code"]=207
			rs["msg"] = "很遗憾请下次再试"
			return rs
		}
	}
	//不同编码的优惠券发放
	if prizeGift.Gtype == conf.GtypeCodeDiff{
		code := utils.PrizeCodeDiff(prizeGift.Id,c.ServiceCode)
		if code == ""{
			rs["code"] = 208
			rs["msg"] = "很遗憾，没有中奖，请下次再试"
			return rs
		}
		prizeGift.Gdata = code
	}
	//记录中奖纪录
	//返回抽奖结果
	return rs
}

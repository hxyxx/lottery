
/**
 * 年会抽奖程序
 * 不是线程安全
 * 基础功能：
 * 1 /import 导入参与名单作为抽奖的用户
 * 2 /lucky 从名单中随机抽取用户
 * 测试方法：
 * curl http://localhost:8080/
 * curl --data "users=yifan,yifan2" http://localhost:8080/import
 * curl http://localhost:8080/lucky
 * @author 一凡Sir
 */
package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"strings"
	"time"
)

var usrList []string
type lotteryController struct{
	Ctx iris.Context
}
func newApp() *iris.Application{
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main(){
	app := newApp()
	usrList = []string{}
	//记得加冒号
	app.Run(iris.Addr(":8080"))
}
// GET http://localhost:8080/
func (lc *lotteryController) Get() string{
	count := len(usrList)
	return fmt.Sprintf("当前抽奖总人数为： %d\n",count)
}
// POST http://localhost:8080/import
func (lc *lotteryController) PostImport() string{
	strUser := lc.Ctx.FormValue("users")
	fmt.Print(strUser+"hhh")
	users := strings.Split(strUser,",")
	count := len(usrList)
	for _,u := range users{
		u = strings.TrimSpace(u)
		if len(u)>0{
			usrList = append(usrList, u)
		}
	}
	count1 := len(usrList)
	return fmt.Sprintf("倒入之前抽奖总人数为 :%d,倒入之后抽奖总人数为 : %d\n",count,count1)
}
//GET http://localhost:8080/lucky
func (lc *lotteryController) GetLucky() string{
	fmt.Sprintf("hhhhhhh")
	count := len(usrList)
	if count>1{
		seed := time.Now().UnixNano()
		//取到幸运儿
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		user := usrList[index]
		usrList = append(usrList[0:index],usrList[index+1:]...)
		return fmt.Sprintf("当前中奖用户为 :%s,抽奖剩余总人数为:%d\n",user,len(usrList))
	}else if count ==1{//奖池中只有一个用户
		user := usrList[0]
		usrList= []string{}
		return fmt.Sprintf("当前中奖用户为 :%s,抽奖剩余总人数为:%d\n",user,len(usrList))
	}else{//奖池中没有用户
		return fmt.Sprintf("奖池中已经没有用户，请先通过 /import 倒入用户\n")

	}
}

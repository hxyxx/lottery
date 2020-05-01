package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"time"
)

/**
即开即得型彩票 http://localhost:8080/
双色球自选型彩票 http://localhost:8080/twoball
*/
type lotteryController struct{
	Ctx *iris.Context
}
func newApp() *iris.Application{
	app := iris.New()
	mvc.New((app).Party("/")).Handle(&lotteryController{})
	return app
}
func main(){
	app := newApp()
	app.Run(iris.Addr(":8080"))
}
//即开即得型
func (lc *lotteryController) Get() string{
	var prize string
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Int31n(10)
	switch  {
	case code == 1:
		prize = "一等奖"
	case code >= 2 && code <= 3:
		prize = "二等奖"
	case code >= 4 && code <= 6:
		prize = "三等奖"
	default:
		return fmt.Sprintf("尾号为1获得一等奖<br/>"+
			"尾号为2或者3获得二等奖<br/>"+
			"尾号为4、5、6获得三等奖<br/>"+
			"您的尾号为%d"+"很遗憾您没有中奖",code)
	}
	return fmt.Sprintf("尾号为1获得一等奖<br/>"+
		"尾号为2或者3获得二等奖<br/>"+
		"尾号为4、5、6获得三等奖<br/>"+
		"您的尾号为%d"+"恭喜您获得%s",code,prize)
}
//双色球型
func (lc *lotteryController) GetTwoball() string{
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	//中奖号码
	var code [7]int
	for i:= 0;i<6;i++{
		code[i] = r.Intn(33)+1
	}
	//最后以为蓝色球
	code[6]=r.Intn(16)+1
	return fmt.Sprintf("本期中奖号码为 %v",code)
}

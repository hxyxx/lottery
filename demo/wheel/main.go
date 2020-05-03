/**
 * 大转盘程序
 * curl http://localhost:8080/
 * curl http://localhost:8080/debug
 * curl http://localhost:8080/prize
 * 固定几个奖品，不同的中奖概率或者总数量限制
 * 每一次转动抽奖，后端计算出这次抽奖的中奖情况，并返回对应的奖品信息
 *
 * 线程不安全，因为获奖概率低，并发更新库存的冲突很少能出现，不容易发现线程安全性问题
 * 压力测试：
 * wrk -t10 -c100 -d5 "http://localhost:8080/prize"
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

// 奖品中奖概率
type Prate struct {
	Rate int		// 万分之N的中奖概率
	Total int		// 总数量限制，0 表示无限数量
	CodeA int		// 中奖概率起始编码（包含）
	CodeB int		// 中奖概率终止编码（包含）
	Left int 		// 剩余数
}
// 奖品列表
var prizeList []string = []string{
	"一等奖，火星单程船票",
	"二等奖，凉飕飕南极之旅",
	"三等奖，iPhone一部",
	"",							// 没有中奖
}
// 奖品的中奖概率设置，与上面的 prizeList 对应的设置
var rateList []Prate = []Prate{
	Prate{1, 1, 0, 0, 1},
	Prate{2, 2, 1, 2, 2},
	Prate{5, 10, 3, 5, 10},
	Prate{100,0, 0, 9999, 0},
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}

// 抽奖的控制器
type lotteryController struct {
	Ctx iris.Context
}

// GET http://localhost:8080/
//返回奖品信息
func (lc *lotteryController) Get() string {
	lc.Ctx.Header("Content-Type", "text/html")
	return fmt.Sprintf("大转盘奖品列表：<br/> %s", strings.Join(prizeList, "<br/>\n"))
}

//Get http://localhost:8080/debug
//返回奖品中奖概率
func (lc *lotteryController) GetDebug() string{
	return fmt.Sprintf("获奖概率：%v\n",rateList)
}
//Get http://localhost:8080/prize
func (lc *lotteryController) GetPrize() string{
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	code := r.Intn(10000)
	//中奖信息
	var mPrize string
	var prizeRate *Prate
	//从奖品列表匹配是否中奖
	for i,prize := range prizeList{
		if code <= rateList[i].CodeB && code >= rateList[i].CodeA{
			mPrize = prize
			prizeRate = &rateList[i]
			break
		}
	}
	//开始发奖
	if mPrize == ""{
		return fmt.Sprintf("很遗憾，您没有中奖")
	}else if prizeRate.Total >0{
		prizeRate.Left -= 1
		return fmt.Sprintf("恭喜您获得:%v",mPrize)
	}else {
		return fmt.Sprintf("很遗憾，您没有中奖")
	}
}


package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

//支付宝集福卡


type gift struct {
	id       int //奖品id
	name     string
	pic      string
	link     string   //奖品链接
	inuse    bool     //是否在抽奖
	rate     int      //抽奖概率， 万分之rate  0-9999
	rateMin  int      //最小中奖编码
	rateMax  int      //最大中奖编码
}
var logger *log.Logger
const rateMax = 10
func initLog() {
	f, _ := os.Create("lottery.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}
type lotteryController struct {
	Ctx iris.Context
}
func newApp() *iris.Application{
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	initGift()
	initLog()
	return app
}
func initGift() *[5]gift {
	gifts := new([5]gift)
	return gifts
}
// 初始化奖品列表信息（管理后台来维护）
func newGift() *[5]gift {
	giftlist := new([5]gift)
	// 1 实物大奖
	g1 := gift{
		id:      1,
		name:    "富强福",
		pic:     "富强福.jpg",
		link:    "",
		inuse:   true,
		rate:    4,
		rateMin: 0,
		rateMax: 0,
	}
	giftlist[0] = g1
	// 2 实物小奖
	g2 := gift{
		id:      2,
		name:    "和谐福",
		pic:     "和谐福.jpg",
		link:    "",
		inuse:   true,
		rate:    3,
		rateMin: 0,
		rateMax: 0,
	}
	giftlist[1] = g2
	// 3 虚拟券，相同的编码
	g3 := gift{
		id:      3,
		name:    "友善福",
		pic:     "友善福.jpg",
		link:    "",
		inuse:   true,
		rate:    2,
		rateMin: 0,
		rateMax: 0,
	}
	giftlist[2] = g3
	// 4 虚拟券，不相同的编码
	g4 := gift{
		id:      4,
		name:    "爱国福",
		pic:     "爱国福.jpg",
		link:    "",
		inuse:   true,
		rate:    1,
		rateMin: 0,
		rateMax: 0,
	}
	giftlist[3] = g4
	// 5 虚拟币
	g5 := gift{
		id:      5,
		name:    "敬业福",
		pic:     "敬业福.jpg",
		link:    "",
		inuse:   true,
		rate:    0,
		rateMin: 0,
		rateMax: 0,
	}
	giftlist[4] = g5
	return giftlist
}

//导入每个福卡的概率
func giftRate(rate string) *[5]gift{
	giftList := newGift()
	rateList := strings.Split(rate,",")
	ratesLen := len(rateList)
	rateStart := 0
	//设置每个副卡的中奖概率
	for i ,data := range giftList{
		if !data.inuse{
			continue
		}
		grate := 0
		if i<ratesLen{
			grate ,_= strconv.Atoi(rateList[i])
		}
		data.rate = grate
		data.rateMin = rateStart
		data.rateMax = rateStart + grate
		fmt.Printf("data.rate = %d \n data.rateMin = %d\n data = rateMax = %d \n",data.rate,data.rateMin,data.rateMax)
		if giftList[i].rateMax >= rateMax{
			data.rateMax = rateMax
			rateStart = 0
		}else{
			rateStart = rateStart+grate
		}
		giftList[i] = data
	}
	fmt.Printf("giftList=%v",giftList)
	return giftList
}

func (lc *lotteryController) Get() string{
	rate := lc.Ctx.URLParamDefault("rate","4,3,2,1,0")
	giftlistt := giftRate(rate)
	return fmt.Sprintf("%v\n",giftlistt)
}
func (lc *lotteryController) GetLucky() map[string]interface{}{
	uid ,_:= lc.Ctx.URLParamInt("uid")
	rate := lc.Ctx.URLParamDefault("rate","4,3,2,1,0")
	code := luckyCode()
	ok := false
	result := make(map[string]interface{})
	result["success"] = ok
	giftList := giftRate(rate)
	for _,data:= range giftList{
		if !data.inuse{
			continue
		}
		//中奖
		if  int(code) <=data.rateMax&& int (code) >=data.rateMin {
			ok = true
			sendData := data.pic
			if ok {
				saveLuckyData(code, data.id, data.name, data.link, sendData)
				result["success"] = ok
				result["uid"] = uid
				result["id"] = data.id
				result["name"] = data.name
				result["link"] = data.link
				result["data"] = sendData
				break
			}
		}
	}
	return result
}
func saveLuckyData(code int32, id int, name string, link string, senddata string) {
	logger.Printf("Lucky,code : %d,id : %d,name = %s,link = %s,senddata = %s",
		code,id,name,link,senddata)
}
func luckyCode() int32{
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Int31n(int32(rateMax))
	return code
}
func main(){
	app := newApp()
	app.Run(iris.Addr(":8080"))
}
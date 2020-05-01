package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

/**
微信摇一摇
只有一个抽奖的接口 /lucky
压力测试： （可发现线程问题）
wrk -t10 -c10 -d5 http://localhost:8080/lucky
t : 线程数   c : 连接数    d : 持续时间
*/
const (
	giftTypeCoin      = iota //虚拟币
	giftTypeCoupon           //不同卷
	giftTypeCouponfix        //相同卷
	giftTypeRealSmall        //实物小
	giftTypeRealLarge        //实物大
)

type gift struct {
	id       int //奖品id
	name     string
	pic      string
	link     string   //奖品链接
	gtype    int      //奖品类型
	data     string   //币的数值，券对应的代码（特定的配置信息
	datalist []string //奖品数据集合(不同的优惠券编码)
	total    int      //总数 0不限量
	left     int      //剩余数量
	inuse    bool     //是否在抽奖
	rate     int      //抽奖概率， 万分之rate  0-9999
	rateMin  int      //最小中奖编码
	rateMax  int      //最大中奖编码
}
type lotteryController struct {
	Ctx *iris.Context
}

const rateMax = 10000
var mu sync.Mutex
var logget *log.Logger

//奖品列表
var giftList []*gift
var logger *log.Logger
//初始化日志
func initLog() {
	f, _ := os.Create("lottery.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}

// 初始化奖品列表信息（管理后台来维护）
func initGift() {
	giftList = make([]*gift, 5)
	// 1 实物大奖
	g1 := gift{
		id:      1,
		name:    "手机N7",
		pic:     "",
		link:    "",
		gtype:   giftTypeRealLarge,
		data:    "",
		total:   10000,
		left:    10000,
		inuse:   true,
		rate:    10000,
		rateMin: 0,
		rateMax: 0,
	}
	giftList[0] = &g1
	// 2 实物小奖
	g2 := gift{
		id:      2,
		name:    "安全充电 黑色",
		pic:     "",
		link:    "",
		gtype:   giftTypeRealSmall,
		data:    "",
		total:   5,
		left:    5,
		inuse:   false,
		rate:    100,
		rateMin: 0,
		rateMax: 0,
	}
	giftList[1] = &g2
	// 3 虚拟券，相同的编码
	g3 := gift{
		id:      3,
		name:    "商城满2000元减50元优惠券",
		pic:     "",
		link:    "",
		gtype:   giftTypeCouponfix ,
		data:    "mall-coupon-2018",
		total:   5,
		left:    5,
		rate:    5000,
		inuse:   false,
		rateMin: 0,
		rateMax: 0,
	}
	giftList[2] = &g3
	// 4 虚拟券，不相同的编码
	g4 := gift{
		id:       4,
		name:     "商城无门槛直降50元优惠券",
		pic:      "",
		link:     "",
		gtype:    giftTypeCoupon,
		data:     "",
		datalist: []string{"c01", "c02", "c03", "c04", "c05"},
		total:    50,
		left:     50,
		inuse:    false,
		rate:     2000,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[3] = &g4
	// 5 虚拟币
	g5 := gift{
		id:      5,
		name:    "社区10个金币",
		pic:     "",
		link:    "",
		gtype:   giftTypeCoin,
		data:    "10",
		total:   5,
		left:    5,
		inuse:   false,
		rate:    5000,
		rateMin: 0,
		rateMax: 0,
	}
	giftList[4] = &g5

	// 整理奖品数据，把rateMin,rateMax根据rate进行编排
	rateStart := 0
	for _, data := range giftList {
		if !data.inuse {
			continue
		}
		data.rateMin = rateStart
		data.rateMax = data.rateMin + data.rate
		if data.rateMax >= rateMax {
			// 号码达到最大值，分配的范围重头再来
			data.rateMax = rateMax
			rateStart = 0
		} else {
			rateStart += data.rate
		}
	}
	fmt.Printf("giftlist=%v\n", giftList)
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	initLog()
	initGift()
	return app
}
func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))

}
//奖品数量信息 Get http://localhost:8080/
func (lc *lotteryController) Get() string{
	var count int //有效奖品数量
	var total int //限量奖品总数量
	for _,data := range giftList{
		if data.inuse && (data.total ==0 || (data.total>0&&data.left>0)){
			count++
			total += data.left
		}
	}
	return fmt.Sprintf("当前有效商品总数量：%d,当前限量奖品总数量：%d\n",count,total )
}
func (lc *lotteryController) GetLucky() map[string]interface{}{
	//每个人获取一个随机码,抽中giftlist中的一个奖品之后退出，抽一个或全都抽不到
	code := luckycode()
	ok := false
	result := make(map[string]interface{})
	result["success"] = ok
	for _,data := range giftList{
		//奖品没有抽奖或者库存不足
		if !data.inuse || (data.total>0 && data.left<=0){
			continue
		}
		if data.rateMin<=int(code)&&data.rateMax>=int(code){
			//中奖了，抽奖编码在奖品编码范围内
			//开始发奖
			sendData := ""
			//判断奖品类型，做不同的发奖
			switch data.gtype {
			case giftTypeCoin:
				ok,sendData=sendCoin(data)
			case giftTypeCoupon:
				ok,sendData=sendCoupon(data)
			case giftTypeCouponfix:
				ok,sendData=sendCouponFix(data)
			case giftTypeRealLarge:
				ok,sendData=sendRealLarge(data)
			case giftTypeRealSmall:
				ok,sendData=sendRealSmall(data)
			}
			if ok{
				//中奖后，保存中奖记录
				saveLuckyData(code,data.id,data.name,data.link,sendData,data.left)
				result["success"] = ok
				result["id"] = data.id
				result["name"] = data.name
				result["link"] = data.link
				result["data"] = sendData
				result["left"] = data.left
				break
			}
		}
	}

	return result
}

func saveLuckyData(code int32, id int, name string, link string, senddata string, left int) {
	logger.Printf("Lucky,code : %d,id : %d,name = %s,link = %s,senddata = %s,left = %d",
		code,id,name,link,senddata,left)
}

func sendCoin(data *gift) (bool,string) {
	//coin数量无限
	if data.total == 0{
		return true,data.data
	} else{
		//奖品数量不足
		if data.left <= 0 {
			return false, "奖品数量不足"
		} else  {
			//奖品数量足够
			data.left -= 1
			return true, data.data
		}
	}
}
//不同值的优惠券
func sendCoupon(data *gift) (bool,string){
	//datalist长度小于left长度，直接抽奖失败
	if len(data.datalist) < data.left{
		return false,"抽奖失败"
	}
	//有库存,每次发的是倒数最后一个优惠券
	if data.left >0{
		left := data.left-1
		data.left = left
		return true,data.datalist[left]
	}else {
		return false ,"奖品数量不足"
	}
}
//发放同值的优惠券
func sendCouponFix(data *gift) (bool,string){
	//数量无限
	if data.total == 0{
		return true,data.data
	}
	//数量足够
	if data.left>0{
		data.left -=1
		return true,data.data
	}else{
		//数量不足
		return false,"奖品数量不足"
	}
}
func sendRealLarge (data *gift) (bool,string){
	//数量无限
	if data.total == 0{
		return true,data.data
	}
	//数量足够
	if data.left>0{
		mu.Lock()
		defer mu.Unlock()
		data.left -=1
		return true,data.data
	}else{
		//数量不足
		return false,"奖品数量不足"
	}
}
func sendRealSmall(data *gift) (bool,string){
	//数量无限
	if data.total == 0{
		return true,data.data
	}
	//数量足够
	if data.left>0{
		data.left -=1
		return true,data.data
	}else{
		//数量不足
		return false,"奖品数量不足"
	}
}
func luckycode() int32{
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Int31n(int32(rateMax))
	return code
}

package main

/**
设置红包
curl "http://localhost:8080/set?uid=1&money=10000&num=100000"
抢红包
curl "http://localhost:8080/get?uid=1&id=3071854753"
并发压力测试
wrk -t10 -c10 -d5 "http://localhost:8080/get?uid=1&id=3071854753"
*/
//存在两个同步方面的问题
//packageList的同步问题以及list的同步问题
//使用sync.list解决packageList的同步问题，用chan 解决list的同步问题
import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"sync"
	"time"
)
type task struct {
	id uint32
	//数据更新完成之后  回调
	callback chan uint
}
//var packageList map[uint32][]uint = make(map[uint32][]uint)
var packageList *sync.Map = new(sync.Map)
var chTasks chan task = make(chan task)
type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	go fetchPackageListMoney()
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))

}

//查看每个红包红包总量以及总金额数
func (lc *lotteryController) Get() map[uint32][2]int {
	rs := make(map[uint32][2]int)
	packageList.Range(func(key, value interface{}) bool {
		id := key.(uint32)
		list := value.([]uint)
		var money int
		for _, v := range list {
			money += int(v)
		}
		rs[id] = [2]int{len(list), money}
		return true
	})
	return rs
}

//发红包
//http://localhost:8080/get?id=1&money=100&num=100
func (lc *lotteryController) GetSet() string {
	uid, uidErr := lc.Ctx.URLParamInt("uid")
	money, moneyErr := lc.Ctx.URLParamFloat64("money")
	num, numErr := lc.Ctx.URLParamInt("num")
	if uidErr != nil || moneyErr != nil || numErr != nil {
		return fmt.Sprintf("uidErr=%v,moneyErr=%v,numErr=%v\n", uidErr, moneyErr, numErr)
	}
	moneyTotal := int(money * 100)
	if uid < 1 || moneyTotal < num || num < 1 {
		return fmt.Sprintf("参数异常，uid=%v,money=%d,num=%d", uid, money, num)
	}
	//金额分配算法
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rateMax := 0.55
	if num > 1000 {
		rateMax = 0.01
	} else if num >= 100 {
		rateMax = 0.1
	} else if num >= 10 {
		rateMax = 0.3
	}
	list := make([]uint, num)
	leftMoney := moneyTotal
	leftNum := num
	//分配金额到每个红包
	for leftMoney > 0 {
		//分配剩余全部金额
		if leftNum == 1 {
			list[num-1] = uint(leftMoney)
			break
		}
		//红包不可分拆
		if leftNum == leftMoney {
			for i := num - leftNum; i < num; i++ {
				list[i] = 1
			}
			break
		}
		rmoneyMax := int(float64(leftMoney-leftNum) * rateMax)
		rmoney := r.Intn(rmoneyMax)
		if rmoney < 1 {
			rmoney = 1
		}
		list[num-leftNum] = uint(rmoney)
		leftMoney -= rmoney
		leftNum -= 1
	}
	//红包的唯一id
	id := r.Uint32()
	//packageList[id] = list
	packageList.Store(id,list)
	return fmt.Sprintf(
		"/get?id=%d&uid=%d&num=%d", id, uid, num)
}

//抢红包
//http://localhost:8080/get?id=1&uid=1
func (lc *lotteryController) GetGet() string {
	id, idErr := lc.Ctx.URLParamInt("id")
	uid, uidErr := lc.Ctx.URLParamInt("uid")
	if idErr != nil || uidErr != nil {
		return fmt.Sprintf("idErr=%v,uidErr=%v", idErr, uidErr)
	}
	if id < 1 || uid < 1 {
		return fmt.Sprintf("操作异常")
	}
	//list, ok := packageList[uint32(id)]
	listOrigin,ok:=packageList.Load(uint32(id))
	if !ok{
		return fmt.Sprintf("红包不存在")
	}
	list := listOrigin.([]uint)
	if len(list) < 1 {
		return fmt.Sprintf("红包不存在,id=%d", id)
	}
	callback :=make(chan uint)
	//发送任务
	t := task{id :uint32(id),callback: callback}
	//发送任务
	chTasks <- t
	//接收返回结果
	money := <- callback
	if money <= 0{
		return fmt.Sprintf("很遗憾，没有抢到红包")
	}else {
		return fmt.Sprintf("恭喜您抢到一个红包，金额为%d\n", money)
	}
}
//处理抢红包任务
func fetchPackageListMoney(){
	for{
		t:=<-chTasks
		id := t.id
		l,ok :=packageList.Load(id)
		if ok && l !=nil{
			list := l.([]uint)
			//分配随机数
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			i := r.Intn(len(list))
			money := list[i]
			//取出钱之后更新红包数组信息
			if len(list) > 1 {
				if i == len(list)-1 {
					//packageList[uint32(id)] = list[:i]
					packageList.Store(uint32(id),list[:i])
				} else if i == 0 {
					//packageList[uint32(id)] = list[1:]
					packageList.Store(uint32(id),list[1:])
				} else {
					//packageList[uint32(id)] = append(list[:i], list[i+1:]...)
					packageList.Store(uint32(id),append(list[:i], list[i+1:]...))
				}
			} else {
				//delete(packageList, uint32(id))
				packageList.Delete(uint32(id))
			}
			//抢到红包，返回金额
			t.callback <- money
		}else {
			//没抢到红包，返回0
			t.callback <- 0
		}
	}
}


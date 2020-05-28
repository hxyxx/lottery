package controller

import (
	"fmt"
	"hxyxx/lottery/conf"
	"hxyxx/lottery/models"
	"log"
	"strconv"
	"time"
)

func (c IndexController) checkUserday(uid int) bool{
	userdayInfo :=c.ServiceUserDay.GetUserToday(uid)
	if userdayInfo != nil && userdayInfo.Uid == uid{
		if userdayInfo.Num >= conf.UserPrizeMax{
			return false
		}else {
			userdayInfo.Num++
			err103 := c.ServiceUserDay.Update(userdayInfo,nil)
			if err103 != nil {
				log.Println("index_lucky_check_userday ServiceUserDay.Update "+ "err103=",err103)
			}
		}
	}else {
		//创建今天的抽奖记录
		 y,m,d := time.Now().Date()
		 strDay := fmt.Sprintf("%d%02d%02d",y,m,d)
		 day,_ :=strconv.Atoi(strDay)
		 userdayInfo = &models.LtUserday{
		 	Uid: uid,
		 	Day: day,
		 	Num: 1,
		 	SysCreated: int(time.Now().Unix()),
		 }
		 err103 :=c.ServiceUserDay.Create(userdayInfo)
		if err103 != nil {
			log.Println("index_lucky_check_userday ServiceUserDay.Update "+ "err103=",err103)
		}

	}
	return true
}

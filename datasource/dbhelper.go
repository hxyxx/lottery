package datasource

import (
	"fmt"
	"github.com/go-xorm/xorm"
	//不init 会报错
	_ "github.com/go-sql-driver/mysql"
	"hxyxx/lottery/conf"
	"log"
	"sync"
)

var masterInstance *xorm.Engine
var mu sync.Mutex

//单例,并发可能创建多个engine，加锁
func InstanceDbMaster() *xorm.Engine {
	if masterInstance != nil {
		return masterInstance
	}
	mu.Lock()
	defer mu.Unlock()
	//可能并发连接多个被锁住，所以重新进行判断
	if masterInstance != nil{
		return masterInstance
	}
	return NewDbHelper()

}
func NewDbHelper() *xorm.Engine {
	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		conf.DbMaster.User,
		conf.DbMaster.Pwd,
		conf.DbMaster.Host,
		conf.DbMaster.Port,
		conf.DbMaster.Database)
	instance, err := xorm.NewEngine(conf.DriverName, sourceName)
	if err != nil {
		log.Fatal("dbhelper.NewDbMaster NewEngine error", err)
		return nil
	}
	instance.ShowSQL(true) //显示sql语句，调试的时候使用
	masterInstance = instance
	return instance
}

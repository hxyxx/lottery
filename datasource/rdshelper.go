package datasource

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"hxyxx/lottery/conf"
	"log"
	"sync"
	"time"
)

var rdsLock sync.Mutex
var cacheInstance *RedisConn

type  RedisConn struct{
	pool *redis.Pool
	showDebug bool
}

func (rds *RedisConn) Do(commandName string,args ...interface{}) (reply interface{},err error){
	conn :=rds.pool.Get()
	defer conn.Close()
	t1:=time.Now().UnixNano()
	reply ,err = conn.Do(commandName,args)
	if err !=nil{
		e :=conn.Err()
		if e !=nil{
			log.Println("rdshelper.Do",err,e)
		}
	}
	t2 := time.Now().UnixNano()
	if rds.showDebug == true{
		fmt.Printf("[redis] [info] [%dus] cmd=%s,err=%v,reply%s\n",
			(t2-t1)/1000,commandName,err,args,reply)
	}
	return reply,err
}

func (rds *RedisConn) ShowDebug(b bool){
	rds.showDebug = b
}

func InstanceCache() *RedisConn{
	if cacheInstance != nil {
		return cacheInstance
	}
	rdsLock.Lock()
	defer rdsLock.Unlock()
	if cacheInstance != nil{
		return cacheInstance
	}
	return NewCache()
}

func NewCache() *RedisConn{
	pool := redis.Pool{
		Dial: func() (redis.Conn,error){
			c,err := redis.Dial("tcp",fmt.Sprintf("%s:%d",conf.RdsCache.Host,conf.RdsCache.Port))
			if err != nil{
				log.Fatal("rdshelper.NewCache Dial error = ",err)
				return nil,err
			}
			return c,nil
		},
		TestOnBorrow: func(c redis.Conn,t time.Time) error{
			if time.Since(t) < time.Minute{
				return nil
			}
			_,err := c.Do("PING")
			return err
		},
		MaxIdle: 10000,//最多连接数
		MaxActive: 10000,//最多活跃数
		IdleTimeout: 0,//超时时间
		Wait: false,
		MaxConnLifetime: 0, //活跃时间，一直活跃
	}
	instance := &RedisConn{
		pool: &pool,
	}
	cacheInstance = instance
	instance.showDebug = true
	return instance
}
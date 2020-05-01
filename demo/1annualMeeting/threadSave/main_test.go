package main
//线程不安全
import (
	"fmt"
	"github.com/kataras/iris/httptest"
	"sync"
	"testing"
)

func TestMVC(t *testing.T){
	e := httptest.New(t,newApp())
	//协程
	var wg sync.WaitGroup
	//期待的状态值为200，期待的返回值为 抽奖总数为0
	e.GET("/").Expect().Status(httptest.StatusOK).
		Body().Equal("当前抽奖总人数为： 0\n")

	//开启协程，导入用户
	for i := 0;i<100;i++{
		wg.Add(1)
		go func(i int){
			e.POST("/import").WithFormField("users",fmt.Sprintf("user_u%d",i)).
				Expect().Status(httptest.StatusOK)
			wg.Done()
		}(i)
	}
	wg.Wait()
	//查看导入用户数量
	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("当前抽奖总人数为： 100\n")
	//抽奖
	e.GET("/lucky").Expect().Status(httptest.StatusOK)

	//查看导入用户
	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("当前抽奖总人数为： 99\n")

}
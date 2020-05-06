package bootstrap

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	recover2 "github.com/kataras/iris/middleware/recover"
	"hxyxx/lottery/conf"
	"time"
)

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	*iris.Application//继承
	AppName string
	AppOwner string
	AppSpawDate time.Time
}

func New(appName,appOwner string,cfgs ...Configurator) *Bootstrapper{
	b := &Bootstrapper{
		Application :iris.New(),
		AppName:appName,
		AppOwner: appOwner,
		AppSpawDate: time.Now(),
	}
	for _,cfg := range cfgs{
		cfg(b)
	}
	return b
}
//
func (b *Bootstrapper) SetupViews(viewDir string){

	htmlEngine := iris.HTML(viewDir,".html").Layout("shared/layout.html")
	htmlEngine.Reload(true)//开发模式设置为true
	htmlEngine.AddFunc("FromUnixtimeShort",func (t int )string{//当前时间
		dt := time.Unix(int64(t),int64(0))
		return dt.Format(conf.SysTimeFormShort)
	})
	htmlEngine.AddFunc("FromUnixtime",func (t int )string{//当前时间
		dt := time.Unix(int64(t),int64(0))
		return dt.Format(conf.SysTimeForm)
	})
	b.RegisterView(htmlEngine)
}

//异常处理
func (b *Bootstrapper) SetupErrorHandler(){
	b.OnAnyErrorCode(func (ctx iris.Context){
		err := iris.Map{
			"app" : b.AppName,
			"status" : ctx.GetStatusCode(),
		  	"message" : ctx.Values().GetString("message"),
		}
		//json输出
		if jsonOutput := ctx.URLParamExists("json");jsonOutput{
			ctx.JSON(err)
			return
		}
		//使用模板输出
		ctx.ViewData("err",err)
		ctx.ViewData("title","Error")
		//设置模板输出
		ctx.View("shared/error.html")
	})
}

//配置
func (b *Bootstrapper) Configure(cs ...Configurator){
	for _,c := range cs{
		c(b)
	}
}

//计划任务程序
func (b *Bootstrapper) setUpCron(){
	// todo:
}

const  (
	StaticAssets = "./public/"
	Favicon = "favicon.ico"
)

//初始化bootstrapper
func (b *Bootstrapper) Bootstrap() *Bootstrapper{
	b.SetupViews("./views")
	b.SetupErrorHandler()
	b.Favicon(StaticAssets+Favicon)
	//b.StaticWeb(StaticAssets[1:len(StaticAssets)-1],StaticAssets)

	b.setUpCron()
	b.Use(recover2.New())
	b.Use(logger.New())
	return b
}

//提供一个监听的方法
func (b *Bootstrapper) Listen(addr string,cfgs ...iris.Configurator){
	b.Run(iris.Addr(addr),cfgs...)
}
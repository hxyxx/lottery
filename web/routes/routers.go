package routes

import (
	"github.com/kataras/iris/mvc"
	"hxyxx/lottery/bootstrap"
	"hxyxx/lottery/services"
	"hxyxx/lottery/web/controller"
)

func Configure(b *bootstrap.Bootstrapper){
	userService := services.NewUserService()
	giftService := services.NewGiftService()
	codeService := services.NewCodeService()
	resultService := services.NewResultService()
	blackipService := services.NewBlackipService()
	index := mvc.New(b.Party("/"))
	index.Register(userService,giftService,codeService,resultService,blackipService)
	index.Handle(new(controller.IndexController))
}
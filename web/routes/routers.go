package routes

import (
	"github.com/kataras/iris/_examples/mvc/login/web/middleware"
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


	admin := mvc.New(b.Party("/admin"))
	admin.Router.Use(middleware.BasicAuth)
	admin.Register(userService,giftService,codeService,resultService,blackipService)
	admin.Handle(new(controller.AdminController))

	adminGift :=admin.Party("/gift")
	adminGift.Register(giftService)
	adminGift.Handle(new(controller.AdminGiftController))

	adminCode := admin.Party("/code")
	adminCode.Register(codeService)
	adminCode.Handle(new(controller.AdminCodeController))

	adminResult := admin.Party("/result")
	adminResult.Register(resultService)
	adminResult.Handle(new(controller.AdminResultController))

	adminBlackip := admin.Party("/user")
	adminBlackip.Register(userService)
	adminBlackip.Handle(new(controller.AdminUserController))


}
package controller

import "github.com/kataras/iris"

type IndexController struct{
	Ctx iris.Context
	ServiceUser
}

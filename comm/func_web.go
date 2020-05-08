package comm

import (
	"crypto/md5"
	"fmt"
	"hxyxx/lottery/conf"
	"hxyxx/lottery/models"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
)
//获取ip地址
func ClientIP(request *http.Request) string {
	host ,_,_:=net.SplitHostPort(request.RemoteAddr)
	return host
}
//跳转
func Redirect(writer http.ResponseWriter,url string){
	writer.Header().Add("Location",url)
	writer.WriteHeader(http.StatusFound)
}
//
func GetLoginUser(request *http.Request) *models.ObjLoginuser{
	c,err := request.Cookie("lottery_loginuser")
	if err != nil{
		return nil
	}
	//Cookie的值是url类型的
	params ,err := url.ParseQuery(c.Value)
	if err != nil{
		return nil
	}
	uid,err := strconv.Atoi(params.Get("uid"))
	if err !=nil || uid<1{
		return nil
	}
	now,err:=strconv.Atoi(params.Get("now"))
	//如果cookie中的时间超过了30天
	if err != nil||(NowUnix()-now)>86400*30{
		return nil
	}
	loginuser := &models.ObjLoginuser{}
	loginuser.Uid = uid
	loginuser.Username = params.Get("username")
	loginuser.Now = now
	loginuser.Ip = ClientIP(request)
	loginuser.Sign = params.Get("sign")
	//对签名进行验证
	sign := createLoginuserSign(loginuser)
	if sign != loginuser.Sign{
		log.Println("func_web GetLoginuser createloginusersign not sign",sign,loginuser.Sign)
		return nil
	}
	return loginuser
}
//设置cookie
func SetLoginuser(writter http.ResponseWriter,loginuser *models.ObjLoginuser){
	if loginuser == nil || loginuser.Uid <1{
		c := &http.Cookie{
			Name : "lottery_loginuser",
			Value : "",
			Path : "/",
			MaxAge : -1,
		}
		http.SetCookie(writter,c)
		return
	}
	if loginuser.Sign ==""{
		loginuser.Sign = createLoginuserSign(loginuser)
	}
	params := url.Values{}
	params.Add("uid",strconv.Itoa(loginuser.Uid))
	params.Add("username",loginuser.Username)
	params.Add("now",strconv.Itoa(loginuser.Now))
	params.Add("ip",loginuser.Ip)
	params.Add("sign",loginuser.Sign)
	c := &http.Cookie{
		Name: "lottery_loginuser",
		Value: params.Encode(),
		Path: "/",
	}
	http.SetCookie(writter,c)
}
//一致性签名算法
func createLoginuserSign(loginuser *models.ObjLoginuser) string{
	str := fmt.Sprintf("uid=%s&username=%s&secret=%s&now=%d",
		loginuser.Uid,loginuser.Username,conf.CookieSecret,loginuser.Now)
	sign := fmt.Sprintf("%x",md5.Sum([]byte(str)))
	return sign
}
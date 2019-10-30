package controller

import (
	"SecondKill/SecProxy/service"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
	"time"
)

type SkillController struct {
	beego.Controller
}

func (p *SkillController) SecKill() {
	productId, err := p.GetInt("product_id")
	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "success"

	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if err != nil {
		result["code"] = 1001
		result["message"] = "invalid product_id"
		return
	}

	source := p.GetString("src")
	authCode := p.GetString("authcode")
	secTime := p.GetString("time")
	nance := p.GetString("nance")

	secRequest := service.NewSecRequest()
	secRequest.AuthCode = authCode
	secRequest.Source = source
	secRequest.SecTime = secTime
	secRequest.Nance = nance
	secRequest.ProductId = productId
	secRequest.UserAuthSign = p.Ctx.GetCookie("userAuthSign")
	secRequest.UserId, _ = strconv.Atoi(p.Ctx.GetCookie("UserId"))
	secRequest.AccessTime = time.Now()

	if len(p.Ctx.Request.RemoteAddr) > 0 {
		secRequest.ClientAddr = strings.Split(p.Ctx.Request.RemoteAddr, ":")[0]
	}
	secRequest.ClientRefence = p.Ctx.Request.Referer()
	secRequest.CloseNotify = p.Ctx.ResponseWriter.CloseNotify()
	logs.Debug("client request:[%v]", secRequest)

	data, code, err := service.SecKill(secRequest)
	if err != nil {
		result["code"] = code
		result["message"] = err.Error()
		return
	}
	result["data"] = data
	result["code"] = code

	p.ServeJSON()
}

func (p *SkillController) SecInfo() {
	productId, err := p.GetInt("product_id")
	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "success"

	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if err != nil {
		data, code, err := service.SecInfoList()
		if err != nil {
			result["code"] = code
			result["messsage"] = "invalid productId"
			logs.Error("invaild request, get product id failed, err:%v", err)
			return
		}
		result["code"] = code
		result["message"] = data

	} else {
		data, code, err := service.SecInfo(productId)
		if err != nil {
			result["code"] = code
			result["message"] = err.Error()
			logs.Error("invalid request, get product id failed, err:%v", err)
			return
		}
		result["code"] = code
		result["data"] = data
	}

}

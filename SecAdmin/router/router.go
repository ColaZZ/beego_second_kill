package router

import (
	"SecondKill/SecAdmin/controller/product"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("product/list", &product.ProductController{}, "*:ListProduct")
	beego.Router("/", &product.ProductController{}, "*:ListProduct")
	beego.Router("/product/create", &product.ProductController{}, "*:CreateProduct")
	beego.Router("/product/submit", &product.ProductController{}, "*:SubmitProduct")
}

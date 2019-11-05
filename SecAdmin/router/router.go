package router

import (
	"SecondKill/SecAdmin/controller/product"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("product/list", &product.ProductController{}, "*:ListProduct")
}

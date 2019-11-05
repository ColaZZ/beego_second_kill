package product

import (
	"github.com/astaxie/beego"
)

type ProductController struct {
	beego.Controller
}

func (p *ProductController) ListProduct() {
	//productModel := model.NewProductModel()
	//productList, err := productModel.GetProductList()
	//if err != nil {
	//	logs.Warn("get product list failed, err:%v", err)
	//	return
	//}
	//p.Data["product_list"] = productList
	p.TplName = "product/list.html"
	p.Layout = "layout/layout.html"
}

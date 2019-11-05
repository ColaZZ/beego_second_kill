package product

import (
	"SecondKill/SecAdmin/model"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type ProductController struct {
	beego.Controller
}

func (p *ProductController) ListProduct() {
	productModel := model.NewProductModel(model.Db)
	productList, err := productModel.GetProductList()
	if err != nil {
		logs.Warn("get product list failed, err:%v", err)
		return
	}
	p.Data["product_list"] = productList
	p.TplName = "product/list.html"
	p.Layout = "layout/layout.html"
}

func (p *ProductController) CreateProduct() {
	p.TplName = "product/list.html"
	p.Layout = "layout/layout.html"
}

func (p *ProductController) SubmitProduct() {
	productName := p.GetString("product_name")
	productTotal, err := p.GetInt("product_total")

	p.TplName = "product/create.html"
	p.Layout = "layout/layout.html"
	errMsg := "success"

	defer func() {
		if err != nil {
			p.Data["Error"] = errMsg
			p.TplName = "product/error.html"
			p.Layout = "layout/layout.html"
		}
	}()

	if len(productName) == 0 {
		logs.Warn("invalid product name")
		errMsg = fmt.Sprintf("invalid product name")
		err = errors.New("invalid product name")
		return
	}

	if err != nil {
		logs.Warn("invalid product total, err:%v", err)
		errMsg = fmt.Sprintf("invalid product total, err:%v", err)
		return
	}

	productStatus, err := p.GetInt("product_status")
	if err != nil {
		logs.Warn("invalid product status, err:%v", err)
		errMsg = fmt.Sprintf("invalid product status, err:%v", err)
		return
	}

	//model操作
	productModel := model.NewProductModel(model.Db)
	product := model.Product{
		ProductName: productName,
		Total:       productTotal,
		Status:      productStatus,
	}
	// insert
	err = productModel.CreateProduct(&product)
	if err != nil {
		logs.Warn("create product failed, err:%v", err)
		errMsg = fmt.Sprintf("create product failed, err:%v", err)
		return
	}
	logs.Debug("product name[%s], product total[%d], product status[%d]", productName, productTotal, productStatus)
}

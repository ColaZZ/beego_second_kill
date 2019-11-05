package model

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ProductModel struct {
	DB *sqlx.DB
}

type Product struct {
	ProductId   int    `db:"id"`
	ProductName string `db:"name"`
	Total       int    `db:"total"`
	Status      int    `db:"status"`
}

func NewProductModel(db *sqlx.DB) (productModle *ProductModel) {
	productModle = &ProductModel{
		DB: db,
	}
	return
}

// 返回所有商品列表
func (p *ProductModel) GetProductList() (list []*Product, err error) {
	sqlStr := "select id,name,total,status from product"
	err = Db.Select(&list, sqlStr)
	if err != nil {
		logs.Warn("select from mysql failed, err:%v", err)
		return
	}
	return
}

// 插入新建的商品信息
func (p *ProductModel) CreateProduct(product *Product) (err error) {
	sqlSrt := "insert into product(name,total,status) values(?,?,?)"
	_, err = Db.Exec(sqlSrt, product.ProductName, product.Total, product.Status)
	if err != nil {
		logs.Warn("insert into mysql failed, err:%v", err)
		return
	}
	logs.Debug("insert into mysql success")
	return
}

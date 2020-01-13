// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/jinzhu/gorm"
	"rest-gin-gorm/invoice"
	"rest-gin-gorm/product"
)

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Injectors from wire.go:

func InitProductAPI(db *gorm.DB) product.ProductAPI {
	productRepository := product.ProvideProductRepostiory(db)
	productService := product.ProvideProductService(productRepository)
	productAPI := product.ProvideProductAPI(productService)
	return productAPI
}

func InitInvoiceAPI(db *gorm.DB) invoice.InvoiceAPI {
	invoiceRepository := invoice.ProvideInvoiceRepostiory(db)
	invoiceService := invoice.ProvideInvoiceService(invoiceRepository)
	invoiceAPI := invoice.ProvideInvoiceAPI(invoiceService)
	return invoiceAPI
}

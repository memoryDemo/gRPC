package service

import (
	"context"
)

var ProductService = &productService{}

type productService struct {
}

func (p *productService) mustEmbedUnimplementedProdServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (p *productService) GetProductStock(context context.Context, request *ProductRequest) (*ProductResponse, error) {
	// 实现具体的业务逻辑

	stock := p.GetStockById(request.ProdId)

	return &ProductResponse{ProdStock: stock}, nil
}

func (p *productService) GetStockById(id int32) int32 {
	return id
}

package service

import (
	"context"
	"google.golang.org/protobuf/types/known/anypb"
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
	user := User{Username: "memory"}
	content := Content{Msg: "memory msg..."}
	a, _ := anypb.New(&content)
	return &ProductResponse{ProdStock: stock, User: &user, Data: a}, nil
}

func (p *productService) GetStockById(id int32) int32 {
	return id
}

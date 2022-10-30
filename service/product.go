package service

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
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

func (p *productService) UpdateProductStockClientStream(stream ProdService_UpdateProductStockClientStreamServer) error {
	count := 0
	for {
		//源源不断的去接收客户端发来的消息
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Println("服务端接收到的流", recv.ProdId)
		count++
		if count > 10 {
			rsp := &ProductResponse{ProdStock: recv.ProdId}
			err := stream.SendAndClose(rsp)
			if err != nil {
				return err
			}
			return nil
		}
	}
}

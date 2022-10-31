package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/memory-grpc/client/auth"
	"github.com/memory-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

func main() {

	//################ 添加证书  单向认证 #################
	//creds, err2 := credentials.NewClientTLSFromFile("cert/server.pem", "*.memory.com")
	//if err2 != nil {
	//	log.Fatal("证书 错误", err2)
	//}
	//################ 添加证书  单向认证 #################

	// 添加证书-双向认证
	// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	cert, _ := tls.LoadX509KeyPair("cert/client.pem", "cert/client.key")
	// 创建一个新的，空的 CertPool
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.crt")
	// 尝试解析所传入的 PEM 编码的证书。如果成功解析会将其加到 CertPool 中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)
	// 构建基于 TLS 的 TransportCredentials 选项
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		ServerName:   "*.memory.com",
		RootCAs:      certPool,
	})

	// 1. 新建连接，端口是服务端开放的8002端口
	// 没有证书会报错 --  grpc.WithTransportCredentials(insecure.NewCredentials())
	// 无认 证，grpc  http 2

	// Token 认证
	token := &auth.Authentication{
		User:     "admin",
		Password: "admin",
	}
	//  grpc.WithPerRPCCredentials(token)  token-认证
	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(token))
	if err != nil {
		log.Fatal("服务端出错，链接不上", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("连接关闭失败！", err)
		}
	}(conn)

	// 2. 调用Product.pb.go中的NewProdServiceClient方法
	prodServiceClient := service.NewProdServiceClient(conn)

	// 3. 像调用本地方法一样调用GetProductStock方法
	//resp, err := prodServiceClient.GetProductStock(context.Background(), &service.ProductRequest{ProdId: 334455})
	//if err != nil {
	//	log.Fatal("调用gRPC方法错误：", err)
	//}
	//
	//fmt.Println("调用gRPC方法成功，ProdStock = ", resp.ProdStock, resp.User, resp.Data)

	//4. 客户端流
	//stream, err := prodServiceClient.UpdateProductStockClientStream(context.Background())
	//if err != nil {
	//	log.Fatal("获取流出错", err)
	//}
	//rsp := make(chan struct{}, 1)
	//go prodRequest(stream, rsp)
	//select {
	//case <-rsp:
	//	recv, err := stream.CloseAndRecv()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	stock := recv.ProdStock
	//	fmt.Println("客户端收到响应：", stock)
	//}

	//5. 服务端流
	request := &service.ProductRequest{ProdId: 123}
	stream, err := prodServiceClient.GetProductStockServerStream(context.Background(), request)
	if err != nil {
		log.Fatal("获取流出错", err)
	}

	for {
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端数据接收完成")
				err := stream.CloseSend()
				if err != nil {
					log.Fatal(err)
				}
				break
			}
			log.Fatal(err)
		}
		fmt.Println("客户端收到的流", recv.ProdStock)
	}
}

func prodRequest(stream service.ProdService_UpdateProductStockClientStreamClient, rsp chan struct{}) {
	count := 0
	for {
		request := &service.ProductRequest{
			ProdId: int32(rand.Intn(1000)),
		}

		err := stream.Send(request)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		count++
		if count > 10 {
			rsp <- struct{}{}
			break
		}
	}
}

package main

//
//import (
//	"fmt"
//	"google.golang.org/protobuf/proto"
//	"ms-proto/service"
//)
//
//func main() {
//	user := &service.User{
//		Username: "memory",
//		Age:      18,
//	}
//
//	//序列化的构成
//	marshal, err := proto.Marshal(user)
//	if err != nil {
//		panic(err)
//	}
//
//	//反序列化
//	newUser := &service.User{}
//	err = proto.Unmarshal(marshal, newUser)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Print(newUser.String())
//}

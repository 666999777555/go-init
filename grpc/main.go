package register

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func Register(port int, res func(s *grpc.Server)) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
		return err
	}
	s := grpc.NewServer()
	//反射接口支持查询
	reflection.Register(s)
	res(s)
	log.Printf("sever listening at %v", err)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen:%v\",err")
		return err
	}
	return err
}

package main

import (
	"encoding/json"
	"fmt"
	"context"
	"time"

	"google.golang.org/grpc"
	"github.com/satori/go.uuid"
	
	pb "audi/protobuf/message"
)

const Addr = ":3721"

type User struct {
	Id int64
	Name string
}

func NewRequest(srv string, evt string, data interface{}) (*pb.CreateMessageReq, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req := pb.CreateMessageReq{
		Srv:   srv,
		Event: evt,
		Data:  b,
	}

	return &req, nil
}

func main() {
	conn, err := grpc.Dial(Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewMessageServiceClient(conn)

	ctx , cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	req, err := NewRequest("account", "user.created", User{Id: 123, Name: "admin"})
	if err != nil {
		panic(err)
	}

	resp, err := client.CreateMessage(ctx, req)

	fmt.Printf("%#v \n", resp)
	fmt.Printf("%#v \n", err)
}
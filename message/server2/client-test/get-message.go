package main

import (
	// "encoding/json"
	"fmt"
	"context"
	"time"

	"google.golang.org/grpc"
	pb "audi/protobuf/message"
)

const Addr = ":3721"


func main() {
	conn, err := grpc.Dial(Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewMessageServiceClient(conn)

	ctx , cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	req := pb.GetMessageReq{Id: 2}

	resp, err := client.GetMessage(ctx, &req)

	fmt.Printf("%#v \n", resp)
	fmt.Printf("%#v \n", err)
	fmt.Printf("Data type %T \n", resp.Data)
	fmt.Println("Data string:", string(resp.Data))
}
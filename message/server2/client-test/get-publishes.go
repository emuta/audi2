package main

import (
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

	req := pb.PublishMessageReq{Id: 2}

	resp, err := client.GetPublishHistory(ctx, &req)

	fmt.Printf("%#v \n", resp)
	fmt.Printf("%#v \n", err)

	for _, publish := range resp.Publishes {
		fmt.Printf("%#v \n", publish)
	}
}
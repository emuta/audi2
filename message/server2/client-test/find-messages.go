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

	// req := pb.FindMessageReq{Id: 2}
	// req := pb.FindMessageReq{CreatedAtFrom: ptypes.TimestampNow()}
	req := pb.FindMessageReq{Srv: "account"}

	resp, err := client.FindMessage(ctx, &req)

	fmt.Printf("%#v \n", resp)
	fmt.Printf("%#v \n", err)

	fmt.Println()

	for _, msg := range resp.Messages {
		fmt.Println(msg)
	}
}
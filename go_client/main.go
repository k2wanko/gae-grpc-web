package main

import (
	"context"
	"log"

	"github.com/k2wanko/gae-grpc-web/echo"
	"google.golang.org/grpc"
)

func main() {
	addr := "localhost:8080"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := echo.NewEchoServiceClient(conn)

	resp, err := c.Echo(context.Background(), &echo.EchoRequest{Message: "Hello"})
	if err != nil {
		log.Fatalf("could not echo: %v", err)
	}
	log.Printf("Echo: %v", resp.GetMessage())
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/karismapa/greeter/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client is running...")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	unary(c)
	serverStream(c)
}

func unary(c greetpb.GreetServiceClient) {
	fmt.Println("Invoke Greet function...")

	ctx := context.Background()
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Karisma",
			LastName:  "Pratama",
		},
	}

	res, err := c.Greet(ctx, req)
	if err != nil {
		log.Fatalf("Error while calling Greet: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func serverStream(c greetpb.GreetServiceClient) {
	fmt.Println("Invoke GetManyTimes function...")

	ctx := context.Background()
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Karisma",
			LastName:  "Pratama",
		},
	}
	stream, err := c.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/karismapa/greeter/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client main is running...")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

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

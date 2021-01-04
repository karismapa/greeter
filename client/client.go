package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/karismapa/greeter/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client is running...")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	// unary(c)
	// serverStream(c)
	// clientStream(c)
	bidirectionalStream(c)
}

func unary(c greetpb.GreetServiceClient) {
	fmt.Println("Invoke unary function...")

	ctx := context.Background()
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Karisma",
			LastName:  "Pratama",
		},
	}

	res, err := c.Greet(ctx, req)
	if err != nil {
		log.Fatalf("Error while calling Greet: %v\n", err)
	}
	log.Printf("Response from Greet: %v\n", res.Result)
}

func serverStream(c greetpb.GreetServiceClient) {
	fmt.Println("Invoke server stream function...")

	ctx := context.Background()
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Karisma",
			LastName:  "Pratama",
		},
	}
	stream, err := c.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes: %v\n", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}
		log.Printf("Response from GreetManyTimes: %v\n", msg.GetResult())
	}
}

func clientStream(c greetpb.GreetServiceClient) {
	fmt.Println("Invoke client stream function...")

	ctx := context.Background()
	stream, err := c.LongGreet(ctx)
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v\n", err)
	}

	reqs := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Karis",
				LastName:  "Pratama",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Karisma",
				LastName:  "Pratama",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Karismapa",
				LastName:  "Pratama",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mapa",
				LastName:  "Pratama",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Kupe",
				LastName:  "Pratama",
			},
		},
	}
	for _, req := range reqs {
		fmt.Printf("Sending request: %v\n", req)
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Error while sending stream: %v\n", err)
		}

		time.Sleep(700 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v\n", err)
	}
	fmt.Printf("Response from LongGreet: %v\n", res)
}

func bidirectionalStream(c greetpb.GreetServiceClient) {
	fmt.Println("Invoke bidirectional stream function...")

	// create stream
	ctx := context.Background()
	stream, err := c.GreetAll(ctx)
	if err != nil {
		log.Fatalf("Error while calling GreetAll: %v\n", err)
		return
	}

	reqs := []*greetpb.GreetAllRequest{
		&greetpb.GreetAllRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Karis",
				LastName:  "Pratama",
			},
		},
		&greetpb.GreetAllRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Karisma",
				LastName:  "Pratama",
			},
		},
		&greetpb.GreetAllRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Karismapa",
				LastName:  "Pratama",
			},
		},
		&greetpb.GreetAllRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mapa",
				LastName:  "Pratama",
			},
		},
		&greetpb.GreetAllRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Kupe",
				LastName:  "Pratama",
			},
		},
	}

	waitc := make(chan struct{})

	// send message to stream
	go func() {
		for _, req := range reqs {
			fmt.Printf("Sending request: %v\n", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("Error while sending stream: %v\n", err)
			}
			time.Sleep(700 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive messages from client
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while reading stream: %v\n", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

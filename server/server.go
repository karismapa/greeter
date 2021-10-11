package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/karismapa/greeter/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName + "!"
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function invoked with %v\n", req)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("Request has been canceled!")
			return nil, status.Error(codes.Canceled, "The client has canceled the request")
		}
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Request deadline exceeded!")
			return nil, status.Error(codes.DeadlineExceeded, "Request deadline has been exceeded")
		}

		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName + "!"
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)

		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		fmt.Printf("Send response: %v\n", res)
		if err := stream.Send(res); err != nil {
			fmt.Printf("Response %d got err %v\n", i, err)
		}
		time.Sleep(700 * time.Millisecond)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function invoked with a streaming request\n")

	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// finished reading the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
			return err
		}

		fmt.Printf("Receive request: %v\n", req)
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetAll(stream greetpb.GreetService_GreetAllServer) error {
	fmt.Printf("GreetAll function invoked with a streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client strea,: %v\n", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "
		sendErr := stream.Send(&greetpb.GreetAllResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

func main() {
	fmt.Println("Server is running...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../greetpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello from client!")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Client could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	//fmt.Printf("Created client: %f", c)
	// doServerStream(c)
	// doClientStream(c)
	doAllStream(c)
}

func doUnary(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Léo",
			LastName:  "Besançon",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Client could not greet: %v", err)
	}

	fmt.Printf("Response: %v\n", res.Result)

}

func doServerStream(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetStreamRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Léo",
			LastName:  "Besançon",
		},
	}

	resStream, err := c.GreetStream(context.Background(), req)
	if err != nil {
		log.Fatalf("Client could not greet: %v", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving stream: %v", err)
		}

		fmt.Printf("Response: %v\n", msg.GetResult())
	}
}

func doClientStream(c greetpb.GreetServiceClient) {

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Léo",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Léo2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Léo3",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Léo4",
			},
		},
	}

	resStream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while receiving stream: %v", err)
	}

	for _, req := range requests {
		resStream.Send(req)
	}

	res, err := resStream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving stream: %v", err)
	}

	fmt.Println(res)
}

func doAllStream(c greetpb.GreetServiceClient) {

	requests := []*greetpb.AllGreetRequest{
		&greetpb.AllGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Leo",
			},
		},
		&greetpb.AllGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Leo2",
			},
		},
		&greetpb.AllGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Leo3",
			},
		},
		&greetpb.AllGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Leo4",
			},
		},
	}

	resStream, err := c.AllGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while receiving stream: %v", err)
	}

	waitchannel := make(chan struct{})

	go func() {
		for _, req := range requests {

			err := resStream.Send(req)

			if err != nil {
				log.Fatalf("Error while sending stream: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			res, err := resStream.Recv()

			if err != nil {
				log.Fatalf("Error while receiving stream: %v", err)
			}

			fmt.Println(res)
		}
	}()

	<-waitchannel

}

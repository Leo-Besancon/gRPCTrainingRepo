package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../../calculator"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error in client: %v", err)
	}

	defer conn.Close()

	// c := calculator.NewAdditionServiceClient(conn)
	// c := calculator.NewDecompositionServiceClient(conn)
	// c := calculator.NewAverageServiceClient(conn)
	c := calculator.NewMaximumServiceClient(conn)

	doAllStream(c)
}

func doUnary(c calculator.AdditionServiceClient) {

	req := &calculator.AdditionRequest{
		Add: &calculator.Addition{
			FirstNum:  3,
			SecondNum: 10,
		},
	}

	ctx := context.Background()

	res, err := c.Add(ctx, req)

	if err != nil {
		log.Fatalf("Client could not greet: %v", err)
	}

	fmt.Printf("Added: %v", res.GetResult())
}

func doStreamDecomp(c calculator.DecompositionServiceClient) {

	req := &calculator.DecompositionRequest{
		Decomp: &calculator.Decomposition{
			Num: 45645123,
		},
	}

	resStream, err := c.Decomp(context.Background(), req)
	if err != nil {
		log.Fatalf("Client could not greet: %v", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			log.Fatalf("Error while receiving stream: %v", err)
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving stream: %v", err)
		}

		fmt.Printf("Response: %v\n", msg.GetResult())
	}
}

func doClientStream(c calculator.AverageServiceClient) {

	requests := []*calculator.AverageRequest{
		&calculator.AverageRequest{
			Avg: &calculator.Average{
				Num: 1,
			},
		},
		&calculator.AverageRequest{
			Avg: &calculator.Average{
				Num: 4,
			},
		},
		&calculator.AverageRequest{
			Avg: &calculator.Average{
				Num: 5,
			},
		},
		&calculator.AverageRequest{
			Avg: &calculator.Average{
				Num: 10,
			},
		},
		&calculator.AverageRequest{
			Avg: &calculator.Average{
				Num: 10,
			},
		},
	}

	resstream, err := c.Avg(context.Background())

	if err != nil {
		log.Fatalf("Error in client3: %v", err)
	}

	for _, req := range requests {

		resstream.Send(req)
	}

	res, err := resstream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error in client1: %v", err)
	}

	fmt.Printf("Response: %v", res.GetResult())
}

func doAllStream(c calculator.MaximumServiceClient) {
	requests := []*calculator.MaximumRequest{
		&calculator.MaximumRequest{
			Max: &calculator.Maximum{
				Num: 1,
			},
		},
		&calculator.MaximumRequest{
			Max: &calculator.Maximum{
				Num: 4,
			},
		},
		&calculator.MaximumRequest{
			Max: &calculator.Maximum{
				Num: 5,
			},
		},
		&calculator.MaximumRequest{
			Max: &calculator.Maximum{
				Num: 20,
			},
		},
		&calculator.MaximumRequest{
			Max: &calculator.Maximum{
				Num: 10,
			},
		},
	}

	waitc := make(chan struct{})

	resStream, err := c.Max(context.Background())

	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	go func() {

		for _, req := range requests {
			fmt.Printf("Sending number : %v \n", req.GetMax().GetNum())
			err := resStream.Send(req)

			if err != nil {
				log.Fatalf("Error : %v", err)
			}
			time.Sleep(time.Second)
		}

	}()

	go func() {
		for {
			res, err := resStream.Recv()

			if err != nil {
				log.Fatalf("Error : %v", err)
			}

			fmt.Printf("Current Max : %v \n", res.GetResult())
		}
	}()

	<-waitc

}

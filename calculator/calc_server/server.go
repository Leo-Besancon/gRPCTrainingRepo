package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"../../calculator"
	"google.golang.org/grpc"
)

type server struct {
}

func (*server) Add(ctx context.Context, req *calculator.AdditionRequest) (*calculator.AdditionResponse, error) {
	firstNum := req.GetAdd().GetFirstNum()
	secondNum := req.GetAdd().GetSecondNum()

	res := &calculator.AdditionResponse{Result: firstNum + secondNum}

	return res, nil
}

func (*server) Decomp(req *calculator.DecompositionRequest, stream calculator.DecompositionService_DecompServer) error {
	N := req.GetDecomp().GetNum()

	k := int64(2)

	for {
		if N < 1 {
			break
		}
		if N%k == 0 {
			res := &calculator.DecompositionResponse{
				Result: k,
			}
			stream.Send(res)
			N = N / k
		} else {
			k = k + 1
		}
	}

	return nil
}

func (*server) Avg(stream calculator.AverageService_AvgServer) error {

	curSum := int64(0)
	curReq := int64(0)

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			res := &calculator.AverageResponse{
				Result: float64(curSum) / float64(curReq),
			}

			fmt.Println(" result (server) ", float64(curSum)/float64(curReq))

			return stream.SendAndClose(res)
		}

		if err != nil {
			log.Fatalf("Failed to recieve stream! %v", err)
		}

		curSum += req.GetAvg().Num
		curReq++
	}
}

func (*server) Max(stream calculator.MaximumService_MaxServer) error {

	curMax := int64(0)

	for {
		req, err := stream.Recv()

		if err != nil {
			log.Fatalf("Failed to recieve stream! %v", err)
			return nil
		}

		if req.GetMax().GetNum() >= curMax {
			curMax = req.GetMax().GetNum()

			stream.Send(&calculator.MaximumResponse{
				Result: curMax,
			})
		}
	}
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("COULD NOT LISTEN! %v", err)
	}

	s := grpc.NewServer()
	// calculator.RegisterAdditionServiceServer(s, &server{})
	// calculator.RegisterDecompositionServiceServer(s, &server{})
	// calculator.RegisterAverageServiceServer(s, &server{})
	calculator.RegisterMaximumServiceServer(s, &server{})

	err2 := s.Serve(lis)

	if err2 != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

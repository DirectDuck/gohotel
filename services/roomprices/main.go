package main

import (
	"context"
	"log"
	"net"
	"os"

	"hotel/services/roomprices/rpc"

	"google.golang.org/grpc"
)

type RoomPriceServer struct {
	rpc.UnimplementedRoomPricesServiceServer
}

func (self *RoomPriceServer) GetRoomPrice(
	ctx context.Context, request *rpc.RoomPriceRequest,
) (*rpc.RoomPriceResponse, error) {
	roomType := request.GetType()
	return &rpc.RoomPriceResponse{
		Price: float64(roomType * 2),
	}, nil
}

func main() {
	listenUrl := os.Getenv("ROOMPRICES_LISTEN_URL")
	if len(listenUrl) == 0 {
		listenUrl = "0.0.0.0:8100"
	}

	listener, err := net.Listen("tcp", listenUrl)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	rpc.RegisterRoomPricesServiceServer(server, &RoomPriceServer{})

	log.Printf("gRPC server started on %s\n", listenUrl)
	err = server.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

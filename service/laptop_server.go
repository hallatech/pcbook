package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hallatech/pcbook/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LaptopServer is the server that provides laptop services
type LaptopServer struct {
	Store LaptopStore
}

// NewLaptopServer returns a new LaptopServer
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: store,
	}
}

// CreateLaptop is a unary RPC to create a new laptop
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("received a create-laptop request with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		// check if its a valid UUID
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	// some heavy processing for testing timeouts
	time.Sleep((6 * time.Second))

	//Handle the case where the client cancels the request
	if ctx.Err() == context.Canceled {
		log.Printf("request was canceled")
		return nil, status.Error(codes.Canceled, "request was canceled")
	}
	// If we don't want the server to continue processing and save the order we need to cater for it
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("deadline exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	}

	// save the laptop to DB/in mem storage
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	log.Printf("saved laptop with id: %s", laptop.Id)

	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}

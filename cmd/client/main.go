package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hallatech/pcbook/pb"
	"github.com/hallatech/pcbook/sample"
)

func main() {
	serverAddress := flag.String("address", "the server address", "0.0.0.0:<port>")
	flag.Parse()
	log.Printf("dial server: %s", serverAddress)

	// Use an insecure connection to avoid no transport security set error for this exercise
	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	laptop := sample.NewLaptop()
	// test with an empty string - should still create one with a new id
	// laptop.Id = ""
	// test with a copy of a previous Id - should return with laptop already exists message
	// laptop.Id = "e68f54c2-5e9a-44db-9dab-c3e8cc14439a"
	// test with an invalid ID - should return with an invalid UUID message
	// laptop.Id = "some-invalid-UUID"
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	// Basic context.Background
	// res, err := laptopClient.CreateLaptop(context.Background(), req)

	// Test with a set timeout - add a 6 sec wait in the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := laptopClient.CreateLaptop(ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			//not a big deal here
			log.Print("laptop already exists")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
		return
	}

	log.Printf("created laptop with id: %s", res.Id)
}

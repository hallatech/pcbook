package service_test

import (
	"context"
	"net"
	"testing"

	"github.com/hallatech/pcbook/pb"
	"github.com/hallatech/pcbook/sample"
	"github.com/hallatech/pcbook/serializer"
	"github.com/hallatech/pcbook/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, serverAddress := startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	// check that the laptop is saved to the store
	other, err := laptopServer.Store.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	// check that the saved laptop is the same as the one we send
	requireSameLaptop(t, laptop, other)

}

func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0") //random available port
	require.NoError(t, err)

	//grpcServer.Serve(listener) //blocking
	// so use go routine
	go grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()

}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	// use dial to make a connection
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure()) // test so use insecure connection
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func requireSameLaptop(t *testing.T, laptop1 *pb.Laptop, laptop2 *pb.Laptop) {
	// A simple require.Equal will fail because of the extra generated fields which will be different
	// require.Equal(t, laptop1, laptop2)

	// Another way is to serialize to JSON and compare
	json1, err := serializer.ProtobufToJSON(laptop1)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJSON(laptop2)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}

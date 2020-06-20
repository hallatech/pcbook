package serializer_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"

	"github.com/hallatech/pcbook/pb"
	"github.com/hallatech/pcbook/sample"
	"github.com/hallatech/pcbook/serializer"
)

func TestWriteProtobufToBinary(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := sample.NewLaptop()
	err := serializer.WriteProtobufToBinary(laptop1, binaryFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = serializer.ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop1, laptop2))

	err = serializer.WriteProtobufToJSON(laptop2, jsonFile)
	require.NoError(t, err)
}

package main

import (
	"testing"

	"github.com/place1/openapi-mock-server/generator"

	"github.com/stretchr/testify/require"
)

func TestStubResponse(t *testing.T) {
	require := require.New(t)

	stub, err := generator.NewStubGenerator("./petstore.yaml", generator.StubGeneratorOptions{})
	require.NoError(err)

	data, err := stub.StubResponse("/v1/pets", "GET")
	require.NoError(err)

	require.NotNil(data, "data should not be nil")
}

func TestStubResponseWithOverlay(t *testing.T) {
	require := require.New(t)

	stub, err := generator.NewStubGenerator("./petstore.yaml", generator.StubGeneratorOptions{
		Overlay:  "",
		BasePath: "/test-base-path",
	})
	require.NoError(err)

	_, err = stub.StubResponse("/test-base-path/pets", "GET")
	require.NoError(err)
}

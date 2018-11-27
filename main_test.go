package main

import (
	"openapimockserver/core"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStubResponse(t *testing.T) {
	require := require.New(t)

	stub, err := core.NewStubGenerator("./petstore.yaml", core.StubGeneratorOptions{})
	require.NoError(err)

	data, err := stub.StubResponse("/pets", "GET")
	require.NoError(err)

	require.NotNil(data, "data should not be nil")
}

func TestStubResponseWithOverlay(t *testing.T) {
	require := require.New(t)

	overlayPath := "./overlay.yaml"
	stub, err := core.NewStubGenerator("./petstore.yaml", core.StubGeneratorOptions{
		Overlay: &overlayPath,
	})
	require.NoError(err)

	data, err := stub.StubResponse("/pets", "GET")
	require.NoError(err)

	require.IsType([]interface{}{}, data)
	require.Len(data, 1)
	require.Contains(data.([]interface{})[0], "surprise")
}

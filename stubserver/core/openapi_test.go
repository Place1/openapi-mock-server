package core

import (
	"testing"

	"github.com/go-openapi/loads"

	"github.com/stretchr/testify/require"
)

func TestApplyResponseOverlay(t *testing.T) {
	require := require.New(t)

	data := map[string]interface{}{
		"example": "value",
	}

	overlay := Response{
		Content: `{"hello": "world"}`,
	}

	err := ApplyResponseOverlay(overlay, &data)
	require.NoError(err)

	require.Contains(data, "example")
	require.Contains(data, "hello")
}

func TestApplyResponseOverlayWithPrimatives(t *testing.T) {
	require := require.New(t)

	data := "my value"

	overlay := Response{
		Content: "my override",
	}

	err := ApplyResponseOverlay(overlay, &data)
	require.NoError(err)

	require.Equal("my override", data)
}

func TestExpandOperationIDs(t *testing.T) {
	require := require.New(t)

	document, err := loads.Spec("../../petstore.yaml")
	require.NoError(err)

	ExpandOperationIDs(document)

	require.Equal(document.Spec().Paths.Paths["/pets/{petId}"].Get.ID, "Get: /pets/{petId}")
}

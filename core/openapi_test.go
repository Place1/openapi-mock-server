package core

import (
	"testing"

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

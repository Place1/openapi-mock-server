package core

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-openapi/spec"
)

var ISO8601_DATE_STRING_FULL_RE = `^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d(\.\d+)?(([+-]\d\d:\d\d)|Z)?$`
var ISO8601_DATE_STRING_RE = `^\d{4}(-\d\d(-\d\d(T\d\d:\d\d(:\d\d)?(\.\d+)?(([+-]\d\d:\d\d)|Z)?)?)?)?$`

func TestStringStubDateTime(t *testing.T) {
	require := require.New(t)

	result := stringStub(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Format: "date-time",
		},
	})

	require.Regexp(
		regexp.MustCompile(ISO8601_DATE_STRING_FULL_RE),
		result,
		"date-time string should be a full ISO8601 date string",
	)
}

func TestStringStubDate(t *testing.T) {
	require := require.New(t)

	result := stringStub(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Format: "date",
		},
	})

	require.Regexp(
		regexp.MustCompile(ISO8601_DATE_STRING_RE),
		result,
		"date-time string should be a ISO8601 date string without time component",
	)
}

package main

import (
	"openapimockserver/core"
	"testing"

	"github.com/pkg/errors"
)

func TestStubResponse(t *testing.T) {
	stub, err := core.NewStubGenerator("./petstore.yaml")
	if err != nil {
		t.Fatal(errors.Wrap(err, "creating stub generator"))
	}

	data, err := stub.StubResponse("/pets", "GET")
	if err != nil {
		t.Fatal(errors.Wrap(err, "stubbing data"))
	}

	if data == nil {
		t.Error("data should not be nil")
	}
}

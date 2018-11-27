package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"openapimockserver/core"

	"github.com/pkg/errors"
)

type Options struct {
	Host string
	Port int
}

func OpenAPIStubServer(generator *core.StubGenerator, options *Options) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", options.Host, options.Port),
		Handler: &StubServerMux{generator: generator},
	}
	return server
}

type StubServerMux struct {
	generator *core.StubGenerator
}

func (mux *StubServerMux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	response, err := mux.generator.StubResponse(req.URL.Path, req.Method)
	if err != nil {
		log.Println(errors.Wrap(err, "unable to stub response"))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(response)
	if err != nil {
		log.Println(errors.Wrap(err, "unable to serialize generated response stub"))
		return
	}
}

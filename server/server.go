package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"openapimockserver/core"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
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
	mux.ValidationMiddleware(res, req)

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

func (mux *StubServerMux) ValidationMiddleware(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST", "PUT", "PATCH":
		operation, err := mux.generator.FindOperation(req.URL.Path, req.Method)
		if err != nil {
			log.Println(errors.Wrap(err, "finding operation schema for path and method"))
			break
		}

		bodyParam, err := core.FindBodyParam(operation)
		if err != nil {
			log.Println(errors.Wrap(err, "finding body parameter"))
			break
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(errors.Wrap(err, "reading response body"))
			break
		}

		jsonValue := map[string]interface{}{}
		err = json.Unmarshal(body, &jsonValue)
		if err != nil {
			log.Println(errors.Wrap(err, "decoding request body"))
			break
		}

		// run the validation
		err = validate.AgainstSchema(bodyParam.Schema, jsonValue, strfmt.Default)
		if err != nil {
			log.Println(errors.Wrapf(err, "%v: %v: %v", req.Method, req.URL.Path, string(body)))
			break
		}
	}
}

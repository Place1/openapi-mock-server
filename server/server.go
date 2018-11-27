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

// Options for the OpenAPIStubServer
type Options struct {
	Host string
	Port int
}

// OpenAPIStubServer returns an http.Server that pretends to be the API
// defined in the StubGenerator
func OpenAPIStubServer(generator *core.StubGenerator, options *Options) *http.Server {

	handler := createHandler(generator)
	handler = validationMiddleware(handler, generator)

	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", options.Host, options.Port),
		Handler: handler,
	}

	return server
}

func createHandler(generator *core.StubGenerator) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		response, err := generator.StubResponse(req.URL.Path, req.Method)
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
	})
}

func validationMiddleware(handler http.Handler, generator *core.StubGenerator) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "POST", "PUT", "PATCH":
			operation, err := generator.FindOperation(req.URL.Path, req.Method)
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
	})
}

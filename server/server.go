package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/place1/openapi-mock-server/generator"

	"github.com/pkg/errors"
)

// Options for the OpenAPIMockServer
type Options struct {
	Host string
	Port int
}

// OpenAPIMockServer returns an http.Server that pretends to be the API
// defined in the StubGenerator
func OpenAPIMockServer(generator *generator.StubGenerator, options *Options) *http.Server {

	handler := createHandler(generator)
	handler = requestLogger(handler)
	handler = validationMiddleware(handler, generator)

	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", options.Host, options.Port),
		Handler: handler,
	}

	return server
}

func createHandler(generator *generator.StubGenerator) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		response, err := generator.StubResponse(req.URL.Path, req.Method)
		if err != nil {
			log.Println(errors.Wrap(err, "unable to stub response"))
			http.Error(res, "stub server error - check the logs", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(res).Encode(response)
		if err != nil {
			log.Println(errors.Wrap(err, "unable to serialize generated response stub"))
			http.Error(res, "stub server error - check the logs", http.StatusInternalServerError)
			return
		}
	})
}

func validationMiddleware(handler http.Handler, generator *generator.StubGenerator) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "POST", "PUT", "PATCH":
			operation, err := generator.FindOperation(req.URL.Path, req.Method)
			if err != nil {
				log.Println(errors.Wrap(err, "finding operation schema for path and method"))
				break
			}

			err = ValidateConsumes(*operation, *req)
			if err != nil {
				log.Println(errors.Wrap(err, "validating content type header"))
			}

			err = ValidateParameters(*operation, *req)
			if err != nil {
				log.Println(errors.Wrap(err, "validating parameters"))
			}
		}
		handler.ServeHTTP(res, req)
	})
}

func requestLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%v %v", req.Method, req.URL.Path)
		handler.ServeHTTP(res, req)
	})
}

func errorLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(res, req)
	})
}

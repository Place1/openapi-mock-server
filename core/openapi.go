package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

// StubGenerator is the main type used to interact with this
// library's feature set
type StubGenerator struct {
	spec *spec.Swagger
}

// NewStubGenerator loads an OpenAPI spec from the given url/path
// and returns a StubGenerator
func NewStubGenerator(urlOrPath string) (*StubGenerator, error) {
	document, err := loads.Spec(urlOrPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load input file")
	}

	document, err = document.Expanded()
	if err != nil {
		return nil, errors.Wrap(err, "expanding spec refs")
	}

	stub := &StubGenerator{
		spec: document.Spec(),
	}

	return stub, nil
}

// StubResponse returns data that matches the schema for a given Operation
// in the OpenAPI spec. The Operation is determined by a path and method
func (stub *StubGenerator) StubResponse(path string, method string) (interface{}, error) {
	operation, err := FindOperation(stub.spec, path, method)
	if err != nil {
		return nil, errors.Wrap(err, "finding operation from path and method")
	}

	response, err := FindResponse(operation)
	if err != nil {
		return nil, errors.Wrap(err, "finding response for operation")
	}

	return StubSchema(response.Schema), nil
}

// FindOperation returns the best matching OpenAPI operation
// from the Spec given an HTTP Request
func FindOperation(openAPISpec *spec.Swagger, httpPath string, httpMethod string) (*spec.Operation, error) {
	// for every path, calculate a match score (most specific wins)
	scores := make(map[string]int)
	for path := range openAPISpec.Paths.Paths {
		matcher := pathToRegexp(httpPath)
		matches := matcher.FindAllString(path, -1)
		scores[path] = len(matches)
	}

	// pick the best matching path
	// a lower score means less matching segments and a more specific path.
	// we'll choose the most specific path
	var bestPath *string
	for path, score := range scores {
		if score != 0 && (bestPath == nil || score < scores[*bestPath]) {
			copy := string(path)
			bestPath = &copy
		}
	}

	if bestPath == nil {
		return nil, fmt.Errorf("unknown path %s", httpPath)
	}

	// find the operation from the pathItem using http method
	var operation *spec.Operation
	switch strings.ToUpper(httpMethod) {
	case "GET":
		operation = openAPISpec.Paths.Paths[*bestPath].Get
	case "POST":
		operation = openAPISpec.Paths.Paths[*bestPath].Post
	case "PUT":
		operation = openAPISpec.Paths.Paths[*bestPath].Put
	case "PATCH":
		operation = openAPISpec.Paths.Paths[*bestPath].Patch
	case "HEAD":
		operation = openAPISpec.Paths.Paths[*bestPath].Head
	case "OPTIONS":
		operation = openAPISpec.Paths.Paths[*bestPath].Options
	default:
		operation = nil
	}

	if operation == nil {
		return nil, fmt.Errorf("no operation for HTTP %s %s", httpMethod, httpPath)
	}

	return operation, nil
}

// pathToRegexp will convert an openapi path i.e. /api/{param}/thing/
// into a regexp like /api/(.*)/thing/
func pathToRegexp(path string) *regexp.Regexp {
	quotedPath := regexp.QuoteMeta(path)
	result := regexp.MustCompile(`(\\{\w+\\})`).ReplaceAllString(quotedPath, "(.*)")
	return regexp.MustCompile("^" + result + "$")
}

// FindResponse returns either the default response from an operation
// or the response with the lowest HTTP status code (i.e. success codes over error codes)
func FindResponse(operation *spec.Operation) (*spec.Response, error) {
	var response *spec.Response
	if operation.Responses.Default != nil {
		response = operation.Responses.Default
	} else {
		lowestCode := 999
		for code, res := range operation.Responses.StatusCodeResponses {
			if code < lowestCode {
				response = &res
			}
		}
	}

	if response == nil {
		return nil, fmt.Errorf("no response definition found for operation %s", operation.ID)
	}

	return response, nil
}

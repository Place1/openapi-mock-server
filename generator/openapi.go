package generator

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

// StubGeneratorOptions that can configure the stub generator
type StubGeneratorOptions struct {
	Overlay  string
	BasePath string
}

// StubGenerator is the main type used to interact with this
// library's feature set
type StubGenerator struct {
	spec    spec.Swagger
	overlay Overlay
}

// NewStubGenerator loads an OpenAPI spec from the given url/path
// and returns a StubGenerator
func NewStubGenerator(urlOrPath string, options StubGeneratorOptions) (*StubGenerator, error) {
	document, err := loads.Spec(urlOrPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load input file")
	}

	document, err = document.Expanded()
	if err != nil {
		return nil, errors.Wrap(err, "expanding spec refs")
	}

	// the openapi libraries suggest that the base path is
	// prefixed to paths: https://godoc.org/github.com/go-openapi/spec#Paths
	// but it doesn't seem to be happening in practice.
	// we'll expand them here
	ExpandPaths(document, options.BasePath)

	ExpandOperationIDs(document)

	var overlay *Overlay
	if options.Overlay != "" {
		overlay, err = LoadOverlayFile(options.Overlay)
		if err != nil {
			return nil, errors.Wrap(err, "loading overlay")
		}
	} else {
		tmp := EmptyOverlay()
		overlay = &tmp
	}

	stub := &StubGenerator{
		spec:    *document.Spec(),
		overlay: *overlay,
	}

	return stub, nil
}

// StubResponse returns data that matches the schema for a given Operation
// in the OpenAPI spec. The Operation is determined by a path and method
func (stub *StubGenerator) StubResponse(path string, method string) (interface{}, error) {
	operation, err := stub.FindOperation(path, method)
	if err != nil {
		return nil, errors.Wrap(err, "finding operation from path and method")
	}

	response, statusCode, err := stub.FindResponse(operation)
	if err != nil {
		return nil, errors.Wrap(err, "finding response for operation")
	}

	stubbedData := StubSchema(*response.Schema)

	if responseOverlay, err := stub.overlay.FindResponse(path, method, *statusCode); err == nil {
		ApplyResponseOverlay(*responseOverlay, &stubbedData)
	}

	return stubbedData, nil
}

// FindOperation returns the best matching OpenAPI operation
// from the Spec given an HTTP Request
func (stub *StubGenerator) FindOperation(httpPath string, httpMethod string) (*spec.Operation, error) {
	// for every path, calculate a match score
	// more matching path params means a higher score, 1 point per matching path param
	scores := make(map[string]int)
	for path := range stub.spec.Paths.Paths {
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
		operation = stub.spec.Paths.Paths[*bestPath].Get
	case "POST":
		operation = stub.spec.Paths.Paths[*bestPath].Post
	case "PUT":
		operation = stub.spec.Paths.Paths[*bestPath].Put
	case "PATCH":
		operation = stub.spec.Paths.Paths[*bestPath].Patch
	case "HEAD":
		operation = stub.spec.Paths.Paths[*bestPath].Head
	case "OPTIONS":
		operation = stub.spec.Paths.Paths[*bestPath].Options
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
func (stub *StubGenerator) FindResponse(operation *spec.Operation) (*spec.Response, *int, error) {
	var response *spec.Response

	lowestCode := 999
	for code, res := range operation.Responses.StatusCodeResponses {
		if code < lowestCode {
			tmp := res
			response = &tmp
			lowestCode = code
		}
	}

	if response == nil {
		return nil, nil, fmt.Errorf("no response definition found for operation %s", operation.ID)
	}

	return response, &lowestCode, nil
}

// ExpandPaths modifies all the paths in the openapi document
// by prefixing them with the basePath
func ExpandPaths(document *loads.Document, basePath string) {
	paths := map[string]spec.PathItem{}
	if basePath == "" {
		basePath = document.BasePath()
	}
	for apiPath, pathItem := range document.Spec().Paths.Paths {
		expandedPath := applyBasePath(basePath, apiPath)
		paths[expandedPath] = pathItem
	}
	document.Spec().Paths.Paths = paths
}

func applyBasePath(prefix string, suffix string) string {
	joint := path.Join(prefix, suffix)
	if strings.HasSuffix(suffix, "/") {
		// path.Join will clean a trailing slash.
		// some webservers actually care about the trailing slash
		// and so we want to preserve the "trailing slash-ness"
		// of the spec
		joint += "/"
	}
	return joint
}

// ExpandOperationIDs the operationId field can be omitted in
// the spec. Codegen tools will automatically generate a default
// value for this field but go-openapi does not.
// ExpandOperationIDs will expand empty operationId fields into
// a useful name for usage in error/logging.
func ExpandOperationIDs(document *loads.Document) {
	for path, pathItem := range document.Spec().Paths.Paths {
		if op := pathItem.Head; op != nil && op.ID == "" {
			op.ID = "Head: " + path
		}
		if op := pathItem.Options; op != nil && op.ID == "" {
			op.ID = "Options: " + path
		}
		if op := pathItem.Put; op != nil && op.ID == "" {
			op.ID = "Put: " + path
		}
		if op := pathItem.Get; op != nil && op.ID == "" {
			op.ID = "Get: " + path
		}
		if op := pathItem.Post; op != nil && op.ID == "" {
			op.ID = "Post: " + path
		}
		if op := pathItem.Patch; op != nil && op.ID == "" {
			op.ID = "Patch: " + path
		}
		if op := pathItem.Delete; op != nil && op.ID == "" {
			op.ID = "Delete: " + path
		}
	}
}

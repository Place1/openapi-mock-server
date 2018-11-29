package core

import (
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

type DocumentVisitor func(node interface{}, key string)

func Walk(node interface{}, key string, visitor DocumentVisitor) {
	visitor(node, key)

	switch v := node.(type) {
	case *loads.Document:
		Walk(v.Spec(), "", visitor)

	case *spec.Swagger:
		for apiPath, pathItem := range v.Paths.Paths {
			Walk(&pathItem, apiPath, visitor)
		}
		for definitionName, schema := range v.Definitions {
			Walk(&schema, definitionName, visitor)
		}
		for parameterName, parameter := range v.Parameters {
			Walk(&parameter, parameterName, visitor)
		}
		for responseName, response := range v.Responses {
			Walk(&response, responseName, visitor)
		}

	case *spec.PathItem:
		if v.Get != nil {
			Walk(v.Get, "Get", visitor)
		}
		if v.Post != nil {
			Walk(v.Post, "Post", visitor)
		}
		if v.Put != nil {
			Walk(v.Put, "Put", visitor)
		}
		if v.Patch != nil {
			Walk(v.Patch, "Patch", visitor)
		}
		if v.Delete != nil {
			Walk(v.Delete, "Delete", visitor)
		}
		if v.Options != nil {
			Walk(v.Options, "Options", visitor)
		}
		if v.Head != nil {
			Walk(v.Head, "Head", visitor)
		}

	case *spec.Operation:
		for _, parameter := range v.Parameters {
			Walk(&parameter, "", visitor)
		}
		for statusCode, response := range v.Responses.StatusCodeResponses {
			Walk(&response, string(statusCode), visitor)
		}

	case *spec.Response:
		if v.Schema != nil {
			Walk(v.Schema, "", visitor)
		}

	case *spec.Parameter:
		if v.Schema != nil {
			Walk(v.Schema, "", visitor)
		}
	}
}

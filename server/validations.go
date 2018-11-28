package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/pkg/errors"
)

func ValidateConsumes(operation spec.Operation, req http.Request) error {
	contentType := req.Header["Content-Type"]

	// if the spec says there's no content types consumed,
	// validate that no content type header exists
	if len(operation.Consumes) == 0 {
		if contentType != nil {
			return fmt.Errorf("operation %v expected no content type header but %v was found", operation.ID, contentType)
		}
		return nil
	}

	// if the spec says there a content type consumed
	// validate that the request content type is correct
	for _, consumes := range operation.Consumes {
		for _, item := range contentType {
			if item == consumes {
				return nil
			}
		}
	}

	return fmt.Errorf("operation %v expected to consume %v but found a content type of %v", operation.ID, operation.Consumes, contentType)
}

func ValidateParameters(operation spec.Operation, req http.Request) error {
	for _, parameter := range operation.Parameters {
		switch parameter.In {
		case "body":
			// TODO: confirm that a swagger spec can only have 1 body param
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return errors.Wrap(err, "reading response body")
			}

			if len(body) == 0 && parameter.Required {
				return fmt.Errorf("missing required body parameter")
			}

			// parse the request body
			jsonValue := map[string]interface{}{}
			err = json.Unmarshal(body, &jsonValue)
			if err != nil {
				return errors.Wrap(err, "decoding request body")
			}

			// run the validation
			err = validate.AgainstSchema(parameter.Schema, jsonValue, strfmt.Default)
			if err != nil {
				return errors.Wrap(err, "validating parameter")
			}

		case "query":
			paramName := parameter.Name
			value := req.URL.Query()[paramName]
			err := validate.AgainstSchema(parameter.Schema, value, strfmt.Default)
			if err != nil {
				return errors.Wrap(err, "validating parameter")
			}

		case "path":
			// TODO: implement path param validation
			break
		}
	}

	return nil
}

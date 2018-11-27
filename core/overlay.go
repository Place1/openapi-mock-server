package core

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/imdario/mergo"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Overlay struct {
	Paths map[string]PathItem `yaml:"paths"`
}

type PathItem struct {
	Get     *Operation `yaml:"get,omitempty"`
	Put     *Operation `yaml:"put,omitempty"`
	Post    *Operation `yaml:"post,omitempty"`
	Patch   *Operation `yaml:"patch,omitempty"`
	Options *Operation `yaml:"options,omitempty"`
	Head    *Operation `yaml:"head,omitempty"`
}

type Operation struct {
	Responses map[int]Response `yaml:"responses"`
}

type Response struct {
	Content string `yaml:"content"`
}

// LoadOverlayFile reads an overlay.yaml file into an Overlay struct
func LoadOverlayFile(path string) (*Overlay, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading overlay file")
	}

	overlay := &Overlay{}
	err = yaml.Unmarshal(content, overlay)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling overlay file")
	}

	return overlay, nil
}

// EmptyOverlay is used when the user doesn't provide
// and overlay file. it's just used inplace of a nil value.
func EmptyOverlay() Overlay {
	return Overlay{}
}

// ApplyResponseOverlay expects data to be passed by reference.
// The response overlay will be applied by merging/overriding data.
func ApplyResponseOverlay(response Response, data interface{}) error {
	switch reflect.Indirect(reflect.ValueOf(data)).Kind() {
	case reflect.String:
		*data.(*string) = string(response.Content)
		return nil

	case reflect.Map:
		var override map[string]interface{}
		err := json.Unmarshal([]byte(response.Content), &override)
		if err != nil {
			return errors.Wrap(err, "unmarshalling object response overlay")
		}
		err = mergo.Merge(data, override, mergo.WithOverride)
		if err != nil {
			return errors.Wrap(err, "merging response overlay with generated response stub")
		}
		return nil

	default:
		// I don't know why, but I can't match a reflect.Slice
		// so instead i'm handling slices in the default case
		// TODO: actually solve my problems...
		var override interface{}
		err := json.Unmarshal([]byte(response.Content), &override)
		if err != nil {
			return errors.Wrap(err, "unmarshalling array response overlay")
		}
		*data.(*interface{}) = override
		return nil
	}
}

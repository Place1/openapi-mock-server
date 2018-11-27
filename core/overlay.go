package core

import (
	"io/ioutil"

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
	Content     string  `yaml:"content"`
	ContentType *string `yaml:"contentType,omitempty"`
}

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

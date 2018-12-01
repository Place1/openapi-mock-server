package speclint

import (
	"openapimockserver/core"
	"openapimockserver/utils"
	"path"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

func NoEmptyDescriptions() core.DocumentVisitor {
	return func(node interface{}, data core.NodeData) {
		switch node := node.(type) {
		case *spec.Operation:
			if node.Description == "" {
				logrus.Warnf(`operation "%v" must have a description`, data.Ref)
			}
		}
	}
}

func NoEmptyOperationID() core.DocumentVisitor {
	return func(node interface{}, data core.NodeData) {
		switch node := node.(type) {
		case *spec.Operation:
			if node.ID == "" {
				logrus.Warnf(`operation "%v" must have an operationId`, data.Ref)
			}
		}
	}
}

func SlashTerminatedPaths() core.DocumentVisitor {
	return func(node interface{}, data core.NodeData) {
		switch node.(type) {
		case *spec.PathItem:
			if !strings.HasSuffix(data.Key, "/") {
				logrus.Warnf(`path "%v" must end with a slash`, data.Ref)
			}
		}
	}
}

type NamingConvention = string

const (
	PascalCase NamingConvention = "PascalCase"
	SnakeCase  NamingConvention = "SnakeCase"
	CamelCase  NamingConvention = "CamelCase"
	KebabCase  NamingConvention = "KebabCase"
)

func GetNamingChecker(convention NamingConvention) func(string) bool {
	var checker func(string) bool
	switch convention {
	case PascalCase:
		checker = utils.IsPascalCase
	case CamelCase:
		checker = utils.IsCamelCase
	case SnakeCase:
		checker = utils.IsSnakeCase
	case KebabCase:
		checker = utils.IsKebabCase
	}
	return checker
}

func PathNamingConvention(convention NamingConvention) core.DocumentVisitor {
	checker := GetNamingChecker(convention)
	return func(node interface{}, data core.NodeData) {
		switch node.(type) {
		case *spec.PathItem:
			for _, segment := range strings.Split(data.Key, "/") {
				if segment != "" && !checker(segment) {
					logrus.Warnf(`path "%v" must follow the %v naming convention`, data.Ref, convention)
				}
			}
		}
	}
}

func DefinitionNamingConvention(convention NamingConvention) core.DocumentVisitor {
	checker := GetNamingChecker(convention)
	return func(node interface{}, data core.NodeData) {
		switch node := node.(type) {
		case *spec.Definitions:
			for name := range *node {
				if !checker(name) {
					logrus.Warnf(`definition "%v" must follow the %v naming convention`, path.Join(data.Ref, name), convention)
				}
			}
		}
	}
}

func PropertyNamingConvention(convention NamingConvention) core.DocumentVisitor {
	checker := GetNamingChecker(convention)
	return func(node interface{}, data core.NodeData) {
		switch node := node.(type) {
		case *spec.Definitions:
			for name, schema := range *node {
				core.Walk(&schema, core.NodeData{Ref: path.Join(data.Ref, name)}, func(node interface{}, data core.NodeData) {
					switch node.(type) {
					case *spec.Schema:
						if !checker(data.Key) {
							logrus.Warnf(`property "%v" must follow the %v naming convention`, data.Ref, convention)
						}
					}
				})
			}
		}
	}
}

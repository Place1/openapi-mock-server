package speclint

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Spec string
}

func RunSpecLint(options Options) {
	document, err := loads.Spec(options.Spec)
	if err != nil {
		log.Fatalln(err)
	}

	validator := validate.NewSpecValidator(document.Schema(), strfmt.Default)
	validator.SetContinueOnErrors(true)
	result, _ := validator.Validate(document)
	for _, specError := range result.Errors {
		logrus.Error(specError.Error())
	}
	for _, warning := range result.Warnings {
		logrus.Warn(warning.Error())
	}
}

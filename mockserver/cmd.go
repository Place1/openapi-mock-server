package mockserver

import (
	"log"

	"github.com/place1/openapi-mock-server/generator"
	"github.com/place1/openapi-mock-server/server"
)

type Options struct {
	Spec     string
	Overlay  string
	BasePath string
	Host     string
	Port     int
}

func Runmockserver(options Options) {
	stub, err := generator.NewStubGenerator(options.Spec, generator.StubGeneratorOptions{
		Overlay:  options.Overlay,
		BasePath: options.BasePath,
	})
	if err != nil {
		log.Fatalln(err)
	}

	server := server.OpenAPIMockServer(stub, &server.Options{
		Host: options.Host,
		Port: options.Port,
	})

	log.Printf("listening on %v:%v\n", options.Host, options.Port)
	log.Fatalln(server.ListenAndServe())
}

package main

import (
	"log"

	"github.com/place1/openapi-mock-server/server"

	"github.com/place1/openapi-mock-server/generator"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	serveSpec     = kingpin.Arg("openapi-spec", "the path to an openapi spec yaml file").Required().String()
	serveHost     = kingpin.Flag("host", "the host or ip address that the server should listen on.").Default("127.0.0.1").String()
	servePort     = kingpin.Flag("port", "the port that the server should listen on.").Default("8000").Int()
	serveOverlay  = kingpin.Flag("overlay", "path to an overlay.yaml file.").Default("").String()
	serveBasePath = kingpin.Flag("base-path", "override the basePath defined in the spec. defaults to the value defined in the spec.").Default("").String()
)

func main() {
	kingpin.Parse()
	Runmockserver(Options{
		Spec:     *serveSpec,
		Host:     *serveHost,
		Port:     *servePort,
		Overlay:  *serveOverlay,
		BasePath: *serveBasePath,
	})
}

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

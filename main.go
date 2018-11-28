package main

import (
	"log"
	"openapimockserver/stubserver/core"
	"openapimockserver/stubserver/server"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	spec := kingpin.Arg("openapi-spec", "the path to an openapi spec yaml file").Required().String()
	host := kingpin.Flag("host", "the host or ip address that the server should listen on.").Default("127.0.0.1").String()
	port := kingpin.Flag("port", "the port that the server should listen on.").Default("8000").Int()
	overlay := kingpin.Flag("overlay", "path to an overlay.yaml file.").Default("").String()
	basePath := kingpin.Flag("base-path", "override the basePath defined in the spec. defaults to the value defined in the spec.").Default("").String()
	kingpin.Parse()

	if *spec == "" {
		log.Fatalln("missing positional argument <openapi-spec.yaml>")
	}

	stub, err := core.NewStubGenerator(*spec, core.StubGeneratorOptions{
		Overlay:  *overlay,
		BasePath: *basePath,
	})
	if err != nil {
		log.Fatalln(err)
	}

	server := server.OpenAPIStubServer(stub, &server.Options{
		Host: *host,
		Port: *port,
	})

	log.Printf("listening on %v:%v\n", *host, *port)
	log.Fatalln(server.ListenAndServe())
}

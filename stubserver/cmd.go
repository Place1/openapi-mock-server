package stubserver

import (
	"log"
	"openapimockserver/stubserver/core"
	"openapimockserver/stubserver/server"
)

type Options struct {
	Spec     string
	Overlay  string
	BasePath string
	Host     string
	Port     int
}

func RunStubServer(options Options) {
	stub, err := core.NewStubGenerator(options.Spec, core.StubGeneratorOptions{
		Overlay:  options.Overlay,
		BasePath: options.BasePath,
	})
	if err != nil {
		log.Fatalln(err)
	}

	server := server.OpenAPIStubServer(stub, &server.Options{
		Host: options.Host,
		Port: options.Port,
	})

	log.Printf("listening on %v:%v\n", options.Host, options.Port)
	log.Fatalln(server.ListenAndServe())
}

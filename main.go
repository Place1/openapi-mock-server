package main

import (
	"openapimockserver/stubserver"

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
	stubserver.RunStubServer(stubserver.Options{
		Spec:     *serveSpec,
		Host:     *serveHost,
		Port:     *servePort,
		Overlay:  *serveOverlay,
		BasePath: *serveBasePath,
	})
}

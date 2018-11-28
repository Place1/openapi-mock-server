package main

import (
	"openapimockserver/speclint"
	"openapimockserver/stubserver"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// serve command
	serve         = kingpin.Command("serve", "serve a API stub server")
	serveSpec     = serve.Arg("openapi-spec", "the path to an openapi spec yaml file").Required().String()
	serveHost     = serve.Flag("host", "the host or ip address that the server should listen on.").Default("127.0.0.1").String()
	servePort     = serve.Flag("port", "the port that the server should listen on.").Default("8000").Int()
	serveOverlay  = serve.Flag("overlay", "path to an overlay.yaml file.").Default("").String()
	serveBasePath = serve.Flag("base-path", "override the basePath defined in the spec. defaults to the value defined in the spec.").Default("").String()

	// lint command
	lint     = kingpin.Command("lint", "lint an openapi spec")
	lintSpec = lint.Arg("openapi-spec", "the path to an openapi spec yaml file").Required().String()
)

func main() {
	switch kingpin.Parse() {
	case "serve":
		stubserver.RunStubServer(stubserver.Options{
			Spec:     *serveSpec,
			Host:     *serveHost,
			Port:     *servePort,
			Overlay:  *serveOverlay,
			BasePath: *serveBasePath,
		})

	case "lint":
		speclint.RunSpecLint(speclint.Options{
			Spec: *lintSpec,
		})
	}
}

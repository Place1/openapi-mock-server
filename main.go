package main

import (
	"log"
	"openapimockserver/core"
	"openapimockserver/server"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Description = "An automatic server stub powered by OpenAPI and Go"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
		cli.StringFlag{
			Name:  "host",
			Value: "127.0.0.1",
			Usage: "the host or ip address that the server should listen on",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 8000,
			Usage: "the port that the server should listen on",
		},
		cli.StringFlag{
			Name:  "overlay",
			Value: "",
			Usage: "path to an overlay.yaml file",
		},
	}

	app.Action = func(context *cli.Context) error {
		// if c.NArg() > 0 {
		// 	name = c.Args().Get(0)
		// }
		// if c.String("lang") == "spanish" {
		// 	fmt.Println("Hola", name)
		// } else {
		// 	fmt.Println("Hello", name)
		// }
		// return nil

		openAPISpec := context.Args().First()
		host := context.String("host")
		port := context.Int("port")
		overlay := context.String("overlay")

		stub, err := core.NewStubGenerator(openAPISpec, core.StubGeneratorOptions{
			Overlay: &overlay,
		})
		if err != nil {
			log.Fatalln(err)
		}

		server := server.OpenAPIStubServer(stub, &server.Options{
			Host: host,
			Port: port,
		})

		log.Printf("listening on %v:%v\n", host, port)
		return server.ListenAndServe()
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

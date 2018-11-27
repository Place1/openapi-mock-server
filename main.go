package main

import (
	"log"
	"openapimockserver/core"
	"openapimockserver/server"

	"github.com/spf13/cobra"
)

var overlay string

var rootCmd = &cobra.Command{
	Use:   "openapi-server-stub",
	Short: "An automatic server stub powered by OpenAPI and Go",
	Run: func(cmd *cobra.Command, args []string) {
		stub, err := core.NewStubGenerator("./petstore.yaml", core.StubGeneratorOptions{
			Overlay: &overlay,
		})
		if err != nil {
			log.Fatalln(err)
		}

		server := server.OpenAPIStubServer(stub, &server.Options{
			Host: "127.0.0.1",
			Port: 8000,
		})

		log.Println("listening on :8000")
		log.Fatalln(server.ListenAndServe())
	},
}

func init() {
	rootCmd.Flags().StringVarP(&overlay, "overlay", "", "", "path to an overlay.yaml file")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

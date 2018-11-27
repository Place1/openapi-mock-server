package main

import (
	"log"
	"openapimockserver/core"
	"openapimockserver/server"
)

func main() {
	stub, err := core.NewStubGenerator("./petstore.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	server := server.OpenAPIStubServer(stub, &server.Options{
		Host: "127.0.0.1",
		Port: 8000,
	})

	log.Println("listening on :8000")
	log.Fatalln(server.ListenAndServe())
}

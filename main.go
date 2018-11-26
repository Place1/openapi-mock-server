package main

import (
	"fmt"
	"log"
	"openapimockserver/runtime"
)

func main() {
	stub, err := runtime.NewStubGenerator("./petstore.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	data, err := stub.StubResponse("/pets", "GET")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(data)
}

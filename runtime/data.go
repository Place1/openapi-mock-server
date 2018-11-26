package runtime

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/go-openapi/spec"
)

// StubSchema returns a struct that matches the openapi schema
// with values filled with randomly generated data
func StubSchema(schema *spec.Schema) interface{} {
	if schema.Type.Contains("object") {
		return objectStub(schema)

	} else if schema.Type.Contains("array") {
		return arrayStub(schema)

	} else if schema.Type.Contains("string") {
		return stringStub()

	} else if schema.Type.Contains("number") {
		return 10

	} else if schema.Type.Contains("integer") {
		return 10

	} else if schema.Type.Contains("boolean") {
		return true
	} else if len(schema.Properties) != 0 {
		// there was no `type` field on the schema
		// but there were `properties`.
		// we'll log a warning and then hope it's actually an object
		log.Printf("unknown schema type %s. assuming type object", schema.Type)
		return objectStub(schema)
	}

	panic(fmt.Sprintf("unknown schema type %s for schema %s", schema.Type, schema.ID))
}

func objectStub(schema *spec.Schema) interface{} {
	obj := map[string]interface{}{}
	for property, propSchema := range schema.Properties {
		obj[property] = StubSchema(&propSchema)
	}
	return obj
}

func arrayStub(schema *spec.Schema) []interface{} {
	items := make([]interface{}, 1)
	items[0] = StubSchema(schema.Items.Schema)
	return items
}

func stringStub() string {
	values := []string{
		"lorem",
		"ipsum",
		"hello world",
	}
	return values[randInt(0, len(values)-1)]
}

func integerStub() int {
	return randInt(-1000, 1000)
}

func numberStub() float32 {
	return rand.Float32() * 1000
}

func booleanStub() bool {
	// Intn is [0, n) non inclusive, so we need n=2 to get [0,1]
	return rand.Intn(2) == 0
}

func randInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

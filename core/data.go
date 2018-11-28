package core

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"time"

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
		return stringStub(schema)

	} else if schema.Type.Contains("number") {
		return numberStub()

	} else if schema.Type.Contains("integer") {
		return integerStub()

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

func stringStub(schema *spec.Schema) interface{} {
	if len(schema.Enum) != 0 {
		// if the schema defines an enum we will choose
		// one of the values at random
		return schema.Enum[randInt(0, len(schema.Enum)-1)]
	}

	switch schema.Format {
	case "date":
		return time.Now().Format("2006-01-02")
	case "date-time":
		return time.Now().Format(time.RFC3339)
	case "byte":
		return base64.StdEncoding.EncodeToString([]byte(FakeString()))
	case "binary":
		return []byte(FakeString())
	default:
		return FakeString()
	}
}

func integerStub() int {
	return randInt(0, 100)
}

func numberStub() float32 {
	return rand.Float32() * 100
}

func booleanStub() bool {
	return randInt(0, 1) == 0
}

// randInt returns a random number within the given
// bounds [min, max] inclusive.
func randInt(min int, max int) int {
	// rand.Intn is non-inclusive of the upper bound
	// so we +1 to get an inclusive upper bound
	return rand.Intn(max-min+1) + min
}

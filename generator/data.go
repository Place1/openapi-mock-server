package generator

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
func StubSchema(schema spec.Schema) interface{} {
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
		log.Printf("unknown schema type \"%v\". assuming type object", schema.Type)
		return objectStub(schema)
	}

	panic(fmt.Sprintf("unknown schema type \"%v\" for schema \"%v\"", schema.Type, schema.ID))
}

func objectStub(schema spec.Schema) interface{} {
	obj := map[string]interface{}{}
	for property, propSchema := range schema.Properties {
		obj[property] = StubSchema(propSchema)
	}
	return obj
}

func arrayStub(schema spec.Schema) []interface{} {
	if schema.Items.Schema == nil {
		log.Printf("schema \"%v\" of type array missing items schema - stub will be an empty array", schema.ID)
		return []interface{}{}
	}

	size := randInt(0, 10)
	items := make([]interface{}, size)
	for i := 0; i < size; i++ {
		items[i] = StubSchema(*schema.Items.Schema)
	}

	return items
}

func stringStub(schema spec.Schema) interface{} {
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
		return base64.StdEncoding.EncodeToString([]byte(generateString()))
	case "binary":
		return []byte(generateString())
	default:
		return generateString()
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

func generateString() string {
	values := []string{
		"lorem ipsum",
		"neque porro",
		"dolorem ipsum",
	}
	return values[randInt(0, len(values)-1)]
}

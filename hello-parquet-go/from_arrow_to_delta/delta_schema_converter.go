package main

import (
	"encoding/json"
	"fmt"
	delta "parquet-project/delta"
	"strings"

	"github.com/apache/arrow/go/v13/arrow"
)

func main() {
	// Your Arrow schema
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "archer", Type: arrow.BinaryTypes.String},
			{Name: "location", Type: arrow.BinaryTypes.String},
			{Name: "year", Type: arrow.PrimitiveTypes.Int16},
		},
		nil,
	)

	// Convert to Delta Lake schema string
	deltaSchemaString, err := delta.ArrowSchemaToDeltaSchemaString(schema)
	if err != nil {
		panic(fmt.Sprintf("failed to convert schema: %v", err))
	}

	fmt.Println("Arrow Schema converted to Delta Lake Schema String:")
	fmt.Println(deltaSchemaString)

	// Pretty print for readability
	var prettyJSON map[string]interface{}
	json.Unmarshal([]byte(deltaSchemaString), &prettyJSON)
	prettyBytes, _ := json.MarshalIndent(prettyJSON, "", "  ")
	fmt.Println("\nPretty formatted:")
	fmt.Println(string(prettyBytes))

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("COMPREHENSIVE TESTS FOR ALL DELTA LAKE DATA TYPES")
	fmt.Println(strings.Repeat("=", 60))

	// Test all scalar types
	fmt.Println("\n1. Testing Scalar Types:")
	scalarSchema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "null_field", Type: arrow.Null, Nullable: true},
			{Name: "string_field", Type: arrow.BinaryTypes.String, Nullable: true},
			{Name: "binary_field", Type: arrow.BinaryTypes.Binary, Nullable: false},
			{Name: "byte_field", Type: arrow.PrimitiveTypes.Int8, Nullable: true},
			{Name: "short_field", Type: arrow.PrimitiveTypes.Int16, Nullable: true},
			{Name: "integer_field", Type: arrow.PrimitiveTypes.Int32, Nullable: true},
			{Name: "long_field", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
			{Name: "float_field", Type: arrow.PrimitiveTypes.Float32, Nullable: true},
			{Name: "double_field", Type: arrow.PrimitiveTypes.Float64, Nullable: true},
			{Name: "boolean_field", Type: arrow.FixedWidthTypes.Boolean, Nullable: true},
			{Name: "date32_field", Type: arrow.FixedWidthTypes.Date32, Nullable: true},
			{Name: "date64_field", Type: arrow.FixedWidthTypes.Date64, Nullable: true},
			{Name: "timestamp_field", Type: arrow.FixedWidthTypes.Timestamp_ms, Nullable: true},
			{Name: "decimal128_field", Type: &arrow.Decimal128Type{Precision: 10, Scale: 2}, Nullable: true},
		},
		nil,
	)

	scalarDeltaSchemaString, err := delta.ArrowSchemaToDeltaSchemaString(scalarSchema)
	if err != nil {
		panic(fmt.Sprintf("failed to convert scalar schema: %v", err))
	}

	var scalarPrettyJSON map[string]interface{}
	json.Unmarshal([]byte(scalarDeltaSchemaString), &scalarPrettyJSON)
	scalarPrettyBytes, _ := json.MarshalIndent(scalarPrettyJSON, "", "  ")
	fmt.Println(string(scalarPrettyBytes))

	// Test complex types
	fmt.Println("\n2. Testing Complex Types:")

	// Create a struct type
	structType := arrow.StructOf(
		arrow.Field{Name: "nested_string", Type: arrow.BinaryTypes.String, Nullable: true},
		arrow.Field{Name: "nested_int", Type: arrow.PrimitiveTypes.Int32, Nullable: true},
	)

	// Create a list type
	listType := arrow.ListOf(arrow.BinaryTypes.String)

	// Create a map type
	mapType := arrow.MapOf(arrow.BinaryTypes.String, arrow.PrimitiveTypes.Int32)

	complexSchema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "struct_field", Type: structType, Nullable: true},
			{Name: "array_field", Type: listType, Nullable: true},
			{Name: "map_field", Type: mapType, Nullable: true},
		},
		nil,
	)

	complexDeltaSchemaString, err := delta.ArrowSchemaToDeltaSchemaString(complexSchema)
	if err != nil {
		panic(fmt.Sprintf("failed to convert complex schema: %v", err))
	}

	var complexPrettyJSON map[string]interface{}
	json.Unmarshal([]byte(complexDeltaSchemaString), &complexPrettyJSON)
	complexPrettyBytes, _ := json.MarshalIndent(complexPrettyJSON, "", "  ")
	fmt.Println(string(complexPrettyBytes))

	// Test edge cases and nullable variations
	fmt.Println("\n3. Testing Nullable Variations:")
	nullableSchema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "nullable_true", Type: arrow.BinaryTypes.String, Nullable: true},
			{Name: "nullable_false", Type: arrow.BinaryTypes.String, Nullable: false},
			{Name: "nullable_default", Type: arrow.BinaryTypes.String}, // Default is false
		},
		nil,
	)

	nullableDeltaSchemaString, err := delta.ArrowSchemaToDeltaSchemaString(nullableSchema)
	if err != nil {
		panic(fmt.Sprintf("failed to convert nullable schema: %v", err))
	}

	var nullablePrettyJSON map[string]interface{}
	json.Unmarshal([]byte(nullableDeltaSchemaString), &nullablePrettyJSON)
	nullablePrettyBytes, _ := json.MarshalIndent(nullablePrettyJSON, "", "  ")
	fmt.Println(string(nullablePrettyBytes))

	// Summary
	fmt.Println("\n4. Type Mapping Summary:")
	fmt.Println("\nArrow Type        → Delta Lake Type")
	fmt.Println(strings.Repeat("=", 35))
	typeTests := map[arrow.Type]string{
		arrow.NULL:       "null",
		arrow.STRING:     "string",
		arrow.BINARY:     "binary",
		arrow.INT8:       "byte",
		arrow.INT16:      "short",
		arrow.INT32:      "integer",
		arrow.INT64:      "long",
		arrow.FLOAT32:    "float",
		arrow.FLOAT64:    "double",
		arrow.BOOL:       "boolean",
		arrow.DATE32:     "date",
		arrow.DATE64:     "date",
		arrow.TIMESTAMP:  "timestamp",
		arrow.DECIMAL128: "decimal",
		arrow.DECIMAL256: "decimal",
		arrow.STRUCT:     "struct",
		arrow.LIST:       "array",
		arrow.MAP:        "map",
	}

	for arrowType, deltaType := range typeTests {
		fmt.Printf("%-17s → %s\n", arrowType.String(), deltaType)
	}
}

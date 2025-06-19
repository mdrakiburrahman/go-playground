package parquet_project

import (
	"encoding/json"
	"fmt"

	"github.com/apache/arrow/go/v13/arrow"
)

// DeltaField represents a field in Delta Lake schema format
type DeltaField struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Nullable bool        `json:"nullable"`
	Metadata interface{} `json:"metadata"`
}

// DeltaSchema represents the root schema structure in Delta Lake format
type DeltaSchema struct {
	Type   string       `json:"type"`
	Fields []DeltaField `json:"fields"`
}

// arrowTypeToDeltaType converts Arrow data types to Delta Lake type strings
func arrowTypeToDeltaType(arrowType arrow.DataType) string {
	switch arrowType.ID() {
	case arrow.NULL:
		return "null"
	case arrow.STRING:
		return "string"
	case arrow.BINARY:
		return "binary"
	case arrow.INT8:
		return "byte"
	case arrow.INT16:
		return "short"
	case arrow.INT32:
		return "integer"
	case arrow.INT64:
		return "long"
	case arrow.FLOAT32:
		return "float"
	case arrow.FLOAT64:
		return "double"
	case arrow.BOOL:
		return "boolean"
	case arrow.DATE32, arrow.DATE64:
		return "date"
	case arrow.TIMESTAMP:
		return "timestamp"
	case arrow.DECIMAL128, arrow.DECIMAL256:
		return "decimal"
	case arrow.STRUCT:
		return "struct"
	case arrow.LIST:
		return "array"
	case arrow.MAP:
		return "map"
	default:
		// For unsupported types, return string as fallback
		return "string"
	}
}

// ArrowSchemaToDeltaSchemaString converts an Arrow schema to Delta Lake schema JSON string
func ArrowSchemaToDeltaSchemaString(schema *arrow.Schema) (string, error) {
	fields := make([]DeltaField, len(schema.Fields()))

	for i, field := range schema.Fields() {
		fields[i] = DeltaField{
			Name:     field.Name,
			Type:     arrowTypeToDeltaType(field.Type),
			Nullable: field.Nullable,
			Metadata: map[string]interface{}{}, // Empty metadata as requested
		}
	}

	deltaSchema := DeltaSchema{
		Type:   "struct",
		Fields: fields,
	}

	jsonBytes, err := json.Marshal(deltaSchema)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

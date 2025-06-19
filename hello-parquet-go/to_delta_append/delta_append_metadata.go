package main

import (
	"fmt"
	"parquet-project/delta"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/google/uuid"
)

func main() {
	fmt.Println("Delta Lake Append Metadata Generator Demo")
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")

	// Create an Arrow schema similar to your original example
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "Body", Type: arrow.BinaryTypes.Binary, Nullable: true},
			{Name: "BatchEnqueuedUnixTimeMs", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
			{Name: "BatchIngestionUnixTimeMs", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
			{Name: "BatchOffset", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
			{Name: "BatchSequenceNumber", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
			{Name: "BatchPartitionId", Type: arrow.BinaryTypes.String, Nullable: true},
			{Name: "YearMonthDate", Type: arrow.BinaryTypes.String, Nullable: true},
		},
		nil,
	)

	// Example 1: Using the convenience function with Arrow schema
	fmt.Println("\n1. Using Arrow Schema Convenience Function:")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	deltaMetadata1, err := delta.GenerateDeltaAppendMetadataFromArrowSchema(
		schema,
		uuid.New().String(),
		[]string{"YearMonthDate"},
		time.Now().UnixMilli(),
		1, // minReaderVersion
		2, // minWriterVersion
		"YearMonthDate=20250612/part-00000-f3ffe72a-98fd-4dfe-9e12-96030934ce3f.c000.snappy.parquet",
		map[string]string{"YearMonthDate": "20250612"},
		448876,
		time.Now().UnixMilli(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to generate metadata from Arrow schema: %v", err))
	}

	fmt.Println("Generated Delta Lake Transaction Log:")
	fmt.Println(deltaMetadata1)

	// Example 2: Using manual configuration
	fmt.Println("\n\n2. Using Manual Configuration:")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	// First convert the schema to string
	schemaString, err := delta.ArrowSchemaToDeltaSchemaString(schema)
	if err != nil {
		panic(fmt.Sprintf("failed to convert schema: %v", err))
	}

	config := delta.AppendFileGeneratorConfig{
		ID:               "bbf06f99-2376-4d41-b4fb-367dd31df3de",
		SchemaString:     schemaString,
		PartitionColumns: []string{"YearMonthDate"},
		CreatedTime:      1749738853893,
		MinReaderVersion: 1,
		MinWriterVersion: 2,
		Path:             "YearMonthDate=20250612/part-00000-f3ffe72a-98fd-4dfe-9e12-96030934ce3f.c000.snappy.parquet",
		PartitionValues:  map[string]string{"YearMonthDate": "20250612"},
		Size:             448876,
		ModificationTime: 1749738853808,
	}

	deltaMetadata2, err := delta.GenerateDeltaAppendMetadata(config)
	if err != nil {
		panic(fmt.Sprintf("failed to generate metadata: %v", err))
	}

	fmt.Println("Generated Delta Lake Transaction Log (Manual Config):")
	fmt.Println(deltaMetadata2)

	// Example 3: Simple schema example
	fmt.Println("\n\n3. Simple Schema Example:")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	simpleSchema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "archer", Type: arrow.BinaryTypes.String, Nullable: true},
			{Name: "location", Type: arrow.BinaryTypes.String, Nullable: true},
			{Name: "year", Type: arrow.PrimitiveTypes.Int16, Nullable: true},
		},
		nil,
	)

	deltaMetadata3, err := delta.GenerateDeltaAppendMetadataFromArrowSchema(
		simpleSchema,
		uuid.New().String(),
		[]string{}, // No partition columns
		time.Now().UnixMilli(),
		1, // minReaderVersion
		2, // minWriterVersion
		"part-00000-simple-data.snappy.parquet",
		map[string]string{}, // No partition values
		12345,
		time.Now().UnixMilli(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to generate simple metadata: %v", err))
	}

	fmt.Println("Generated Simple Delta Lake Transaction Log:")
	fmt.Println(deltaMetadata3)

	// Example 4: Multiple partition columns
	fmt.Println("\n\n4. Multiple Partition Columns Example:")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	multiPartitionSchema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
			{Name: "name", Type: arrow.BinaryTypes.String, Nullable: true},
			{Name: "amount", Type: arrow.PrimitiveTypes.Float64, Nullable: true},
			{Name: "year", Type: arrow.PrimitiveTypes.Int32, Nullable: true},
			{Name: "month", Type: arrow.PrimitiveTypes.Int32, Nullable: true},
			{Name: "day", Type: arrow.PrimitiveTypes.Int32, Nullable: true},
		},
		nil,
	)

	deltaMetadata4, err := delta.GenerateDeltaAppendMetadataFromArrowSchema(
		multiPartitionSchema,
		uuid.New().String(),
		[]string{"year", "month", "day"},
		time.Now().UnixMilli(),
		1, // minReaderVersion
		2, // minWriterVersion
		"year=2025/month=06/day=19/part-00000-multi-partition.snappy.parquet",
		map[string]string{
			"year":  "2025",
			"month": "06",
			"day":   "19",
		},
		987654,
		time.Now().UnixMilli(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to generate multi-partition metadata: %v", err))
	}

	fmt.Println("Generated Multi-Partition Delta Lake Transaction Log:")
	fmt.Println(deltaMetadata4)

	fmt.Println("\n\nDemo completed successfully!")
}

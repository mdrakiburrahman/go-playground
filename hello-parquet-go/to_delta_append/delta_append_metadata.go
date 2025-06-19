package main

import (
	"encoding/json"
	"fmt"
	"parquet-project/delta"
	"strings"

	"github.com/apache/arrow/go/v13/arrow"
)

func main() {
	fmt.Println("Delta Lake Append Transaction Notification Demo")
	fmt.Println(strings.Repeat("=", 60))

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

	schemaString, err := delta.ArrowSchemaToDeltaSchemaString(schema)
	if err != nil {
		panic(fmt.Sprintf("failed to convert schema: %v", err))
	}

	fmt.Println("\nGenerated Delta Schema String:")
	fmt.Println(schemaString)

	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("Expected Delta Lake Transaction Log:")

	deltaMetadata, err := delta.GenerateDeltaAppendMetadataFromArrowSchema(
		schema,
		"bbf06f99-2376-4d41-b4fb-367dd31df3de",
		[]string{"YearMonthDate"},
		1749738853893,
		1,
		2,
		"YearMonthDate=20250612/part-00000-f3ffe72a-98fd-4dfe-9e12-96030934ce3f.c000.snappy.parquet",
		map[string]string{"YearMonthDate": "20250612"},
		448876,
		1749738853808,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to generate metadata: %v", err))
	}

	fmt.Println(deltaMetadata)

	// Now generate the append-only transaction notification
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("APPEND-ONLY TRANSACTION NOTIFICATION")
	fmt.Println(strings.Repeat("=", 60))

	transactionNotification, err := delta.GenerateAppendOnlyTransactionNotification(
		"WorkloadIdentityCredential",
		"arcdatasynapsedogfood",
		"onelake",
		"synapse/workspaces/arcdatasynapsedogfood/raw/eventhub/arndatareplicate/armlinkednotifications-uksouth",
		"dfs.core.windows.net",
		"72f988bf-86f1-41af-91ab-2d7cd011db47",
		"DeltaLakeStandaloneDotnet/V1",
		schema,
		"bbf06f99-2376-4d41-b4fb-367dd31df3de",
		[]string{"YearMonthDate"},
		1749738853893,
		1,
		2,
		"YearMonthDate=20250612/part-00000-f3ffe72a-98fd-4dfe-9e12-96030934ce3f.c000.snappy.parquet",
		map[string]string{"YearMonthDate": "20250612"},
		448876,
		1749738853808,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to generate transaction notification: %v", err))
	}

	fmt.Println("Generated Transaction Notification (compact JSON):")
	fmt.Println(transactionNotification)

	// Pretty print the transaction notification for readability
	var transactionPrettyJSON map[string]interface{}
	json.Unmarshal([]byte(transactionNotification), &transactionPrettyJSON)
	transactionPrettyBytes, _ := json.MarshalIndent(transactionPrettyJSON, "", "  ")
	fmt.Println("\nPretty formatted Transaction Notification:")
	fmt.Println(string(transactionPrettyBytes))

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Demo completed successfully!")
	fmt.Println("The serializedTransactionCommandListBase64 field contains the base64-encoded")
	fmt.Println("Delta Lake transaction log that was shown above.")
}

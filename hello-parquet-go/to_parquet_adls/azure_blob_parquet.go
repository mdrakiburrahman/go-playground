package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/apache/arrow/go/v13/parquet"
	"github.com/apache/arrow/go/v13/parquet/compress"
	"github.com/apache/arrow/go/v13/parquet/pqarrow"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: azure_blob_parquet_streaming <accountName> <containerName>")
		return
	}
	accountName := os.Args[1]
	containerName := os.Args[2]
	blobName := "flat_record_compressed.parquet"

	// Create Azure credentials using Azure CLI
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create Azure CLI credential: %v", err))
	}

	// Create blob service client
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	client, err := azblob.NewClient(serviceURL, cred, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create blob client: %v", err))
	}

	// Create Arrow records
	var records []arrow.Record
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "archer", Type: arrow.BinaryTypes.String},
			{Name: "location", Type: arrow.BinaryTypes.String},
			{Name: "year", Type: arrow.PrimitiveTypes.Int16},
		},
		nil,
	)

	rb := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer rb.Release()

	for i := 0; i < 3; i++ {
		postfix := strconv.Itoa(i)
		rb.Field(0).(*array.StringBuilder).AppendValues([]string{"tony" + postfix, "amy" + postfix, "jim" + postfix}, nil)
		rb.Field(1).(*array.StringBuilder).AppendValues([]string{"beijing" + postfix, "shanghai" + postfix, "chengdu" + postfix}, nil)
		rb.Field(2).(*array.Int16Builder).AppendValues([]int16{1992 + int16(i), 1993 + int16(i), 1994 + int16(i)}, nil)
		rec := rb.NewRecord()
		records = append(records, rec)
	}

	// Use bytes.Buffer as io.Writer to write parquet data in memory first
	buf := &bytes.Buffer{}

	// Create parquet writer with compression
	props := parquet.NewWriterProperties(
		parquet.WithCompression(compress.Codecs.Zstd),
		parquet.WithCompressionFor("year", compress.Codecs.Brotli),
	)

	writer, err := pqarrow.NewFileWriter(schema, buf, props, pqarrow.DefaultWriterProps())
	if err != nil {
		panic(fmt.Sprintf("failed to create parquet writer: %v", err))
	}

	// Write records to parquet
	for _, rec := range records {
		if err := writer.Write(rec); err != nil {
			panic(fmt.Sprintf("failed to write record: %v", err))
		}
		rec.Release()
	}

	// Close the writer to finalize the parquet file
	if err := writer.Close(); err != nil {
		panic(fmt.Sprintf("failed to close writer: %v", err))
	}

	// Upload the parquet data to Azure Blob Storage
	ctx := context.Background()
	_, err = client.UploadBuffer(ctx, containerName, blobName, buf.Bytes(), &azblob.UploadBufferOptions{
		BlockSize:   4 * 1024 * 1024, // 4MB blocks
		Concurrency: 16,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to upload to blob storage: %v", err))
	}

	fmt.Printf("Successfully uploaded parquet file to Azure Blob Storage: %s/%s\n", containerName, blobName)
}

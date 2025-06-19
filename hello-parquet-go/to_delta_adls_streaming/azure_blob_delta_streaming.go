package main

import (
	"context"
	"fmt"
	"os"
	"parquet-project/delta"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/apache/arrow/go/v13/parquet"
	"github.com/apache/arrow/go/v13/parquet/compress"
	"github.com/apache/arrow/go/v13/parquet/pqarrow"
)

// BlobWriteCloser implements io.WriteCloser for Azure Blob Storage
type BlobWriteCloser struct {
	client        *azblob.Client
	containerName string
	blobName      string
	ctx           context.Context
	buffer        []byte
}

func NewBlobWriteCloser(client *azblob.Client, containerName, blobName string) *BlobWriteCloser {
	return &BlobWriteCloser{
		client:        client,
		containerName: containerName,
		blobName:      blobName,
		ctx:           context.Background(),
		buffer:        make([]byte, 0),
	}
}

func (bw *BlobWriteCloser) Write(p []byte) (n int, err error) {
	bw.buffer = append(bw.buffer, p...)
	return len(p), nil
}

func (bw *BlobWriteCloser) Close(kustoDatabase string, kustoTable string) error {
	metadata := map[string]*string{
		"rawSizeBytes":    stringPtr(fmt.Sprintf("%d", len(bw.buffer))),
		"kustoDatabase":   stringPtr(kustoDatabase),
		"kustoTable":      stringPtr(kustoTable),
		"kustoDataFormat": stringPtr("parquet"),
	}

	_, err := bw.client.UploadBuffer(bw.ctx, bw.containerName, bw.blobName, bw.buffer, &azblob.UploadBufferOptions{
		BlockSize:   4 * 1024 * 1024, // 4MB blocks
		Concurrency: 16,
		Metadata:    metadata,
	})
	return err
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	if len(os.Args) < 7 {
		fmt.Println("Usage: azure_blob_parquet_streaming <accountName> <containerName> <folderName> <kustoDatabase> <kustoTable> <eventHubNamespace>")
		return
	}
	accountName := os.Args[1]
	containerName := os.Args[2]
	folderName := os.Args[3]
	kustoDatabase := os.Args[4]
	kustoTable := os.Args[5]
	eventHubNamespace := os.Args[6]

	yearMonthDate := time.Now().Format("20060102")

	rootFolderPath := fmt.Sprintf(
		"warehouse/%s",
		folderName,
	)
	blobName := fmt.Sprintf(
		"%s/YearMonthDate=%s/flat_record_compressed_streaming_%d.parquet",
		rootFolderPath,
		yearMonthDate,
		time.Now().Unix(),
	)

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

	// Create blob writer
	blobWriter := NewBlobWriteCloser(client, containerName, blobName)
	defer blobWriter.Close(kustoDatabase, kustoTable)

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

	// Create parquet writer that streams directly to Azure Blob
	props := parquet.NewWriterProperties(
		parquet.WithCompression(compress.Codecs.Zstd),
		parquet.WithCompressionFor("year", compress.Codecs.Brotli),
	)

	writer, err := pqarrow.NewFileWriter(schema, blobWriter, props, pqarrow.DefaultWriterProps())
	if err != nil {
		panic(fmt.Sprintf("failed to create parquet writer: %v", err))
	}

	// Write records to parquet (streaming to blob)
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

	fmt.Printf("Successfully streamed parquet file to Azure Blob Storage: %s/%s\n", containerName, blobName)

	// Generate the transaction notification
	transactionNotification, err := delta.GenerateAppendOnlyTransactionNotification(
		"ManagedIdentityCredential",
		accountName,
		containerName,
		rootFolderPath,
		"dfs.core.windows.net",
		"72f988bf-86f1-41af-91ab-2d7cd011db47",
		"DeltaLakeStandaloneDotnet/V1",
		schema,
		uuid.New().String(),
		[]string{"YearMonthDate"},
		time.Now().UnixMilli(),
		1,
		2,
		blobName,
		map[string]string{
			"YearMonthDate": yearMonthDate,
		},
		int64(len(blobWriter.buffer)),
		time.Now().UnixMilli(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to generate transaction notification: %v", err))
	}

	fmt.Println("\nSending transaction notification to Event Hub...")
	err = sendToEventHub(cred, transactionNotification, eventHubNamespace)
	if err != nil {
		panic(fmt.Sprintf("failed to send to Event Hub: %v", err))
	}
	fmt.Println("Successfully sent transaction notification to Event Hub")
}

// sendToEventHub sends the transaction notification to the specified Event Hub
func sendToEventHub(cred *azidentity.AzureCLICredential, message string, namespace string) error {

	eventHubNamespace := fmt.Sprintf(
		"%s.servicebus.windows.net",
		namespace,
	)
	eventHubName := "delta-bulk-loader"

	// Create Event Hub producer client
	producerClient, err := azeventhubs.NewProducerClient(eventHubNamespace, eventHubName, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create Event Hub producer client: %w", err)
	}
	defer producerClient.Close(context.Background())

	// Create event data
	eventData := &azeventhubs.EventData{
		Body: []byte(message),
		Properties: map[string]any{
			"ContentType": "application/json",
			"MessageType": "DeltaLakeTransactionNotification",
			"Timestamp":   time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Create event batch
	batch, err := producerClient.NewEventDataBatch(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to create event batch: %w", err)
	}

	// Add event to batch
	err = batch.AddEventData(eventData, nil)
	if err != nil {
		return fmt.Errorf("failed to add event to batch: %w", err)
	}

	// Send the batch
	err = producerClient.SendEventDataBatch(context.Background(), batch, nil)
	if err != nil {
		return fmt.Errorf("failed to send event batch: %w", err)
	}

	return nil
}

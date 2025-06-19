package delta

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apache/arrow/go/v13/arrow"
)

// DeltaMetadata represents the metadata action in a Delta Lake transaction log
type DeltaMetadata struct {
	MetaData DeltaMetaData `json:"metaData"`
}

// DeltaMetaData contains the table metadata
type DeltaMetaData struct {
	ID               string            `json:"id"`
	Format           DeltaFormat       `json:"format"`
	SchemaString     string            `json:"schemaString"`
	PartitionColumns []string          `json:"partitionColumns"`
	CreatedTime      int64             `json:"createdTime"`
	Configuration    map[string]string `json:"configuration"`
}

// DeltaFormat specifies the storage format
type DeltaFormat struct {
	Provider string            `json:"provider"`
	Options  map[string]string `json:"options"`
}

// DeltaProtocol represents the protocol action in a Delta Lake transaction log
type DeltaProtocol struct {
	Protocol DeltaProtocolInfo `json:"protocol"`
}

// DeltaProtocolInfo contains protocol version information
type DeltaProtocolInfo struct {
	MinReaderVersion int      `json:"minReaderVersion"`
	MinWriterVersion int      `json:"minWriterVersion"`
	ReaderFeatures   []string `json:"readerFeatures"`
	WriterFeatures   []string `json:"writerFeatures"`
}

// DeltaAdd represents the add action in a Delta Lake transaction log
type DeltaAdd struct {
	Add DeltaAddInfo `json:"add"`
}

// DeltaAddInfo contains information about a file being added
type DeltaAddInfo struct {
	Path             string            `json:"path"`
	PartitionValues  map[string]string `json:"partitionValues"`
	Size             int64             `json:"size"`
	ModificationTime int64             `json:"modificationTime"`
	DataChange       bool              `json:"dataChange"`
	Tags             map[string]string `json:"tags"`
}

// AppendFileGeneratorConfig contains all the parameters needed to generate Delta Lake append metadata
type AppendFileGeneratorConfig struct {
	// Metadata fields
	ID               string
	SchemaString     string
	PartitionColumns []string
	CreatedTime      int64

	// Protocol fields
	MinReaderVersion int
	MinWriterVersion int

	// Add action fields
	Path             string
	PartitionValues  map[string]string
	Size             int64
	ModificationTime int64
}

// GenerateDeltaAppendMetadata generates a complete Delta Lake append transaction log entry
func GenerateDeltaAppendMetadata(config AppendFileGeneratorConfig) (string, error) {
	// Create metadata action
	metadata := DeltaMetadata{
		MetaData: DeltaMetaData{
			ID: config.ID,
			Format: DeltaFormat{
				Provider: "parquet",
				Options:  map[string]string{},
			},
			SchemaString:     config.SchemaString,
			PartitionColumns: config.PartitionColumns,
			CreatedTime:      config.CreatedTime,
			Configuration:    map[string]string{},
		},
	}

	// Create protocol action
	protocol := DeltaProtocol{
		Protocol: DeltaProtocolInfo{
			MinReaderVersion: config.MinReaderVersion,
			MinWriterVersion: config.MinWriterVersion,
			ReaderFeatures:   nil,
			WriterFeatures:   nil,
		},
	}

	// Create add action
	add := DeltaAdd{
		Add: DeltaAddInfo{
			Path:             config.Path,
			PartitionValues:  config.PartitionValues,
			Size:             config.Size,
			ModificationTime: config.ModificationTime,
			DataChange:       true,
			Tags:             map[string]string{},
		},
	}

	// Marshal each action to JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	protocolJSON, err := json.Marshal(protocol)
	if err != nil {
		return "", fmt.Errorf("failed to marshal protocol: %w", err)
	}

	addJSON, err := json.Marshal(add)
	if err != nil {
		return "", fmt.Errorf("failed to marshal add action: %w", err)
	}

	// Combine all actions with newlines (Delta Lake transaction log format)
	result := strings.Join([]string{
		string(metadataJSON),
		string(protocolJSON),
		string(addJSON),
	}, "\n")

	return result, nil
}

// GenerateDeltaAppendMetadataFromArrowSchema is a convenience function that creates Delta Lake metadata from an Arrow schema
func GenerateDeltaAppendMetadataFromArrowSchema(
	schema *arrow.Schema,
	id string,
	partitionColumns []string,
	createdTime int64,
	minReaderVersion int,
	minWriterVersion int,
	path string,
	partitionValues map[string]string,
	size int64,
	modificationTime int64,
) (string, error) {
	// Convert Arrow schema to Delta Lake schema string
	schemaString, err := ArrowSchemaToDeltaSchemaString(schema)
	if err != nil {
		return "", fmt.Errorf("failed to convert Arrow schema: %w", err)
	}

	// Create config
	config := AppendFileGeneratorConfig{
		ID:               id,
		SchemaString:     schemaString,
		PartitionColumns: partitionColumns,
		CreatedTime:      createdTime,
		MinReaderVersion: minReaderVersion,
		MinWriterVersion: minWriterVersion,
		Path:             path,
		PartitionValues:  partitionValues,
		Size:             size,
		ModificationTime: modificationTime,
	}

	return GenerateDeltaAppendMetadata(config)
}

// NotifyTransaction sends a transaction notification to the Delta Lake transaction log
func NotifyTransaction(
	deltaLogPath string,
	protocolVersion DeltaProtocolInfo,
	metadata DeltaMetaData,
	addFiles []DeltaAddInfo,
) error {
	// Create transaction log entry
	var actions []interface{}

	// Add metadata action
	actions = append(actions, metadata)

	// Add protocol action
	actions = append(actions, DeltaProtocol{Protocol: protocolVersion})

	// Add add actions
	for _, addFile := range addFiles {
		actions = append(actions, DeltaAdd{Add: addFile})
	}

	// Marshal each action to JSON
	var jsonActions []string
	for _, action := range actions {
		actionJSON, err := json.Marshal(action)
		if err != nil {
			return fmt.Errorf("failed to marshal action: %w", err)
		}
		jsonActions = append(jsonActions, string(actionJSON))
	}

	// Combine all actions with newlines (Delta Lake transaction log format)
	logEntry := strings.Join(jsonActions, "\n")

	// Encode the log entry in base64
	encodedLogEntry := base64.StdEncoding.EncodeToString([]byte(logEntry))

	// Here you would write the encodedLogEntry to the Delta Lake transaction log file
	// For example purposes, we just print it
	fmt.Println("Transaction Log Entry (base64):", encodedLogEntry)

	return nil
}

// TransactionDestination represents the destination configuration for the transaction
type TransactionDestination struct {
	StorageAccountAuthType    string `json:"storageAccountAuthType"`
	StorageAccountName        string `json:"storageAccountName"`
	StorageContainerName      string `json:"storageContainerName"`
	TableRelativePath         string `json:"tableRelativePath"`
	StorageAccountDfsEndpoint string `json:"storageAccountDfsEndpoint"`
	StorageAccountTenantId    string `json:"storageAccountTenantId"`
	EngineInfo                string `json:"engineInfo"`
}

// AppendOnlyTransactionNotification represents the full transaction notification structure
type AppendOnlyTransactionNotification struct {
	TransactionDestination                 TransactionDestination `json:"transactionDestination"`
	SerializedTransactionCommandListBase64 string                 `json:"serializedTransactionCommandListBase64"`
}

// GenerateAppendOnlyTransactionNotification generates a transaction notification with base64 encoded metadata
func GenerateAppendOnlyTransactionNotification(
	storageAccountAuthType string,
	storageAccountName string,
	storageContainerName string,
	tableRelativePath string,
	storageAccountDfsEndpoint string,
	storageAccountTenantId string,
	engineInfo string,
	schema *arrow.Schema,
	id string,
	partitionColumns []string,
	createdTime int64,
	minReaderVersion int,
	minWriterVersion int,
	path string,
	partitionValues map[string]string,
	size int64,
	modificationTime int64,
) (string, error) {
	// Generate the Delta append metadata
	deltaMetadata, err := GenerateDeltaAppendMetadataFromArrowSchema(
		schema,
		id,
		partitionColumns,
		createdTime,
		minReaderVersion,
		minWriterVersion,
		path,
		partitionValues,
		size,
		modificationTime,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate delta metadata: %w", err)
	}

	// Base64 encode the serialized transaction command list
	serializedBase64 := base64.StdEncoding.EncodeToString([]byte(deltaMetadata))

	// Create the transaction destination
	destination := TransactionDestination{
		StorageAccountAuthType:    storageAccountAuthType,
		StorageAccountName:        storageAccountName,
		StorageContainerName:      storageContainerName,
		TableRelativePath:         tableRelativePath,
		StorageAccountDfsEndpoint: storageAccountDfsEndpoint,
		StorageAccountTenantId:    storageAccountTenantId,
		EngineInfo:                engineInfo,
	}

	// Create the full notification
	notification := AppendOnlyTransactionNotification{
		TransactionDestination:                 destination,
		SerializedTransactionCommandListBase64: serializedBase64,
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(notification)
	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction notification to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

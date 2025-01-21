package marshaler

import (
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	// OTLP content types and encodings
	encodingCsv    = "otlp_csv"
	contentTypeCsv = "application/csv"
)

// otlpLogs defines a struct for marshaling logs into bytes using
// internal implementations of plog.Marshaler.
type otlpLogs struct {
	// logsMarshaler is the internal implementation of plog.Marshaler.
	logsMarshaler plog.Marshaler

	// encoding is the name of the encoding that this marshaler supports.
	encoding string

	// contentType is the indicates the general MIME type of the marshaled data.
	contentType string
}

// NewOtlpCsvLogs creates a new otlpLogs that uses csv as the encoding.
func NewOtlpCsvLogs() Logs {
	return &otlpLogs{
		logsMarshaler: &CSVMarshaler{},
		encoding:      encodingCsv,
		contentType:   contentTypeCsv,
	}
}

// Marshal serializes logs into bytes.
func (o *otlpLogs) Marshal(logs plog.Logs) ([]byte, error) {
	return o.logsMarshaler.MarshalLogs(logs)
}

// Encoding is the name of the encoding that this marshaler supports.
func (o *otlpLogs) Encoding() string {
	return o.encoding
}

// ContentType is the indicates the marshaled format of the data.
func (o *otlpLogs) ContentType() string {
	return o.contentType
}

// otlpMetrics defines a struct for marshaling metrics into bytes using
// internal implementations of pmetric.Marshaler.
type otlpMetrics struct {
	metricsMarshaler pmetric.Marshaler
	encoding         string
	contentType      string
}

// NewOtlpCsvMetrics creates a new otlpMetrics that uses csv as the encoding.
func NewOtlpCsvMetrics() Metrics {
	return &otlpMetrics{
		metricsMarshaler: &CSVMarshaler{},
		encoding:         encodingCsv,
		contentType:      contentTypeCsv,
	}
}

// Marshal serializes metrics into bytes.
func (o *otlpMetrics) Marshal(metrics pmetric.Metrics) ([]byte, error) {
	return o.metricsMarshaler.MarshalMetrics(metrics)
}

// Encoding is the name of the encoding that this marshaler supports.
func (o *otlpMetrics) Encoding() string {
	return o.encoding
}

// ContentType is the indicates the marshaled format of the data.
func (o *otlpMetrics) ContentType() string {
	return o.contentType
}

// otlpTraces defines a struct for marshaling traces into bytes using
// internal implementations of ptrace.Marshaler.
type otlpTraces struct {
	tracesMarshaler ptrace.Marshaler
	encoding        string
	contentType     string
}

// NewOtlpCsvTraces creates a new otlpTraces that uses csv as the encoding.
func NewOtlpCsvTraces() Traces {
	return &otlpTraces{
		tracesMarshaler: &CSVMarshaler{},
		encoding:        encodingCsv,
		contentType:     contentTypeCsv,
	}
}

// Marshal serializes traces into bytes.
func (o *otlpTraces) Marshal(traces ptrace.Traces) ([]byte, error) {
	return o.tracesMarshaler.MarshalTraces(traces)
}

// Encoding is the name of the encoding that this marshaler supports.
func (o *otlpTraces) Encoding() string {
	return o.encoding
}

// ContentType is the indicates the marshaled format of the data.
func (o *otlpTraces) ContentType() string {
	return o.contentType
}

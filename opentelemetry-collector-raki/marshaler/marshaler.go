package marshaler

import (
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// Logs defines an interface for Marshaling logs into bytes.
type Logs interface {
	// Marshal serializes logs into bytes.
	Marshal(logs plog.Logs) ([]byte, error)

	// Encoding is the name of the encoding that this marshaler supports.
	Encoding() string

	// ContentType is the indicates the original format of the data
	// before encodings e.g. x-protobuf
	ContentType() string
}

// Metrics defines an interface for Marshaling metrics into bytes.
type Metrics interface {
	// Marshal serializes metrics into bytes.
	Marshal(metrics pmetric.Metrics) ([]byte, error)

	// Encoding is the name of the encoding that this marshaler supports.
	Encoding() string

	// ContentType is the indicates the original format of the data
	// before encodings e.g. x-protobuf
	ContentType() string
}

// Traces defines an interface for Marshaling traces into bytes.
type Traces interface {
	// Marshal serializes traces into bytes.
	Marshal(traces ptrace.Traces) ([]byte, error)

	// Encoding is the name of the encoding that this marshaler supports.
	Encoding() string

	// ContentType is the indicates the original format of the data
	// before encodings e.g. x-protobuf
	ContentType() string
}

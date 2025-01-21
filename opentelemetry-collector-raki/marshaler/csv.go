package marshaler // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awss3exporter"

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type CSVMarshaler struct{}

func NewCSVMarshaler() CSVMarshaler {
	return CSVMarshaler{}
}

// MarshalLogs converts OpenTelemetry logs into a CSV format.
func (CSVMarshaler) MarshalLogs(ld plog.Logs) ([]byte, error) {
	buf := bytes.Buffer{}
	writer := csv.NewWriter(&buf)

	// Write CSV header
	header := []string{"timestamp", "severity", "body", "attributes"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rl := rls.At(i)
		ills := rl.ScopeLogs()
		for j := 0; j < ills.Len(); j++ {
			ils := ills.At(j)
			logs := ils.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)

				// Extract fields for CSV
				timestamp := lr.Timestamp().AsTime().Format("2006-01-02T15:04:05Z")
				severity := lr.SeverityText()
				body := lr.Body().AsString()
				attributes, err := attributesToCSVString(lr.Attributes())
				if err != nil {
					return nil, fmt.Errorf("failed to serialize attributes: %w", err)
				}

				// Write log entry as a CSV row
				record := []string{timestamp, severity, body, attributes}
				if err := writer.Write(record); err != nil {
					return nil, fmt.Errorf("failed to write CSV record: %w", err)
				}
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error during CSV writing: %w", err)
	}

	return buf.Bytes(), nil
}

// UnmarshalLogs converts a CSV byte array into OpenTelemetry logs.
func (CSVMarshaler) UnmarshalLogs(buf []byte) (plog.Logs, error) {
	reader := csv.NewReader(bytes.NewReader(buf))
	lines, err := reader.ReadAll()
	if err != nil {
		return plog.NewLogs(), fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(lines) < 2 {
		return plog.NewLogs(), errors.New("no log records found")
	}

	ld := plog.NewLogs()
	for _, line := range lines[1:] {
		if len(line) != 4 {
			return plog.NewLogs(), errors.New("invalid CSV format")
		}
		logRecord := ld.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
		timestamp, err := time.Parse("2006-01-02T15:04:05Z", line[0])
		if err != nil {
			return plog.NewLogs(), fmt.Errorf("failed to parse timestamp: %w", err)
		}
		logRecord.SetTimestamp(pcommon.NewTimestampFromTime(timestamp))
		logRecord.SetSeverityText(line[1])
		logRecord.Body().SetStr(line[2])
		logRecord.Attributes().FromRaw(attributesFromCSVString(line[3]).AsRaw())
	}

	return ld, nil
}

// MarshalMetrics converts OpenTelemetry metrics into a CSV format.
func (CSVMarshaler) MarshalMetrics(md pmetric.Metrics) ([]byte, error) {
	buf := bytes.Buffer{}
	writer := csv.NewWriter(&buf)
	header := []string{"timestamp", "metric_name", "value"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	metrics := md.ResourceMetrics()
	for i := 0; i < metrics.Len(); i++ {
		rm := metrics.At(i)
		ilms := rm.ScopeMetrics()
		for j := 0; j < ilms.Len(); j++ {
			ilm := ilms.At(j)
			metrics := ilm.Metrics()
			for k := 0; k < metrics.Len(); k++ {
				m := metrics.At(k)
				switch m.Type() {
				case pmetric.MetricTypeGauge:
					gauge := m.Gauge().DataPoints()
					for l := 0; l < gauge.Len(); l++ {
						dp := gauge.At(l)
						timestamp := dp.Timestamp().AsTime().Format("2006-01-02T15:04:05Z")
						metricName := m.Name()
						value := fmt.Sprintf("%f", dp.DoubleValue())
						record := []string{timestamp, metricName, value}
						if err := writer.Write(record); err != nil {
							return nil, fmt.Errorf("failed to write CSV record: %w", err)
						}
					}
				}
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error during CSV writing: %w", err)
	}

	return buf.Bytes(), nil
}

// UnmarshalMetrics converts a CSV byte array into OpenTelemetry metrics.
func (CSVMarshaler) UnmarshalMetrics(buf []byte) (pmetric.Metrics, error) {
	return pmetric.NewMetrics(), errors.New("unmarshaling metrics is not implemented")
}

// MarshalTraces converts OpenTelemetry traces into a CSV format.
func (CSVMarshaler) MarshalTraces(td ptrace.Traces) ([]byte, error) {
	buf := bytes.Buffer{}
	writer := csv.NewWriter(&buf)
	header := []string{"trace_id", "span_id", "name", "status"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	traces := td.ResourceSpans()
	for i := 0; i < traces.Len(); i++ {
		rspan := traces.At(i)
		illSpans := rspan.ScopeSpans()
		for j := 0; j < illSpans.Len(); j++ {
			span := illSpans.At(j).Spans()
			for k := 0; k < span.Len(); k++ {
				s := span.At(k)
				record := []string{
					s.TraceID().String(),
					s.SpanID().String(),
					s.Name(),
					s.Status().Message(),
				}
				if err := writer.Write(record); err != nil {
					return nil, fmt.Errorf("failed to write CSV record: %w", err)
				}
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error during CSV writing: %w", err)
	}

	return buf.Bytes(), nil
}

// UnmarshalTraces converts a CSV byte array into OpenTelemetry traces.
func (CSVMarshaler) UnmarshalTraces(buf []byte) (ptrace.Traces, error) {
	return ptrace.NewTraces(), errors.New("unmarshaling traces is not implemented")
}

// attributesFromCSVString parses a semicolon-separated key-value string and returns a map.
func attributesFromCSVString(csvStr string) pcommon.Map {
	attrs := pcommon.NewMap()
	if csvStr == "" {
		return attrs
	}

	pairs := strings.Split(csvStr, ";")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue // Skip malformed pairs
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		attrs.PutStr(key, value)
	}

	return attrs
}

// attributesToCSVString serializes attributes into a single string for CSV.
func attributesToCSVString(attrs pcommon.Map) (string, error) {
	var sb strings.Builder
	first := true
	attrs.Range(func(k string, v pcommon.Value) bool {
		if !first {
			sb.WriteString("; ")
		}
		first = false
		attrValue, err := attributeValueToString(v)
		if err != nil {
			return false
		}
		sb.WriteString(fmt.Sprintf("%s=%s", k, attrValue))
		return true
	})

	if !first {
		return sb.String(), nil
	}
	return "", nil
}

// attributeValueToString converts a pcommon.Value to a string representation.
func attributeValueToString(v pcommon.Value) (string, error) {
	switch v.Type() {
	case pcommon.ValueTypeStr:
		return v.Str(), nil
	case pcommon.ValueTypeBool:
		return fmt.Sprintf("%t", v.Bool()), nil
	case pcommon.ValueTypeInt:
		return fmt.Sprintf("%d", v.Int()), nil
	case pcommon.ValueTypeDouble:
		return fmt.Sprintf("%f", v.Double()), nil
	default:
		return "", errors.New("unsupported attribute type")
	}
}

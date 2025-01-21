package marshaler

// Marshalers is a collection of marshalers for logs, metrics, and traces
type Marshalers struct {
	Logs    map[string]Logs
	Metrics map[string]Metrics
	Traces  map[string]Traces
}

// baseMarshalers returns the set of supported marshalers
func BaseMarshalers() *Marshalers {
	return &Marshalers{
		Logs:    BaseLogsMarshalers(),
		Metrics: BaseMetricsMarshalers(),
		Traces:  BaseTracesMarshalers(),
	}
}

// baseLogsMarshalers returns the set of supported logs marshalers
func BaseLogsMarshalers() map[string]Logs {
	otlpCsv := NewOtlpCsvLogs()
	return map[string]Logs{
		otlpCsv.Encoding(): otlpCsv,
	}
}

// baseMetricsMarshalers returns the set of supported metrics marshalers
func BaseMetricsMarshalers() map[string]Metrics {
	otlpCsv := NewOtlpCsvMetrics()
	return map[string]Metrics{
		otlpCsv.Encoding(): otlpCsv,
	}
}

// baseTracesMarshalers returns the set of supported traces marshalers
func BaseTracesMarshalers() map[string]Traces {
	otlpCsv := NewOtlpCsvTraces()
	return map[string]Traces{
		otlpCsv.Encoding(): otlpCsv,
	}
}

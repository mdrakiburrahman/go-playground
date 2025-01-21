package emptyexporter

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-tutorials/marshaler"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

var (
	typeStr = component.MustNewType("emptyexporter")
)

type emptyExporterFactory struct {
	Marshalers *marshaler.Marshalers
	sender     senderFunc
}

// FactoryOption is used to configure a factory.
type FactoryOption func(*emptyExporterFactory)

// WithMarshalers adds additional marshalers to the factory or overrides
// existing marshalers if present.
func WithLogsMarshalers(marshalers ...marshaler.Logs) FactoryOption {
	return func(f *emptyExporterFactory) {
		for _, m := range marshalers {
			f.Marshalers.Logs[m.Encoding()] = m
		}
	}
}

// WithMetricsMarshalers adds additional marshalers to the factory or overrides
// existing marshalers if present.
func WithMetricsMarshalers(marshalers ...marshaler.Metrics) FactoryOption {
	return func(f *emptyExporterFactory) {
		for _, m := range marshalers {
			f.Marshalers.Metrics[m.Encoding()] = m
		}
	}
}

// WithTracesMarshalers adds additional marshalers to the factory or overrides
// existing marshalers if present.
func WithTracesMarshalers(marshalers ...marshaler.Traces) FactoryOption {
	return func(f *emptyExporterFactory) {
		for _, m := range marshalers {
			f.Marshalers.Traces[m.Encoding()] = m
		}
	}
}

// WithSender overrides the default sender with a custom sender.
func WithSender(sender senderFunc) FactoryOption {
	return func(f *emptyExporterFactory) {
		f.sender = sender
	}
}

func NewFactory(options ...FactoryOption) exporter.Factory {
	f := &emptyExporterFactory{
		Marshalers: marshaler.BaseMarshalers(),
		sender:     nil,
	}

	for _, opt := range options {
		opt(f)
	}
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithTraces(f.createTracesExporter, component.StabilityLevelDevelopment),
		exporter.WithMetrics(f.createMetricsExporter, component.StabilityLevelDevelopment),
		exporter.WithLogs(f.createLogsExporter, component.StabilityLevelDevelopment),
	)
}

func (f *emptyExporterFactory) createTracesExporter(
	ctx context.Context,
	params exporter.Settings,
	config component.Config) (exporter.Traces, error) {
	cfg := config.(*Config)
	s, err := newEmptyexporter(params.Logger, config.(*Config))
	if err != nil {
		return nil, err
	}
	if marshaler, ok := f.Marshalers.Traces[cfg.Encoding]; ok {
		s.registerTracesMarshaler(marshaler)
	} else {
		return nil, fmt.Errorf("marshaler %s not found", cfg.Encoding)
	}
	return exporterhelper.NewTraces(ctx, params, cfg, s.pushTraces)
}

func (f *emptyExporterFactory) createMetricsExporter(
	ctx context.Context,
	params exporter.Settings,
	config component.Config) (exporter.Metrics, error) {
	cfg := config.(*Config)
	s, err := newEmptyexporter(params.Logger, config.(*Config))
	if err != nil {
		return nil, err
	}
	if marshaler, ok := f.Marshalers.Metrics[cfg.Encoding]; ok {
		s.registerMetricsMarshaler(marshaler)
	} else {
		return nil, fmt.Errorf("marshaler %s not found", cfg.Encoding)
	}
	return exporterhelper.NewMetrics(ctx, params, cfg, s.pushMetrics)
}

func (f *emptyExporterFactory) createLogsExporter(
	ctx context.Context,
	params exporter.Settings,
	config component.Config) (exporter.Logs, error) {
	cfg := config.(*Config)
	s, err := newEmptyexporter(params.Logger, config.(*Config))
	if err != nil {
		return nil, err
	}
	if marshaler, ok := f.Marshalers.Logs[cfg.Encoding]; ok {
		s.registerLogsMarshaler(marshaler)
	} else {
		return nil, fmt.Errorf("marshaler %s not found", cfg.Encoding)
	}
	return exporterhelper.NewLogs(ctx, params, cfg, s.pushLogs)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

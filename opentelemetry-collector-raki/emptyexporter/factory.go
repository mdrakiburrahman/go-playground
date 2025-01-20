package emptyexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

var (
	typeStr = component.MustNewType("emptyexporter")
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithTraces(createTracesExporter, component.StabilityLevelDevelopment),
		exporter.WithMetrics(createMetricsExporter, component.StabilityLevelDevelopment),
		exporter.WithLogs(createLogsExporter, component.StabilityLevelDevelopment),
	)
}

func createTracesExporter(
	ctx context.Context,
	params exporter.Settings,
	config component.Config) (exporter.Traces, error) {

	cfg := config.(*Config)
	s, _ := newEmptyexporter(params.Logger, config.(*Config))
	return exporterhelper.NewTraces(ctx, params, cfg, s.pushTraces)
}

func createMetricsExporter(
	ctx context.Context,
	params exporter.Settings,
	config component.Config) (exporter.Metrics, error) {

	cfg := config.(*Config)
	s, _ := newEmptyexporter(params.Logger, config.(*Config))
	return exporterhelper.NewMetrics(ctx, params, cfg, s.pushMetrics)
}

func createLogsExporter(
	ctx context.Context,
	params exporter.Settings,
	config component.Config) (exporter.Logs, error) {

	cfg := config.(*Config)
	s, _ := newEmptyexporter(params.Logger, config.(*Config))
	return exporterhelper.NewLogs(ctx, params, cfg, s.pushLogs)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

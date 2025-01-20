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

type Config struct {
}

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
	set exporter.Settings,
	config component.Config) (exporter.Traces, error) {

	cfg := config.(*Config)
	s := NewEmptyexporter()
	return exporterhelper.NewTraces(ctx, set, cfg, s.pushTraces)
}

func createMetricsExporter(
	ctx context.Context,
	set exporter.Settings,
	config component.Config) (exporter.Metrics, error) {

	cfg := config.(*Config)
	s := NewEmptyexporter()
	return exporterhelper.NewMetrics(ctx, set, cfg, s.pushMetrics)
}

func createLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	config component.Config) (exporter.Logs, error) {

	cfg := config.(*Config)
	s := NewEmptyexporter()
	return exporterhelper.NewLogs(ctx, set, cfg, s.pushLogs)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

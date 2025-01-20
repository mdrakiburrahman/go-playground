package emptyexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type emptyexporter struct {
	config Config
	logger *zap.Logger
}

func newEmptyexporter(logger *zap.Logger, config component.Config) (*emptyexporter, error) {
	return &emptyexporter{
		config: *config.(*Config),
		logger: logger,
	}, nil
}

func (s *emptyexporter) pushLogs(_ context.Context, ld plog.Logs) error {
	if s.config.ShouldLog {
		s.logger.Info("pushLogs", zap.Any("logs", ld))
	}
	return nil
}

func (s *emptyexporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	if s.config.ShouldLog {
		s.logger.Info("pushMetrics", zap.Any("metrics", md))
	}
	return nil
}

func (s *emptyexporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	if s.config.ShouldLog {
		s.logger.Info("pushTraces", zap.Any("traces", td))
	}
	return nil
}

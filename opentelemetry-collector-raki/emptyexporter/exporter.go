package emptyexporter

import (
	"context"

	"github.com/open-telemetry/opentelemetry-tutorials/marshaler"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type senderFunc func(ctx context.Context, e *emptyexporter, content string, contentType string) error

type emptyexporter struct {
	config Config
	logger *zap.Logger

	logsMarshaler    marshaler.Logs
	metricsMarshaler marshaler.Metrics
	tracesMarshaler  marshaler.Traces

	sender senderFunc
}

func newEmptyexporter(logger *zap.Logger, config component.Config) (*emptyexporter, error) {
	return &emptyexporter{
		config: *config.(*Config),
		logger: logger,
		sender: sendToScreen,
	}, nil
}

func (s *emptyexporter) pushLogs(_ context.Context, ld plog.Logs) error {
	bytes, err := s.logsMarshaler.Marshal(ld)
	if err != nil {
		return err
	}
	if s.config.ShouldLog {
		s.sender(context.Background(), s, string(bytes), s.logsMarshaler.ContentType())
	}
	return nil
}

func (s *emptyexporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	bytes, err := s.metricsMarshaler.Marshal(md)
	if err != nil {
		return err
	}
	if s.config.ShouldLog {
		s.sender(ctx, s, string(bytes), s.metricsMarshaler.ContentType())
	}
	return nil
}

func (s *emptyexporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	bytes, err := s.tracesMarshaler.Marshal(td)
	if err != nil {
		return err
	}
	if s.config.ShouldLog {
		s.sender(context.Background(), s, string(bytes), s.tracesMarshaler.ContentType())
	}
	return nil
}

// registerTracesMarshaler sets the traces marshaler to use
func (e *emptyexporter) registerTracesMarshaler(marshaler marshaler.Traces) {
	e.tracesMarshaler = marshaler
}

// registerMetricsMarshaler sets the metrics marshaler to use
func (e *emptyexporter) registerMetricsMarshaler(marshaler marshaler.Metrics) {
	e.metricsMarshaler = marshaler
}

// registerLogsMarshaler sets the logs marshaler to use
func (e *emptyexporter) registerLogsMarshaler(marshaler marshaler.Logs) {
	e.logsMarshaler = marshaler
}

func sendToScreen(_ context.Context, e *emptyexporter, content string, contentType string) error {
	e.logger.Info("Empty send ->", zap.String("content", content), zap.String("contentType", contentType))
	return nil
}

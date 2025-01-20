package emptyexporter

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type emptyexporter struct {
}

func NewEmptyexporter() *emptyexporter {
	return &emptyexporter{}
}

func (s *emptyexporter) pushLogs(_ context.Context, ld plog.Logs) error {
	return nil
}

func (s *emptyexporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	return nil
}

func (s *emptyexporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	return nil
}

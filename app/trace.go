package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/trace/datadog"
	"cto-github.cisco.com/NFV-BU/go-msx/trace/jaeger"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, registerTracers)
	OnEvent(EventConfigure, PhaseAfter, activateTracing)
	OnEvent(EventFinal, PhaseAfter, deactivateTracing)
}

func registerTracers(ctx context.Context) error {
	return types.ErrorList{
		datadog.RegisterTracer(ctx),
		jaeger.RegisterTracer(ctx),
	}.Filter()
}

func activateTracing(ctx context.Context) error {
	logger.Info("Activating tracing")
	return trace.ConfigureTracer(ctx)
}

func deactivateTracing(ctx context.Context) error {
	logger.Info("Deactivating tracing")
	return trace.ShutdownTracer(ctx)
}

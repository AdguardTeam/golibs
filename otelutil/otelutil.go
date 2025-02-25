// Package otelutil contains common utilities for working with OpenTelemetry.
package otelutil

import (
	"context"
	"fmt"
	"os"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/service"
	"github.com/AdguardTeam/golibs/validate"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
)

// Supported OpenTelemetry exporter protocols that are expected to be present in
// the environment under [EnvExporterProto].
const (
	ExporterProtoGRPC         = "grpc"
	ExporterProtoHTTPProtobuf = "http/protobuf"
	ExporterProtoStdout       = "stdout"
)

// EnvExporterProto is the name of the environment variable holding the
// OpenTelemetry exporter protocol.
//
// See https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/.
const EnvExporterProto = "OTEL_EXPORTER_OTLP_PROTOCOL"

// Config is the configuration structure for the OpenTelemetry infrastructure.
type Config struct {
	// ServiceName is the name to use for exporting traces.  It must not be
	// empty.
	ServiceName string
}

// type check
var _ validate.Interface = (*Config)(nil)

// Validate implements the [validate.Interface] for *Config.  c may be nil.
func (c *Config) Validate() (err error) {
	if c == nil {
		return errors.ErrNoValue
	}

	return validate.NotEmpty("ServiceName", c.ServiceName)
}

// Init initializes the global OpenTelemetry infrastructure.  svc is the service
// that should be shut down on exit.  c must be valid.
//
// TODO(a.garipov):  See if there are ways to not use globals.
func Init(ctx context.Context, c *Config) (svc service.Interface, err error) {
	err = errors.Annotate(c.Validate(), "c: %w")
	errors.Check(err)

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
	)

	otel.SetTextMapPropagator(prop)

	var exporter trace.SpanExporter
	proto := os.Getenv(EnvExporterProto)
	switch proto {
	case "", ExporterProtoStdout:
		exporter, err = stdouttrace.New()
	case ExporterProtoGRPC:
		exporter, err = otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithDialOption(grpc.WithDisableServiceConfig()),
		)
	case ExporterProtoHTTPProtobuf:
		exporter, err = otlptracehttp.New(ctx)
	default:
		return nil, fmt.Errorf("%s: %w: %q", EnvExporterProto, errors.ErrBadEnumValue, proto)
	}
	if err != nil {
		return nil, fmt.Errorf("creating otel exporter: %s: %w", proto, err)
	}

	// Use the provided service name, since the env one may be incorrect.
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(c.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating otel resource: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
		trace.WithSampler(trace.ParentBased(trace.NeverSample())),
	)

	svc = service.NewShutdownService(tp)

	otel.SetTracerProvider(tp)

	return svc, nil
}

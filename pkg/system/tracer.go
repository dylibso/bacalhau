package system

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/filecoin-project/bacalhau/pkg/logger"
	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// ----------------------------------------
// Tracer helpers
// ----------------------------------------

func GetTracer() oteltrace.Tracer {
	return otel.GetTracerProvider().Tracer("bacalhau")
}

// ----------------------------------------
// Span helpers
// ----------------------------------------

func NewRootSpan(ctx context.Context, t oteltrace.Tracer, name string) (context.Context, oteltrace.Span) {
	// Always include environment info in spans:
	m0, _ := baggage.NewMember("environment", GetEnvironment().String())
	b, _ := baggage.New(m0)
	ctx = baggage.ContextWithBaggage(ctx, b)

	return t.Start(ctx, name)
}

func GetSpanFromRequest(req *http.Request, name string) (context.Context, oteltrace.Span) {
	ctx := req.Context()
	ctx, span := GetTracer().Start(ctx, name)
	return ctx, span
}

// Span creates and starts a new span, and a context containing it.
// For more information see the otel.Tracer.Start(...) docs:
// https://pkg.go.dev/go.opentelemetry.io/otel/trace#Tracer
// ctx: the context to use for the span
// tracerName: the name of the service that the span is for - will be prefixed with "tracer/".
// Will create a new one if one with the same name does not exist
// spanName: the name of the span, inside the service
// opts: additional options to configure the span from trace.SpanStartOption
func Span(ctx context.Context, tracerName, spanName string,
	opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	// Always include environment info in spans:
	opts = append(opts, oteltrace.WithAttributes(
		attribute.String("environment", GetEnvironment().String()),
	))

	spanName = fmt.Sprintf("service/%s", spanName)

	return GetTracer().Start(ctx, spanName, opts...)
}

// ----------------------------------------
// Baggage and Attribute helpers
// ----------------------------------------

func AddNodeIDToBaggage(ctx context.Context, nodeID string) context.Context {
	return addFieldToBaggage(ctx, model.TracerAttributeNameNodeID, nodeID)
}

func AddJobIDToBaggage(ctx context.Context, jobID string) context.Context {
	return addFieldToBaggage(ctx, model.TracerAttributeNameJobID, jobID)
}

func addFieldToBaggage(ctx context.Context, key, value string) context.Context {
	b := baggage.FromContext(ctx)
	m, err := baggage.NewMember(key, value)
	if err != nil {
		log.Warn().Msgf("failed to add key %s to baggage: %s", key, err)
	}

	b, err = b.SetMember(m)
	if err != nil {
		log.Warn().Msgf("failed to add baggage member to baggage: %s", err)
	}

	return baggage.ContextWithBaggage(ctx, b)
}

func AddJobIDFromBaggageToSpan(ctx context.Context, span oteltrace.Span) {
	AddAttributeToSpanFromBaggage(ctx, span, model.TracerAttributeNameJobID)
}

func AddNodeIDFromBaggageToSpan(ctx context.Context, span oteltrace.Span) {
	AddAttributeToSpanFromBaggage(ctx, span, model.TracerAttributeNameNodeID)
}

func AddAttributeToSpanFromBaggage(ctx context.Context, span oteltrace.Span, name string) {
	b := baggage.FromContext(ctx)
	log.Trace().Msgf("adding %s from baggage to span as attribute: %+v", name, b)
	m := b.Member(name)
	if m.Value() != "" {
		span.SetAttributes(attribute.String(name, m.Value()))
	} else {
		log.Trace().Err(errors.WithStack(errors.New("missing value"))).
			Str("baggage_key", name).Msg("No value found for baggage key")
	}
}

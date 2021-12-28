package bifrost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type tracerContext struct {
	Name string
}

func (r *tracerContext) String() string {
	return "context value " + r.Name
}

var TracerContext = tracerContext{Name: "context Respond"}

func HttpTracer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		operation := r.Method + " " + r.URL.Path
		opts := append(
			[]trace.SpanStartOption{
				trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
				trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
				trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(operation, "", r)...),
			},
		) // start with the configured options

		carrier := http.Header{}
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(carrier))
		tr := otel.Tracer("http.tracer")
		ctxSpan, span := tr.Start(ctx, operation, opts...)
		defer span.End()

		defer func() {
			if err := recover(); err != nil {
				span.SetStatus(codes.Error, "recover")
				span.RecordError(err.(error))
				span.SetAttributes(attribute.Key("event").String("error"))
				span.SetAttributes(attribute.Key("error.kind").String("panic"))
				span.SetAttributes(attribute.Key("stack").String(string(debug.Stack())))

				span.End()
				panic(err)
			}
		}()

		span.SetAttributes(attribute.String("request.id", r.Header.Get("X-Request-Id")))

		// There's nothing we can do with any errors here.
		otel.GetTextMapPropagator().Inject(ctxSpan, propagation.HeaderCarrier(carrier))
		r = r.WithContext(trace.ContextWithSpan(ctxSpan, span))

		// check content length
		if r.ContentLength > 0 {
			// Request
			var buf []byte
			if r.Body != nil { // Read
				buf, _ = ioutil.ReadAll(r.Body)
			}

			response := make(map[string]interface{})

			// get content-type
			cType := r.Header.Get(HeaderContentType)

			switch {
			case strings.HasPrefix(cType, MIMETextPlain):
			case strings.HasPrefix(cType, MIMEApplicationForm):
				if err := r.ParseForm(); err != nil {
					log.Error().Err(ErrBadRequest(w, r, err)).Msg("Request body contains badly-formed form-urlencoded")
					_ = ResponseJSONPayload(w, r, http.StatusBadRequest, nil)
					return
				}

				log.Info().
					Str(HeaderContentType, cType).
					Str("body", string(buf)).Msg("resource payload")
				span.SetAttributes(attribute.String("resource.payload", string(buf)))
			case strings.HasPrefix(cType, MIMEMultipartForm):
			case strings.HasPrefix(cType, MIMEApplicationJSON):
				// b := http.MaxBytesReader(w, b, 1048576)
				body := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(buf)))
				body.DisallowUnknownFields()

				if err := body.Decode(&response); err != nil {
					log.Error().Err(ErrBadRequest(w, r, err)).Msg("Request body contains badly-formed JSON")
					_ = ResponseJSONPayload(w, r, http.StatusBadRequest, nil)
					return
				}
				log.Info().
					Str(HeaderContentType, cType).
					Interface("body", response).Msg("resource payload")
				span.SetAttributes(attribute.String("resource.payload", string(buf)))
			default:
				log.Info().
					Str(HeaderContentType, cType).
					Str("body", string(buf)).Msg("resource payload")
				span.SetAttributes(attribute.String("resource.payload", string(buf)))
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		}

		log.Info().Msgf("tracing form middleware endpoint %s", r.URL.Path)

		traceID := "trace-bifrost-id"
		if traceID := r.Header.Get(HeaderUberTraceId); len(traceID) > 0 {
			traceID = strings.Split(traceID, ":")[0]
		}
		sc := trace.SpanContextFromContext(r.Context())
		if sc.TraceID().IsValid() || sc.SpanID().IsValid() {
			traceID = sc.TraceID().String()
		}
		w.Header().Set(HeaderXTraceId, traceID)

		// adds traceID to a context and get from it latter
		r = r.WithContext(context.WithValue(r.Context(), TracerContext, traceID))

		// pass the span through the request context and serve the request to the next middleware
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		// set the status code
		status := ww.Status()
		span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(status))

		if status >= 500 && status < 600 {
			// mark 5xx server error
			span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(status))
			span.SetAttributes(attribute.Key("event").String("error"))
			span.SetAttributes(attribute.Key("message").String(fmt.Sprintf("%d: %s", status, http.StatusText(status))))
		}
	})
}

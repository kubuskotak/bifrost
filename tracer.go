package bifrost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/opentracing/opentracing-go"
	opExt "github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
)

type SpanContext struct {
	// A probabilistically unique identifier for a [multi-span] trace.
	TraceID uint64

	// A probabilistically unique identifier for a span.
	SpanID uint64

	// Whether the trace is sampled.
	Sampled bool

	// The span's associated baggage.
	Baggage map[string]string // initialized on first use
}
type tracerContext struct {
	Name string
}

func (r *tracerContext) String() string {
	return "context value " + r.Name
}

var TracerContext = tracerContext{Name: "context Respond"}

func HttpTracer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var span opentracing.Span
		carrier := opentracing.HTTPHeadersCarrier(r.Header)
		tracerCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			span = opentracing.StartSpan(r.URL.Path)
		} else {
			span = opentracing.StartSpan(r.URL.Path, opentracing.ChildOf(tracerCtx))
		}
		defer span.Finish()
		defer func() {
			if err := recover(); err != nil {
				opExt.HTTPStatusCode.Set(span, uint16(http.StatusInternalServerError))
				opExt.Error.Set(span, true)
				span.SetTag("error.type", "panic")
				span.LogKV(
					"event", "error",
					"error.kind", "panic",
					"message", err,
					"stack", string(debug.Stack()),
				)
				span.Finish()

				panic(err)
			}
		}()

		span.SetTag("request.id", r.Header.Get("X-Request-Id"))

		opExt.HTTPMethod.Set(span, r.Method)
		opExt.HTTPUrl.Set(span, r.URL.Path)

		resourceName := r.URL.Path
		resourceName = r.Method + " " + resourceName
		span.SetTag("resource.name", resourceName)

		// There's nothing we can do with any errors here.
		if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier); err != nil {
			opExt.Error.Set(span, true)
		} else {
			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), span))
		}

		JSONResponse(w)

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
					_ = ResponsePayload(w, r, http.StatusBadRequest, nil)
					return
				}

				log.Info().
					Str(HeaderContentType, cType).
					Str("body", string(buf)).Msg("resource payload")
				span.SetTag("resource.payload", string(buf))
			case strings.HasPrefix(cType, MIMEMultipartForm):
			case strings.HasPrefix(cType, MIMEApplicationJSON):
				// b := http.MaxBytesReader(w, b, 1048576)
				body := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(buf)))
				body.DisallowUnknownFields()

				if err := body.Decode(&response); err != nil {
					log.Error().Err(ErrBadRequest(w, r, err)).Msg("Request body contains badly-formed JSON")
					_ = ResponsePayload(w, r, http.StatusBadRequest, nil)
					return
				}
				log.Info().
					Str(HeaderContentType, cType).
					Interface("body", response).Msg("resource payload")
				span.SetTag("resource.payload", response)
			default:
				log.Info().
					Str(HeaderContentType, cType).
					Str("body", string(buf)).Msg("resource payload")
				span.SetTag("resource.payload", string(buf))
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		}

		log.Info().Msgf("tracing form middleware endpoint %s", r.URL.Path)

		var traceID string
		if traceID = r.Header.Get(HeaderUberTraceId); len(traceID) > 0 {
			traceID = strings.Split(traceID, ":")[0]
			w.Header().Set(HeaderXTraceId, traceID)
		} else if sc, ok := span.Context().(jaeger.SpanContext); ok {
			traceID = sc.TraceID().String()
			w.Header().Set(HeaderXTraceId, traceID)

		}
		// adds traceID to a context and get from it latter
		r = r.WithContext(context.WithValue(r.Context(), TracerContext, traceID))

		// pass the span through the request context and serve the request to the next middleware
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		// set the status code
		status := ww.Status()
		opExt.HTTPStatusCode.Set(span, uint16(status))

		if status >= 500 && status < 600 {
			// mark 5xx server error
			opExt.Error.Set(span, true)
			span.SetTag("error.type", fmt.Sprintf("%d: %s", status, http.StatusText(status)))
			span.LogKV(
				"event", "error",
				"message", fmt.Sprintf("%d: %s", status, http.StatusText(status)),
			)
		}
	})
}

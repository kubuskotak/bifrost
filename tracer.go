package bifrost

import (
	"bytes"
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
)

func HttpTracer(tracer opentracing.Tracer, operationName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			serverSpan := opentracing.SpanFromContext(ctx)
			if serverSpan == nil {
				// All we can do is create a new root span.
				serverSpan = tracer.StartSpan(operationName)
			} else {
				serverSpan.SetOperationName(operationName)
			}
			defer serverSpan.Finish()

			defer func() {
				if err := recover(); err != nil {
					opExt.HTTPStatusCode.Set(serverSpan, uint16(http.StatusInternalServerError))
					opExt.Error.Set(serverSpan, true)
					serverSpan.SetTag("error.type", "panic")
					serverSpan.LogKV(
						"event", "error",
						"error.kind", "panic",
						"message", err,
						"stack", string(debug.Stack()),
					)
					serverSpan.Finish()

					panic(err)
				}
			}()

			opExt.SpanKindRPCServer.Set(serverSpan)
			opExt.HTTPMethod.Set(serverSpan, r.Method)
			opExt.HTTPUrl.Set(serverSpan, r.URL.String())

			resourceName := r.URL.Path
			resourceName = r.Method + " " + resourceName
			serverSpan.SetTag("resource.name", resourceName)

			// There's nothing we can do with any errors here.
			if err := tracer.Inject(
				serverSpan.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header),
			); err != nil {
				log.Error().Err(err).Msg("Tracing error")
			}

			JSONResponse(w)
			ctx = opentracing.ContextWithSpan(ctx, serverSpan)

			// check content length
			if r.ContentLength > 0 {
				// Request
				var buf []byte
				if r.Body != nil { // Read
					buf, _ = ioutil.ReadAll(r.Body)
				}

				readerBody := ioutil.NopCloser(bytes.NewBuffer(buf))
				mediaBody := ioutil.NopCloser(bytes.NewBuffer(buf))

				bufMediaBody := new(bytes.Buffer)
				_, _ = bufMediaBody.ReadFrom(mediaBody)
				r.Body = readerBody

				// get content-type
				s := strings.ToLower(strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]))

				response := make(map[string]interface{}, 0)

				switch MediaType(s) {
				case TextPlain:
				case FormURLEncoded:
					if err := r.ParseForm(); err != nil {
						log.Error().Err(ErrBadRequest(w, r, err)).Msg("Request body contains badly-formed form-urlencoded")
						_ = ResponsePayload(w, r, http.StatusBadRequest, nil)
						return
					}

					log.Info().
						Str("content-type", s).
						Str("body", bufMediaBody.String()).Msg("request payload")
				case MultipartForm:
				case ApplicationJSON:
					// b := http.MaxBytesReader(w, b, 1048576)
					body := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(bufMediaBody.Bytes())))
					body.DisallowUnknownFields()

					if err := body.Decode(&response); err != nil {
						log.Error().Err(ErrBadRequest(w, r, err)).Msg("Request body contains badly-formed JSON")
						_ = ResponsePayload(w, r, http.StatusBadRequest, nil)
						return
					}
					log.Info().
						Str("content-type", s).
						Interface("body", response).Msg("request payload")
				default:
					log.Info().
						Str("content-type", s).
						Str("body", bufMediaBody.String()).Msg("request payload")
				}
			}

			log.Info().Msgf("tracing form middleware endpoint %s", r.URL.String())
			// pass the span through the request context and serve the request to the next middleware
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r.WithContext(ctx))

			// set the status code
			status := ww.Status()
			opExt.HTTPStatusCode.Set(serverSpan, uint16(status))

			if status >= 500 && status < 600 {
				// mark 5xx server error
				opExt.Error.Set(serverSpan, true)
				serverSpan.SetTag("error.type", fmt.Sprintf("%d: %s", status, http.StatusText(status)))
				serverSpan.LogKV(
					"event", "error",
					"message", fmt.Sprintf("%d: %s", status, http.StatusText(status)),
				)
			}
		})
	}
}

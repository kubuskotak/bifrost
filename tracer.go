package bifrost

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

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

			opExt.SpanKindRPCServer.Set(serverSpan)
			opExt.HTTPMethod.Set(serverSpan, r.Method)
			opExt.HTTPUrl.Set(serverSpan, r.URL.String())

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

				r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
				mediaBody := string(buf)

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
						Str("body", mediaBody).Msg("request payload")
				case MultipartForm:
				case ApplicationJSON:
					// b := http.MaxBytesReader(w, b, 1048576)
					body := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer([]byte(mediaBody))))
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
						Str("body", mediaBody).Msg("request payload")
				}
			}

			log.Info().Msgf("tracing form middleware endpoint %s", r.URL.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

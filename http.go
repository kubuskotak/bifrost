package bifrost

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Welkommen() string {
	return `
========================================================================================
   _     _     _     _     _     _     _
  / \   / \   / \   / \   / \   / \   / \
 ( b ) ( i ) ( f ) ( r ) ( o ) ( s ) ( t )
  \_/   \_/   \_/   \_/   \_/   \_/   \_/
========================================================================================
- port    : %d
-----------------------------------------------------------------------------------------
`
}

type (
	WebPort    int
	WebTimeOut int
	Https      bool
	ServeOpts  struct {
		Port     WebPort
		TimeOut  WebTimeOut
		TLS      Https
		CertFile string
		KeyFile  string
	}
)

type Server struct {
	errChan    chan error
	httpServer *http.Server
	Port       WebPort
	TimeOut    WebTimeOut
	TLS        Https
	CertFile   string
	KeyFile    string
}

func NewServerMux(opts ServeOpts) *Server {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return &Server{httpServer: &http.Server{
		Addr:         fmt.Sprintf(":%d", opts.Port),
		ReadTimeout:  time.Duration(opts.TimeOut) * time.Second,
		WriteTimeout: time.Duration(opts.TimeOut) * time.Second,
	}, Port: opts.Port}
}

func (s *Server) Run(handler http.Handler) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.httpServer.Handler = handler
	// Description µ micro service
	fmt.Println(
		fmt.Sprintf(
			Welkommen(),
			s.Port,
		))
	log.Info().Msgf("Now serving at %s", s.httpServer.Addr)
	go func() {
		if s.TLS {
			log.Info().Msg("Secure with HTTPS")
			s.errChan <- s.httpServer.ListenAndServeTLS(s.CertFile, s.KeyFile)
		} else {
			s.errChan <- s.httpServer.ListenAndServe()
		}
	}()
	s.waitForSignals(ctx)
	s.Stop()
	return nil
}

func (s *Server) Stop() {
	if err := s.httpServer.Close(); err != nil {
		log.Error().Err(err).Msg("Server stopping")
	}
}

func (s *Server) Quiet(ctx context.Context) {
	log.Info().Msg("I have to go...")
	log.Info().Msg("Stopping server gracefully")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Wait is over due to error")
		if err = s.httpServer.Close(); err != nil {
			log.Error().Err(err)
		}
	}
	log.Info().Msgf("Stop server at %s", s.httpServer.Addr)
}

func (s *Server) waitForSignals(ctx context.Context) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-interrupt:
			s.Quiet(ctx)
			log.Error().Err(fmt.Errorf("interrupt received, shutting down")).Msg("Server interrupted through context")
			return
		case err := <-s.errChan:
			s.Quiet(ctx)
			log.Error().Err(err)
			return
		}
	}
}

type MediaType string

const (
	ApplicationJSON MediaType = "application/json"
	FormURLEncoded  MediaType = "application/x-www-form-urlencoded"
	MultipartForm   MediaType = "multipart/form-data"
	TextPlain       MediaType = "text/plain"
)

const (
	// StatusSuccess 2xx
	StatusSuccess  = http.StatusOK
	StatusCreated  = http.StatusCreated
	StatusAccepted = http.StatusAccepted
	// StatusPermanentRedirect 3xx
	StatusPermanentRedirect = http.StatusPermanentRedirect
	// StatusBadRequest 4xx
	StatusBadRequest            = http.StatusBadRequest
	StatusUnauthorized          = http.StatusUnauthorized
	StatusPaymentRequired       = http.StatusPaymentRequired
	StatusForbidden             = http.StatusForbidden
	StatusMethodNotAllowed      = http.StatusMethodNotAllowed
	StatusNotAcceptable         = http.StatusNotAcceptable
	StatusInvalidAuthentication = http.StatusProxyAuthRequired
	StatusRequestTimeout        = http.StatusRequestTimeout
	StatusUnsupportedMediaType  = http.StatusUnsupportedMediaType
	StatusUnProcess             = http.StatusUnprocessableEntity
	// StatusInternalError 5xx
	StatusInternalError           = http.StatusInternalServerError
	StatusBadGatewayError         = http.StatusBadGateway
	StatusServiceUnavailableError = http.StatusServiceUnavailable
	StatusGatewayTimeoutError     = http.StatusGatewayTimeout
)

var statusMap = map[int][]string{
	StatusSuccess:  {"STATUS_OK", "Success"},
	StatusCreated:  {"STATUS_CREATED", "Resource has been created"},
	StatusAccepted: {"STATUS_ACCEPTED", "Resource has been accepted"},

	StatusPermanentRedirect: {"STATUS_PERMANENT_REDIRECT", "The resource has moved to a new location"},

	StatusBadRequest:            {"STATUS_BAD_REQUEST", "Invalid data request"},
	StatusUnauthorized:          {"STATUS_UNAUTHORIZED", "Not authorized to access the service"},
	StatusPaymentRequired:       {"STATUS_PAYMENT_REQUIRED", "Payment need to be done"},
	StatusForbidden:             {"STATUS_FORBIDDEN", "Forbidden access the resource "},
	StatusMethodNotAllowed:      {"STATUS_METHOD_NOT_ALLOWED", "The method specified is not allowed"},
	StatusNotAcceptable:         {"STATUS_NOT_ACCEPTABLE", "Request cannot accepted"},
	StatusInvalidAuthentication: {"STATUS_INVALID_AUTHENTICATION", "The resource owner or authorization server denied the request"},
	StatusRequestTimeout:        {"STATUS_REQUEST_TIMEOUT", "Request Timeout"},
	StatusUnsupportedMediaType:  {"STATUS_UNSUPPORTED_MEDIA_TYPE", "Cannot understand request content"},
	StatusUnProcess:             {"STATUS_UNPROCESSABLE_ENTITY", "Unable to process the contained instructions"},

	StatusInternalError:           {"INTERNAL_SERVER_ERROR", "Oops something went wrong"},
	StatusBadGatewayError:         {"STATUS_BAD_GATEWAY_ERROR", "Oops something went wrong"},
	StatusServiceUnavailableError: {"STATUS_SERVICE_UNAVAILABLE_ERROR", "Service Unavailable"},
	StatusGatewayTimeoutError:     {"STATUS_GATEWAY_TIMEOUT_ERROR", "Gateway Timeout"},
}

func StatusCode(code int) string {
	return statusMap[code][0]
}

func StatusText(code int) string {
	return statusMap[code][1]
}

type Adapter func(w http.ResponseWriter, r *http.Request) error

func HandlerAdapter(a Adapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			status, _ := r.Context().Value(CtxError).(int)
			bytes, err := json.Marshal(&Meta{
				Code:    strconv.Itoa(status),
				Type:    http.StatusText(status),
				Message: err.Error(),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = w.Write(bytes)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
	}
}

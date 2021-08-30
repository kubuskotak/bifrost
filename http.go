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
	// Do not make the application hang when it is shutdown.
	ctxOut, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-interrupt:
			s.Quiet(ctxOut)
			log.Error().Err(fmt.Errorf("interrupt received, shutting down")).Msg("Server interrupted through context")
			return
		case err := <-s.errChan:
			s.Quiet(ctxOut)
			log.Error().Err(err)
			return
		}
	}
}

type Adapter func(w http.ResponseWriter, r *http.Request) error

func HandlerAdapter(a Adapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cType, ok := r.Context().Value(ContentTypeCtxKey).(ContentType); ok {
			w.Header().Set(HeaderContentType, GetIdxContentType(cType))
		}
		if err := a(w, r); err != nil {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			code, _ := r.Context().Value(CtxError).(int)
			null := make(map[string]interface{})
			resp := &Response{
				Version: Version{
					Label:  "v1",
					Number: "0.1.0",
				},
				Meta: Meta{
					Code:    strconv.Itoa(code),
					Type:    http.StatusText(code),
					Message: err.Error(),
				},
				Data:       null,
				Pagination: null,
			}
			if ver, ok := r.Context().Value(CtxVersion).(Version); ok {
				resp.Version = ver
			}
			bytes, err := json.Marshal(resp)
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

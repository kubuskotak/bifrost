package bifrost

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	rpc "google.golang.org/grpc"
)

type (
	GRPCPort     int
	GRPCCallback func(*rpc.Server)
	GRPCOpts     struct {
		Port GRPCPort
		Opts []rpc.ServerOption
	}
)
type GRpc struct {
	errChan   chan error
	rpcServer *rpc.Server
	Port      GRPCPort
	Opts      []rpc.ServerOption
}

func NewServerGRPC(opts GRPCOpts) *GRpc {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return &GRpc{rpcServer: rpc.NewServer(opts.Opts...), Port: opts.Port, Opts: opts.Opts}
}

func (g *GRpc) Run(callback GRPCCallback) error {
	n, err := net.Listen("tcp", fmt.Sprintf(":%v", g.Port))
	if err != nil {
		log.Error().Int("port", int(g.Port)).Err(err).Msg("failed to listen:")
		return err
	}
	// Description µ micro service
	fmt.Println(
		fmt.Sprintf(
			Welkommen(),
			g.Port,
		))
	log.Info().Msgf("Now serving at %v", g.rpcServer.GetServiceInfo())
	callback(g.rpcServer)
	go func() {
		g.errChan <- g.rpcServer.Serve(n)
	}()
	g.waitForSignals()
	g.Stop()
	return nil
}

func (g *GRpc) Stop() {
	log.Info().Msgf("Stop server at :%d", g.Port)
	g.rpcServer.Stop()
}

func (g *GRpc) Quiet() {
	log.Info().Msg("I have to go...")
	log.Info().Msg("Stopping server gracefully")
	g.rpcServer.GracefulStop()
	log.Info().Msgf("Stop server at :%d", g.Port)
}

func (g *GRpc) waitForSignals() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-interrupt:
			g.Quiet()
			log.Error().Err(fmt.Errorf("interrupt received, shutting down")).Msg("Server interrupted through context")
			return
		case err := <-g.errChan:
			g.Quiet()
			log.Error().Err(err)
			return
		}
	}
}

package bifrost

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	rpc "google.golang.org/grpc"
)

func TestNewServerGRPC(t *testing.T) {
	port, err := findOpenPort()
	if err != nil {
		assert.Fail(t, "could not find a testing port")
	}
	t.Log("Using port", port)
	srv := NewServerGRPC(GRPCOpts{
		Port: GRPCPort(port),
	})
	done := make(chan struct{})
	// Make sure server exits when receiving TERM signal.
	go func() {
		time.Sleep(2 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		done <- struct{}{}
	}()

	// Testing
	go func() {
		_ = srv.Run(func(s *rpc.Server) error {
			return nil
		})
	}()
}

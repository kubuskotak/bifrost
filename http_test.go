package bifrost

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var portSync sync.Mutex

func findOpenPort() (int, error) {
	portSync.Lock()
	defer portSync.Unlock()

	min := 10000
	max := 65535
	attempts := 10

	for i := 0; i < attempts; i++ {
		bg := big.NewInt(int64(max - min))
		n, err := rand.Int(rand.Reader, bg)
		if err != nil {
			continue
		}
		port := n.Int64() + int64(min)
		if ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
			// Port unavailable
			continue
		} else if err := ln.Close(); err != nil {
			return 0, err
		}
		return int(port), nil
	}
	return 0, fmt.Errorf("could not find port to use for testing (%d attempts)", attempts)
}

func TestHttpListenAndServe(t *testing.T) {
	port, err := findOpenPort()
	if err != nil {
		assert.Fail(t, "could not find a testing port")
	}
	t.Log("Using port", port)
	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//defer cancel()

	srv := NewServerMux(ServeOpts{
		Port:    WebPort(port),
		TimeOut: WebTimeOut(100),
	})
	done := make(chan struct{})
	// Make sure server exits when receiving TERM signal.
	go func() {
		time.Sleep(2 * time.Second)
		p, err := os.FindProcess(os.Getpid())
		assert.NoError(t, err)
		_ = p.Signal(syscall.SIGTERM)
		done <- struct{}{}
	}()

	// Testing
	go func() {
		_ = srv.Run(http.NewServeMux())
	}()

	resp, queryErr := http.Get(fmt.Sprintf("http://localhost:%d", port))

	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	assert.Nil(t, queryErr)
	respStatus := resp.StatusCode
	assert.Equal(t, http.StatusNotFound, respStatus)
}

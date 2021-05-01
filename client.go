package bifrost

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// Transport constructs a HTTP client with keep-alive turned
// off and a dial-timeout of 30 seconds.
func Transport(tlsInsecure bool, timeout int) http.Client {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
			KeepAlive: 0,
		}).DialContext,

		MaxIdleConns:          1,
		DisableKeepAlives:     true,
		IdleConnTimeout:       time.Duration(timeout*2) * time.Millisecond,
		ExpectContinueTimeout: time.Duration(timeout*10) * time.Millisecond,
	}
	// #nosec
	if tlsInsecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: tlsInsecure}
	}

	proxyClient := http.Client{
		Transport: tr,
	}

	return proxyClient
}

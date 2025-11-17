package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"kus/krzysztof/titler/logging"
	"net"
	"net/http"
	"os"
	"time"
)

const CERT_PATH = "/tmp/kk.crt"

var HttpClient *http.Client

func InitClient() {
	HttpClient = createPooledClient()
}

func createPooledClient() *http.Client {
	var caCertPool *x509.CertPool
	caCert, err := os.ReadFile(CERT_PATH)
	if err == nil {
		caCertPool, err = x509.SystemCertPool()
		if err != nil {
			logging.Log(logging.ERROR, "FATAL error, cannot read root CA certificates")
			os.Exit(1)
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		MaxConnsPerHost:       20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

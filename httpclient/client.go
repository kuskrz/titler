package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"kus/krzysztof/titler/environment"
	"kus/krzysztof/titler/logging"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const CERT_PATH = "/tmp/certs"

var HttpClient *http.Client

func InitClient() {
	HttpClient = createPooledClient()
}

func createPooledClient() *http.Client {
	var caCertPool *x509.CertPool

	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		logging.Log(logging.ERROR, "FATAL error, cannot read root CA certificates")
		os.Exit(1)
	}
	userCerts, err := os.ReadDir(CERT_PATH)
	if err != nil {
		logging.Log(logging.ERROR, "Cannot read directory: "+CERT_PATH)
	}
	for _, cert := range userCerts {
		caCert, err := os.ReadFile(filepath.Join(CERT_PATH, cert.Name()))
		if err == nil {
			if ok := caCertPool.AppendCertsFromPEM(caCert); ok {
				logging.Log(logging.INFO, "Loaded certificate from: "+cert.Name())
			}
		} else {
			logging.Log(logging.DEBUG, "Cannot read: "+cert.Name())
		}
	}

	transport := &http.Transport{
		Proxy: getProxyURL,
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

func getProxyURL(req *http.Request) (*url.URL, error) {
	host := environment.EnvVars["PROXY_HOST"]
	port := environment.EnvVars["PROXY_PORT"]
	if host != "" && port != "" {
		//return url.Parse(req.URL.Scheme + "://" + host + ":" + port)
		return url.Parse("http" + "://" + host + ":" + port)
	}
	return nil, nil
}

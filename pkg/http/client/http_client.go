package client

import (
	"net/http"
)

type HttpClientConfig struct {
}

func NewHttpClient(config *HttpClientConfig) http.Client {
	// Trace an HTTP client by wrapping the transport
	client := http.Client{
		Transport: http.DefaultTransport,
	}

	return client
}

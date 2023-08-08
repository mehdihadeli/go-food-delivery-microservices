package client

import (
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	timeout               = 5 * time.Second
	dialContextTimeout    = 5 * time.Second
	tLSHandshakeTimeout   = 5 * time.Second
	xaxIdleConns          = 20
	maxConnsPerHost       = 40
	retryCount            = 3
	retryWaitTime         = 300 * time.Millisecond
	idleConnTimeout       = 120 * time.Second
	responseHeaderTimeout = 5 * time.Second
)

func NewHttpClient() *resty.Client {
	client := resty.New().
		SetTimeout(timeout).
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryWaitTime)

	return client
}

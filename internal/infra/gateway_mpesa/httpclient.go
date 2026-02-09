package gatewaympesa

import (
	"crypto/tls"
	"net/http"
	"time"
)

func NewMpesaHttpClient(isSandbox bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: isSandbox},
	}

	return &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
}

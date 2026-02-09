package gatewaympesa

import (
	"bytes"
	"io"
	"net/http"
)

type MpesaGateway struct {
	Client *http.Client
}

func (g *MpesaGateway) SendPayment(url, token string, jsonBody []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "developer.mpesa.vm.co.mz")

	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (g *MpesaGateway) SendPost(url, token string, jsonBody []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "developer.mpesa.vm.co.mz")

	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

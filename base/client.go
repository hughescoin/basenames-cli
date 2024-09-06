package base

import (
	"net/http"
)

type Client struct {
	HttpClient http.Client
	RpcURL     string
	PrivateKey string
	Address    string
}

func (c *Client) setHeaders(req *http.Request) {

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

func NewClient(rpcURL, privateKey string, address string) *Client {
	return &Client{
		HttpClient: http.Client{},
		RpcURL:     rpcURL,
		PrivateKey: privateKey,
		Address:    address,
	}
}

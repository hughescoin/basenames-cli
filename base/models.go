package base

type Credentials struct {
	RpcUrl     string
	PrivateKey string
	Address    string
}

var GetBlockResponse struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
}

package sdk

import (
	"github.com/ybbus/jsonrpc/v3"
)

var (
	endpoint      = "https://rpc.mainnet.near.org/"
	defaultClient = jsonrpc.NewClient(endpoint)
)

func WithEndpoint(endpoint string) {
	defaultClient = jsonrpc.NewClient(endpoint)
}

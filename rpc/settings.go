package rpc

import "github.com/RainFallsSilent/activation-statistics/rpc/common"

// ELA RPC config
var ElaRpcConfig = &common.RpcConfig{
	HttpUrl: "http://127.0.0.1:20336",
	User:    "948c9e61637cce3cc318ffc00fb4a11a",
	Pass:    "e1a12c2e46e18e5448b729460bebf6c8",
}

// ESC RPC config
var EscRpcConfig = &common.RpcConfig{
	HttpUrl: "https://api.elastos.io/esc",
	User:    "http://149.248.62.252:20636",
	Pass:    "",
}

// Your Infura API key and the Celo RPC URL
const ApiKey = "76891b8517e248fe9a49473d68f8f7f7"
const RpcURL = "https://mainnet.infura.io/v3/" + ApiKey

// Celo RPC config
var CeloRpcConfig = &common.RpcConfig{
	HttpUrl: RpcURL,
	User:    "",
	Pass:    "",
}

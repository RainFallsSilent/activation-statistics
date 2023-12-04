package rpc

import "github.com/RainFallsSilent/activation-statistics/rpc/common"

// ELA RPC config
var ElaRpcConfig = &common.RpcConfig{
	HttpUrl: "https://api.elastos.io/ela",
	User:    "",
	Pass:    "",
}

// ESC RPC config
var EscRpcConfig = &common.RpcConfig{
	HttpUrl: "https://api.elastos.io/esc",
	User:    "",
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

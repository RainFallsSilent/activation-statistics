package common

type RpcConfig struct {
	HttpUrl string `json:"HttpUrl"`
	User    string `json:"User"`
	Pass    string `json:"Pass"`
}

type Response struct {
	ID      int64       `json:"id"`
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	*Error  `json:"error"`
}

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

package ela

type ELABlockInfo struct {
	Hash              string        `json:"hash"`
	Confirmations     uint32        `json:"confirmations"`
	StrippedSize      uint32        `json:"strippedsize"`
	Size              uint32        `json:"size"`
	Weight            uint32        `json:"weight"`
	Height            uint32        `json:"height"`
	Version           uint32        `json:"version"`
	VersionHex        string        `json:"versionhex"`
	MerkleRoot        string        `json:"merkleroot"`
	Tx                []interface{} `json:"tx"`
	Time              uint32        `json:"time"`
	MedianTime        uint32        `json:"mediantime"`
	Nonce             uint32        `json:"nonce"`
	Bits              uint32        `json:"bits"`
	Difficulty        string        `json:"difficulty"`
	ChainWork         string        `json:"chainwork"`
	PreviousBlockHash string        `json:"previousblockhash"`
	NextBlockHash     string        `json:"nextblockhash"`
	AuxPow            string        `json:"auxpow"`
	MinerInfo         string        `json:"minerinfo"`
}

type TransactionContextInfo struct {
	*TransactionInfo
	BlockHash     string `json:"blockhash"`
	Confirmations uint32 `json:"confirmations"`
	Time          uint32 `json:"time"`
	BlockTime     uint32 `json:"blocktime"`
}

type AttributeInfo struct {
	Usage byte   `json:"usage"`
	Data  string `json:"data"`
}

type InputInfo struct {
	TxID     string `json:"txid"`
	VOut     uint16 `json:"vout"`
	Sequence uint32 `json:"sequence"`
}

type OutputPayloadInfo interface{}
type RpcOutputInfo struct {
	Value         string            `json:"value"`
	Index         uint32            `json:"n"`
	Address       string            `json:"address"`
	AssetID       string            `json:"assetid"`
	OutputLock    uint32            `json:"outputlock"`
	OutputType    uint32            `json:"type"`
	OutputPayload OutputPayloadInfo `json:"payload"`
}

type ProgramInfo struct {
	Code      string `json:"code"`
	Parameter string `json:"parameter"`
}

type TransactionInfo struct {
	TxID           string          `json:"txid"`
	Hash           string          `json:"hash"`
	Size           uint32          `json:"size"`
	VSize          uint32          `json:"vsize"`
	Version        byte            `json:"version"`
	TxType         byte            `json:"type"`
	PayloadVersion byte            `json:"payloadversion"`
	Payload        interface{}     `json:"payload"`
	Attributes     []AttributeInfo `json:"attributes"`
	Inputs         []InputInfo     `json:"vin"`
	Outputs        []RpcOutputInfo `json:"vout"`
	LockTime       uint32          `json:"locktime"`
	Programs       []ProgramInfo   `json:"programs"`
}

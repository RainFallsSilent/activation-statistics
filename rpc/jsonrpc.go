package rpc

import (
	"fmt"

	"github.com/RainFallsSilent/activation-statistics/rpc/ela"
	"github.com/RainFallsSilent/activation-statistics/rpc/esc"
)

// ELA RPC
func ELAGetCurrentBlockHeight() (int, error) {
	resp, err := ela.CallAndUnmarshal("getblockcount", nil, ElaRpcConfig)
	if err != nil {
		return 0, err
	}
	var res int
	if err = ela.Unmarshal(&resp, &res); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 0, err
	}

	return res, nil
}

func ELAGetBlockbyheight(height string) (*ela.ELABlockInfo, error) {
	resp, err := ela.CallAndUnmarshal("getblockbyheight", ela.Param("height", height), ElaRpcConfig)
	if err != nil {
		return nil, err
	}
	var res ela.ELABlockInfo
	if err = ela.Unmarshal(&resp, &res); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	return &res, nil
}

func ELAGetRawTransaction(txid string) (*ela.TransactionContextInfo, error) {
	resp, err := ela.CallAndUnmarshal("getrawtransaction", ela.Param("txid", txid).Add("verbose", true), ElaRpcConfig)
	if err != nil {
		return nil, err
	}
	var res ela.TransactionContextInfo
	if err = ela.Unmarshal(&resp, &res); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	return &res, nil
}

// ESC RPC
func ESCGetTransactionByHash(hash string) (*esc.TransactionResult, error) {
	resp, err := esc.CallAndUnmarshal("eth_getTransactionByHash", esc.ParamList(hash), EscRpcConfig)
	if err != nil {
		return nil, err
	}
	var res esc.TransactionResult
	if err = esc.Unmarshal(&resp, &res); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	return &res, nil
}

func ESCGetBlockByNumber(number string) (interface{}, error) {
	resp, err := esc.CallAndUnmarshal("eth_getBlockByNumber", esc.ParamList(number, true), EscRpcConfig)
	if err != nil {
		return nil, err
	}

	// todo change to struct and return according to GetTransactionByHash

	return resp, nil
}

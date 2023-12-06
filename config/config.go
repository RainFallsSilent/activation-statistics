package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/RainFallsSilent/activation-statistics/rpc/common"
)

type ConfigFile struct {
	Days         uint32           `json:"Days"`
	StartHour    uint32           `json:"StartHour"`
	ELARpcConfig common.RpcConfig `json:"ELARpcConfig"`
	ESCRpcConfig common.RpcConfig `json:"ESCRpcConfig"`
}

func InitConfig(path string) *ConfigFile {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Printf("Read config file error: %v\n", e)
		return nil
	}
	i := ConfigFile{
		Days:      2,
		StartHour: 8,
		ELARpcConfig: common.RpcConfig{
			HttpUrl: "https://api.elastos.io/ela",
			User:    "",
			Pass:    "",
		},
		ESCRpcConfig: common.RpcConfig{
			HttpUrl: "https://api.elastos.io/esc",
			User:    "",
			Pass:    "",
		},
	}

	// Remove the UTF-8 Byte Order Mark
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

	e = json.Unmarshal(file, &i)
	if e != nil {
		fmt.Printf("Unmarshal config file error: %v\n", e)
		return nil
	}

	return &i
}

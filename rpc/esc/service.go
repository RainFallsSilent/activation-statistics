package esc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/RainFallsSilent/activation-statistics/rpc/common"
)

func CallAndUnmarshal(method string, params Parameter, config *common.RpcConfig) (interface{}, error) {
	body, err := Call(method, params, config)
	if err != nil {
		return nil, err
	}

	resp := common.Response{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return string(body), nil
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	return resp.Result, nil
}

func Call(method string, params Parameter, config *common.RpcConfig) ([]byte, error) {
	url := config.HttpUrl
	var parm string
	parm = "["
	for _, p := range params {
		parm += "\"" + p + "\""

	}
	parm += "]"
	payload := []byte(`{
		"jsonrpc": "2.0",
		"method": "` + method + `",
		"params": ` + parm + `,
		"id": 1
	}`)

	fmt.Println("call:", string(payload))
	resp, err := post(url, "application/json", strings.NewReader(string(payload)))
	if err != nil {
		fmt.Println("POST requset err:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Unmarshal(result interface{}, target interface{}) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, target); err != nil {
		return err
	}
	return nil
}

func post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	client := *http.DefaultClient
	client.Timeout = time.Minute
	return client.Do(req)
}

package ela

import (
	"encoding/base64"
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

func CallAndUnmarshal(method string, params map[string]interface{}, config *common.RpcConfig) (interface{}, error) {
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

func Call(method string, params map[string]interface{}, config *common.RpcConfig) ([]byte, error) {
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"params": params,
	})
	if err != nil {
		return nil, err
	}

	resp, err := post(config.HttpUrl, "application/json", config.User, config.Pass, strings.NewReader(string(data)))
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

func post(url string, contentType string, user string, pass string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	auth := user + ":" + pass
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", basicAuth)
	req.Header.Set("Content-Type", contentType)

	client := *http.DefaultClient
	client.Timeout = time.Minute
	return client.Do(req)
}

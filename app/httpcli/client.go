package httpcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/osspkg/go-sdk/ioutil"

	"github.com/osspkg/goppy/plugins"
)

var Plugin = plugins.Plugin{
	Inject: NewClient,
}

type (
	Client struct {
		cli *http.Client
	}
)

func NewClient() *Client {
	return &Client{cli: &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
		Timeout: 5 * time.Second,
	}}
}

func (v *Client) Get(uri string) (int, []byte, error) {
	return v.Call(http.MethodGet, uri, nil)
}

func (v *Client) GetJSON(uri string, out interface{}) error {
	code, b, err := v.Call(http.MethodGet, uri, nil)
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return fmt.Errorf("response code: %d", code)
	}
	return json.Unmarshal(b, out)
}

func (v *Client) Call(method, uri string, body []byte) (int, []byte, error) {
	req, err := http.NewRequest(method, uri, bytes.NewReader(body))
	if err != nil {
		return 0, nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Connection", "keep-alive")
	resp, err := v.cli.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("make request: %w", err)
	}
	b, err0 := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, b, err0
}

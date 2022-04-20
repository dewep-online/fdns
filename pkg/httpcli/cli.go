package httpcli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/deweppro/go-errors"
)

type Client struct {
	cli *http.Client
}

func New() *Client {
	cli := &http.Client{
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
	}
	return &Client{cli: cli}
}

func (v *Client) Call(method, uri string, body []byte) (int, []byte, error) {
	req, err := http.NewRequest(method, uri, bytes.NewReader(body))
	if err != nil {
		return 0, nil, errors.WrapMessage(err, "create request")
	}
	req.Header.Set("Connection", "keep-alive")
	resp, err := v.cli.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("make request: %w", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, b, errors.Wrap(err, resp.Body.Close())
}

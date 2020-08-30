// Copyright 2020 PingCAP, Inc. Licensed under Apache-2.0.

package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Client is a REST client interacting with the Grafana API.
//
// This class is limited in scope that only deals with APIs useful for
// `grafsnap`.
type Client struct {
	http   *http.Client
	base   *url.URL
	key    string
	config *Config
}

// Config configures the Grafana client.
type Config struct {
	Key         string
	Base        string
	Concurrency int
}

// NewClient constructs a new Grafana client.
func NewClient(config *Config, client *http.Client) (*Client, error) {
	baseURL, err := url.Parse(config.Base)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL '%s':\n%w", config.Base, err)
	}

	key := config.Key
	if colonPos := strings.IndexByte(key, ':'); colonPos >= 0 {
		username := key[:colonPos]
		password := key[colonPos+1:]
		baseURL.User = url.UserPassword(username, password)
		key = ""
	}

	return &Client{
		http:   client,
		base:   baseURL,
		key:    key,
		config: config,
	}, nil
}

// do performs a JSON request to the Grafana server. This method will fill in
// the authentication details.
//
// `resultJSON` should be a pointer to an object to receive the API output on
// success.
func (cli *Client) do(ctx context.Context, method string, api string, query url.Values, body []byte, resultJSON interface{}) error {
	u := *cli.base
	u.Path = path.Join(u.Path, api)
	u.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot construct request to '%s':\n%w", api, err)
	}
	if len(cli.key) > 0 {
		req.Header.Add("Authorization", "Bearer "+cli.key)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := cli.http.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send request to '%s':\n%w", api, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received failure from '%s' (HTTP %d):\n%s", api, resp.StatusCode, msg)
	}

	if err := json.NewDecoder(resp.Body).Decode(resultJSON); err != nil {
		return fmt.Errorf("received invalid JSON from '%s':\n%w", api, err)
	}
	return nil
}

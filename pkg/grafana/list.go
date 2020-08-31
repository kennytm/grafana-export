// Copyright 2020 PingCAP, Inc. Licensed under Apache-2.0.

package grafana

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// ListDashboardResult is the result of listing dashboards from Grafana.
type ListDashboardResult struct {
	// Title is the human-readable title of the dashboard.
	Title string `json:"title"`
	// UID is the unique identifier of the dashboard.
	UID string `json:"uid"`
	// URL is the URL to the dashboard.
	URL string `json:"url"`
}

// ListDashboards lists all dashboards.
func (cli *Client) ListDashboards(ctx context.Context) ([]ListDashboardResult, error) {
	// TODO: evaluate the templating.

	params := make(url.Values, 1)
	params.Add("type", "dash-db")

	var result []ListDashboardResult
	if err := cli.do(ctx, "GET", "api/search", params, nil, &result); err != nil {
		return nil, fmt.Errorf("cannot list dashboards:\n%w", err)
	}

	// DO NOT use `cli.base` here to avoid exposing the username & password.
	base := strings.TrimRight(cli.config.Base, "/")
	for i, r := range result {
		result[i].URL = base + r.URL
	}

	return result, nil
}

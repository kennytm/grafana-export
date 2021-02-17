// Copyright 2020 PingCAP, Inc. Licensed under Apache-2.0.

package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/kennytm/grafana-export/pkg/archive"
	"github.com/kennytm/grafana-export/pkg/json2"
)

type AnnotationModel struct {
	BuiltIn     int32    `json:"builtIn,omitempty"`
	Datasource  string   `json:"datasource,omitempty"`
	Expr        string   `json:"expr,omitempty"`
	TextFormat  string   `json:"textFormat,omitempty"`
	TitleFormat string   `json:"titleFormat,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Step        string   `json:"step,omitempty"`

	SnapshotData []interface{} `json:"snapshotData"`
}

type DataPoint struct {
	Timestamp uint64
	Value     float64
}

func (dp DataPoint) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%g,%d]", dp.Value, dp.Timestamp)), nil
}

type PanelModel struct {
	Datasource string `json:"datasource,omitempty"`

	Targets []struct {
		Expr string `json:"expr,omitempty"`
	} `json:"targets"`

	SnapshotData []struct {
		DataPoints []DataPoint `json:"datapoints"`
		Query      string      `json:"query"`
		Target     string      `json:"target"`
	} `json:"snapshotData"`
}

type DashboardModel struct {
	ID int64 `json:"id"`

	Snapshot struct {
		Timestamp time.Time `json:"timestamp"`
	} `json:"snapshot"`

	Annotations struct {
		List []json2.Map `json:"list"`
	} `json:"annotations"`

	Panels []json2.Map `json:"panels"`
}

type SnapshotModel struct {
	Meta struct {
		IsSnapshot bool      `json:"isSnapshot"`
		Type       string    `json:"type"`
		Expires    time.Time `json:"expires"`
		Created    time.Time `json:"created"`
	} `json:"meta"`
	Dashboard json2.Map `json:"dashboard"`
}

func (client *Client) ExportDashboards(
	ctx context.Context,
	uids []string,
	from Timestamp,
	to Timestamp,
	variables map[string]string,
	archive archive.Archive,
) error {
	dataSources, err := client.ListDataSources(ctx)
	if err != nil {
		return err
	}

	// TODO respect concurrency
	for _, uid := range uids {
		snapshot, err := client.ExportDashboard(ctx, uid, from, to, variables, dataSources)
		if err != nil {
			return err
		}

		if err := archive.WriteToFile(uid+".json", func(w io.Writer) error {
			return json.NewEncoder(w).Encode(snapshot)
		}); err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) ExportDashboard(
	ctx context.Context,
	uid string,
	from Timestamp,
	to Timestamp,
	variables map[string]string,
	dataSources map[string]DataSource,
) (snapshot SnapshotModel, err error) {
	if err = client.do(ctx, "GET", "api/dashboards/uid/"+uid, nil, nil, &snapshot); err != nil {
		err = fmt.Errorf("cannot get dashboard model for UID '%s':\n%w", uid, err)
		return
	}

	var dashboard DashboardModel
	if err = snapshot.Dashboard.Decode(&dashboard); err != nil {
		err = fmt.Errorf("invalid dashboard model for UID '%s':\n%w", uid, err)
		return
	}

	// fill in annotations

	return
}

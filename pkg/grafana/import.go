// Copyright 2020 PingCAP, Inc. Licensed under Apache-2.0.

package grafana

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ImportSnapshotResult is the result of importing a file into Grafana.
type ImportSnapshotResult struct {
	File    string                 `json:"file"`
	URL     string                 `json:"url,omitempty"`
	Content []ImportSnapshotResult `json:"content,omitempty"`
}

// ImportSnapshot imports the snapshot back into Grafana.
// Returns the URL to access the snapshot.
func (cli *Client) ImportSnapshot(ctx context.Context, name string, body io.Reader) (string, error) {
	var request struct {
		Dashboard json.RawMessage `json:"dashboard"`
		Name      string          `json:"name"`
	}

	var response struct {
		URL string
	}

	if err := json.NewDecoder(body).Decode(&request); err != nil {
		return "", fmt.Errorf("invalid JSON file:\n%w", err)
	}
	request.Name = name
	data, _ := json.Marshal(request)

	if err := cli.do(ctx, "POST", "api/snapshots", nil, data, &response); err != nil {
		return "", fmt.Errorf("cannot upload snapshot:\n%w", err)
	}

	return response.URL, nil
}

// ImportSnapshotsFromFiles imports multiple snapshot files in parallel into
// Grafana. The files can be plain JSON files backed up from the local snapshot,
// or a zip file containing these JSON snapshots.
func (cli *Client) ImportSnapshotsFromFiles(ctx context.Context, files []string) ([]ImportSnapshotResult, error) {
	// TODO respect concurrency

	result := make([]ImportSnapshotResult, 0, len(files))
	for _, fileName := range files {
		res := ImportSnapshotResult{File: filepath.Base(fileName)}

		f, err := os.Open(fileName)
		if err != nil {
			return nil, fmt.Errorf("cannot open snapshot file '%s':\n%w", fileName, err)
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("cannot get file size of '%s':\n%w", fileName, err)
		}

		zipReader, err := zip.NewReader(f, stat.Size())
		if err != nil {
			// should be just raw JSON content. re-read from the start.

			_, err = f.Seek(0, os.SEEK_SET)
			if err != nil {
				return nil, fmt.Errorf("cannot seek '%s' from start:\n%w", fileName, err)
			}

			res.URL, err = cli.ImportSnapshot(ctx, res.File, f)
			if err != nil {
				return nil, fmt.Errorf("cannot import snapshot file '%s':\n%w", fileName, err)
			}

		} else {
			// if the file is a zip file, decompress and read all files in it as JSON.
			res.Content = make([]ImportSnapshotResult, 0, len(zipReader.File))
			for _, subFile := range zipReader.File {
				sf, err := subFile.Open()
				if err != nil {
					return nil, fmt.Errorf("cannot open snapshot file '%s' in zip archive '%s':\n%w", subFile.Name, fileName, err)
				}
				defer sf.Close()

				url, err := cli.ImportSnapshot(ctx, res.File+"//"+subFile.Name, sf)
				if err != nil {
					return nil, fmt.Errorf("cannot import snapshot file '%s' in zip archive '%s':\n%w", subFile.Name, fileName, err)
				}
				res.Content = append(res.Content, ImportSnapshotResult{File: subFile.Name, URL: url})
			}
		}

		result = append(result, res)
	}

	return result, nil
}

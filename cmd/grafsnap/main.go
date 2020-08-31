// Copyright 2020 PingCAP, Inc. Licensed under Apache-2.0.

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kennytm/grafana-export/pkg/grafana"

	"github.com/spf13/cobra"
)

var (
	client *grafana.Client
)

func printPrettyJSON(obj interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(obj)
}

func listCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list dashboards",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := client.ListDashboards(cmd.Context())
			if err != nil {
				return err
			}
			return printPrettyJSON(result)
		},
	}
	return cmd
}

func importCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "import",
		Aliases: []string{"im", "upload", "up"},
		Short:   "import snapshot JSON files back into Grafana",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := client.ImportSnapshotsFromFiles(cmd.Context(), cmd.Flags().Args())
			if err != nil {
				return err
			}
			return printPrettyJSON(result)
		},
	}
	return cmd
}

func rootCommand() *cobra.Command {
	var config grafana.Config

	cmd := &cobra.Command{
		Use:              "grafsnap",
		Short:            "Grafana snapshot taker",
		TraverseChildren: true,
		SilenceUsage:     true,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			if len(config.Key) == 0 {
				config.Key = os.Getenv("GRAFANA_API_KEY")
			}
			if len(config.Key) == 0 {
				config.Key = "admin:admin"
			}

			var err error
			client, err = grafana.NewClient(&config, http.DefaultClient)
			if err != nil {
				return err
			}
			return nil
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringVarP(&config.Key, "api-key", "k", "", "API key or 'username:password'")
	flags.StringVarP(&config.Base, "base", "b", "http://localhost:3000/", "base URL of Grafana instance")
	flags.IntVarP(&config.Concurrency, "concurrency", "j", 4, "maximum number of parallel requests")
	cmd.AddCommand(
		listCommand(),
		importCommand(),
	)
	return cmd
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancelCh := make(chan os.Signal, 1)
	signal.Notify(cancelCh, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		sig := <-cancelCh
		log.Printf("received signal %s, terminating", sig)
		cancel()
	}()

	if rootCommand().ExecuteContext(ctx) != nil {
		os.Exit(1)
	}
}

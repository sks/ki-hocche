/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/sks/kihocche/pkg/formatter"
	"github.com/sks/kihocche/pkg/logger"
	"github.com/sks/kihocche/pkg/models"
	"github.com/sks/kihocche/pkg/scmscanner"
	"github.com/spf13/cobra"
)

var journeOpts = struct {
	config scmscanner.Config
	output formatter.Config
	filter scmscanner.Filter
}{
	config: scmscanner.Config{
		Type: "github",
		Config: scmscanner.SCMConfig{
			Api: "https://api.github.com",
		},
	},
	filter: scmscanner.Filter{
		Since: 7 * 24 * time.Hour,
		Events: []string{
			models.EventTypePush,
			models.EventTypePullRequest,
			models.EventTypeRelease,
		},
		Repos: []string{},
	},
	output: formatter.Config{
		Type: "json",
	},
}

// journeyCmd represents the journey command
var journeyCmd = &cobra.Command{
	Use:   "journey",
	Short: "shows the journey of the git repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logger.GetLogger(cmd.Context()).With("command", "journey")
		defer func(startTime time.Time) {
			logger.Debug("command completed", "duration", time.Since(startTime).String())
		}(time.Now())
		scanner, err := journeOpts.config.Scanner()
		if err != nil {
			return fmt.Errorf("failed to get scanner: %w", err)
		}
		events, err := scanner.Get(cmd.Context(), journeOpts.filter)
		if err != nil {
			return fmt.Errorf("failed to get events: %w", err)
		}
		return journeOpts.output.Write(cmd.Context(), events)
	},
}

func init() {
	rootCmd.AddCommand(journeyCmd)

	journeyCmd.Flags().StringVar(&journeOpts.config.Type, "provider", journeOpts.config.Type, "provider to use")
	journeyCmd.Flags().StringVar(&journeOpts.config.Config.Token, "token", os.Getenv("GH_TOKEN"), "token to use")
	journeyCmd.Flags().StringVar(&journeOpts.config.Config.Api, "api", journeOpts.config.Config.Api, "api to use")
	journeyCmd.Flags().StringSliceVar(&journeOpts.filter.Repos, "repos", journeOpts.filter.Repos, "repos to filter")
	journeyCmd.Flags().StringSliceVar(&journeOpts.filter.Namespace, "namespace", journeOpts.filter.Namespace, "namespace to filter")
	journeyCmd.Flags().StringVarP(&journeOpts.output.Type, "type", "t", journeOpts.output.Type, "output type")
	journeyCmd.Flags().StringVarP(&journeOpts.output.Output, "output", "o", journeOpts.output.Output, "output file")
}

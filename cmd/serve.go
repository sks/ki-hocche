/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"time"

	"gocloud.dev/secrets"
	_ "gocloud.dev/secrets/localsecrets"

	"github.com/sks/kihocche/pkg/osutils"
	"github.com/sks/kihocche/pkg/server"
	"github.com/spf13/cobra"
)

var serverOpts = struct {
	config    server.Config
	secretKey string
}{
	config: server.Config{
		Port:                    "8080",
		ShutdownTimeoutDuration: 5 * time.Second,
	},
	secretKey: osutils.Getenv("ENCRYPTION_KEY", "smGbjm71Nxd1Ig5FS0wj9SlbzAIrnolCz9bQQ6uAhl4="),
}

// serverCmd represents the server command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start webserver that can serve ics files",
	RunE: func(cmd *cobra.Command, args []string) error {
		serveMux := http.NewServeMux()
		randomKeyKeeper, err := secrets.OpenKeeper(cmd.Context(), fmt.Sprintf("base64key://%s", serverOpts.secretKey))
		if err != nil {
			return fmt.Errorf("failed to open keeper: %w", err)
		}
		defer randomKeyKeeper.Close()

		router := server.NewRouter(randomKeyKeeper)
		serveMux.Handle("GET /generate_link", http.HandlerFunc(router.GenerateLinkPage))
		serveMux.Handle("POST /generate_link", http.HandlerFunc(router.GenerateLink))
		serveMux.Handle("GET /subscribe.ics", http.HandlerFunc(router.ServeICal))

		return serverOpts.config.Start(cmd.Context(), serveMux)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&serverOpts.config.Port, "port", "p", "8080", "Port to start the server")
}

package cli

import (
	"os"
	"time"

	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/config"
	"github.com/makeplane/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	flagToken      string
	flagAPIURL     string
	flagWorkspace  string
	flagProject    string
	flagMaxRetries int
	flagTimeout    time.Duration
	flagNoColor    bool
	exitCode       = 0
)

var rootCmd = &cobra.Command{
	Use:   "plane",
	Short: "Plane CLI - interact with Plane API from terminal",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func ExitCode() int { return exitCode }

func init() {
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "API token (overrides PLANE_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&flagAPIURL, "api-url", "", "Base API URL (overrides PLANE_API_URL)")
	rootCmd.PersistentFlags().StringVarP(&flagWorkspace, "workspace", "w", "", "Workspace slug")
	rootCmd.PersistentFlags().StringVarP(&flagProject, "project", "p", "", "Project UUID")
	rootCmd.PersistentFlags().IntVar(&flagMaxRetries, "max-retries", 4, "Max retry attempts")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 30*time.Second, "Total timeout for all retries")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable colored output")
}

func resolveContext() (*api.Client, error) {
	token, err := config.ResolveToken(flagToken, getEnv("PLANE_TOKEN"))
	if err != nil {
		return nil, err
	}
	ws, err := config.ResolveWorkspace(flagWorkspace, getEnv("PLANE_WORKSPACE"))
	if err != nil {
		return nil, err
	}
	prj, err := config.ResolveProject(flagProject, getEnv("PLANE_PROJECT"))
	if err != nil {
		return nil, err
	}
	apiURL := config.ResolveAPIURL(flagAPIURL, getEnv("PLANE_API_URL"))

	cfg := api.DefaultRetryConfig()
	if flagMaxRetries != 4 {
		cfg.MaxRetries = flagMaxRetries
	}
	if flagTimeout != 30*time.Second {
		cfg.TotalTimeout = flagTimeout
	}
	return api.NewClient(token, apiURL, ws, prj, cfg), nil
}

func getEnv(key string) string { return os.Getenv(key) }

func handleAPIError(code int, err error, body []byte) {
	switch {
	case err == api.ErrCircuitBreakerOpen:
		exitCode = 2
		output.WriteError("circuit_breaker_open", "Circuit breaker open after consecutive failures")
	case code == 401 || code == 403:
		exitCode = 1
		output.WriteError("unauthorized", "Invalid or missing API token")
	case code == 404:
		exitCode = 1
		output.WriteError("not_found", "Resource not found")
	case code >= 500:
		exitCode = 1
		output.WriteError("server_error", "Plane API server error")
	default:
		exitCode = 4
		output.WriteError("network", err.Error())
	}
}

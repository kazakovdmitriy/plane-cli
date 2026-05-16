package cli

import (
	"github.com/kazakovdmitriy/plane-cli/internal/config"
	"github.com/kazakovdmitriy/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage Plane context (workspace, project, token)",
}

var contextSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set workspace, project and API token",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, _ := config.Load()
		if ctx == nil {
			ctx = &config.Context{}
		}
		if ws, _ := cmd.Flags().GetString("workspace"); ws != "" {
			ctx.Workspace = ws
		}
		if prj, _ := cmd.Flags().GetString("project"); prj != "" {
			ctx.Project = prj
		}
		if tok, _ := cmd.Flags().GetString("token"); tok != "" {
			ctx.Token = tok
		}
		if url, _ := cmd.Flags().GetString("api-url"); url != "" {
			ctx.APIURL = url
		}
		if err := config.Save(ctx); err != nil {
			return err
		}
		output.WriteItem(map[string]string{"status": "saved"})
		return nil
	},
}

var contextShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current context",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := config.Load()
		if err != nil {
			return err
		}
		display := map[string]string{
			"workspace": ctx.Workspace,
			"project":   ctx.Project,
			"api_url":   ctx.APIURL,
		}
		output.WriteItem(display)
		return nil
	},
}

var contextUnsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Remove saved context",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Delete(); err != nil {
			return err
		}
		output.WriteItem(map[string]string{"status": "removed"})
		return nil
	},
}

func init() {
	contextSetCmd.Flags().String("workspace", "", "Workspace slug")
	contextSetCmd.Flags().String("project", "", "Project UUID")
	contextSetCmd.Flags().String("token", "", "API token")
	contextSetCmd.Flags().String("api-url", "", "Base API URL")
	contextCmd.AddCommand(contextSetCmd, contextShowCmd, contextUnsetCmd)
	rootCmd.AddCommand(contextCmd)
}

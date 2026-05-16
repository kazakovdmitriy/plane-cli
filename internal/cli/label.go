package cli

import (
	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var labelCmd = &cobra.Command{Use: "label", Short: "Manage labels"}

var labelListCmd = &cobra.Command{
	Use: "list", Short: "List labels",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointLabels(), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleCollectionResponse(body, code)
		return nil
	},
}

var labelCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a label",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			output.WriteError("validation", "--name is required")
			return nil
		}
		color, _ := cmd.Flags().GetString("color")
		if color == "" {
			color = "#6b7280"
		}
		payload := map[string]string{"name": name, "color": color}
		if v, _ := cmd.Flags().GetString("parent"); v != "" {
			payload["parent"] = v
		}
		body, code, err := client.Do("POST", api.EndpointLabels(), payload, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

func init() {
	labelCreateCmd.Flags().String("name", "", "Label name (required)")
	labelCreateCmd.Flags().String("color", "", "Color in hex (default: #6b7280)")
	labelCreateCmd.Flags().String("parent", "", "Parent label ID")
	labelCmd.AddCommand(labelListCmd, labelCreateCmd)
	rootCmd.AddCommand(labelCmd)
}

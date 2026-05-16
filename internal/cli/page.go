package cli

import (
	"github.com/kazakovdmitriy/plane-cli/internal/api"
	"github.com/kazakovdmitriy/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var pageCmd = &cobra.Command{Use: "page", Short: "Manage pages"}

var pageListCmd = &cobra.Command{
	Use: "list", Short: "List pages",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointPages(), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleCollectionResponse(body, code)
		return nil
	},
}

var pageGetCmd = &cobra.Command{
	Use: "get <id>", Short: "Get a page", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointPage(args[0]), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

var pageCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a page",
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
		payload := map[string]string{"name": name}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			payload["description_html"] = v
		}
		if v, _ := cmd.Flags().GetString("content"); v != "" {
			payload["description_html"] = v
		}
		body, code, err := client.Do("POST", api.EndpointPages(), payload, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

func init() {
	pageCreateCmd.Flags().String("name", "", "Page name (required)")
	pageCreateCmd.Flags().String("description", "", "Description (HTML)")
	pageCreateCmd.Flags().String("content", "", "Page content (HTML)")
	pageCmd.AddCommand(pageListCmd, pageGetCmd, pageCreateCmd)
	rootCmd.AddCommand(pageCmd)
}

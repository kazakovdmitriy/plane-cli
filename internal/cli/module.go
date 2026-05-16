package cli

import (
	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{Use: "module", Short: "Manage modules"}

var moduleListCmd = &cobra.Command{
	Use: "list", Short: "List modules",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointModules(), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleCollectionResponse(body, code)
		return nil
	},
}

var moduleGetCmd = &cobra.Command{
	Use: "get <id>", Short: "Get a module", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointModule(args[0]), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

func init() {
	moduleCmd.AddCommand(moduleListCmd, moduleGetCmd)
	rootCmd.AddCommand(moduleCmd)
}

package cli

import (
	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var stateCmd = &cobra.Command{Use: "state", Short: "Manage states"}

var stateListCmd = &cobra.Command{
	Use: "list", Short: "List states",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointStates(), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleCollectionResponse(body, code)
		return nil
	},
}

var stateGetCmd = &cobra.Command{
	Use: "get <id>", Short: "Get a state", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointState(args[0]), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

func init() {
	stateCmd.AddCommand(stateListCmd, stateGetCmd)
	rootCmd.AddCommand(stateCmd)
}

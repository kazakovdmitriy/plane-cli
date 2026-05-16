package cli

import (
	"github.com/kazakovdmitriy/plane-cli/internal/api"
	"github.com/kazakovdmitriy/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var memberCmd = &cobra.Command{Use: "member", Short: "Manage members"}

var memberListCmd = &cobra.Command{
	Use: "list", Short: "List members",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointMembers(), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleCollectionResponse(body, code)
		return nil
	},
}

func init() {
	memberCmd.AddCommand(memberListCmd)
	rootCmd.AddCommand(memberCmd)
}

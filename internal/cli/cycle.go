package cli

import (
	"encoding/json"
	"net/url"

	"github.com/kazakovdmitriy/plane-cli/internal/api"
	"github.com/kazakovdmitriy/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var cycleCmd = &cobra.Command{
	Use:   "cycle",
	Short: "Manage cycles",
}

var cycleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cycles",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		query := url.Values{}
		if v, _ := cmd.Flags().GetString("sort-by"); v != "" {
			query.Set("order_by", v)
		}
		body, code, err := client.Do("GET", api.EndpointCycles(), nil, query)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleCollectionResponse(body, code)
		return nil
	},
}

var cycleGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a cycle",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointCycle(args[0]), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

var cycleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a cycle",
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
		payload := map[string]interface{}{"name": name}
		if v, _ := cmd.Flags().GetString("project"); v != "" {
			payload["project"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			payload["description"] = v
		}
		if v, _ := cmd.Flags().GetString("start-date"); v != "" {
			payload["start_date"] = v
		}
		if v, _ := cmd.Flags().GetString("end-date"); v != "" {
			payload["end_date"] = v
		}
		body, code, err := client.Do("POST", api.EndpointCycles(), payload, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

func init() {
	cycleListCmd.Flags().String("sort-by", "", "Sort by: name, start_date, -start_date, end_date, -end_date")
	cycleCreateCmd.Flags().String("name", "", "Cycle name (required)")
	cycleCreateCmd.Flags().String("project", "", "Project UUID")
	cycleCreateCmd.Flags().String("description", "", "Description")
	cycleCreateCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	cycleCreateCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")

	cycleCmd.AddCommand(cycleListCmd, cycleGetCmd, cycleCreateCmd)
	rootCmd.AddCommand(cycleCmd)
}

func handleCollectionResponse(body []byte, code int) {
	var raw interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		output.WriteError("parse_error", "Failed to parse API response")
		return
	}
	switch v := raw.(type) {
	case []interface{}:
		items := make([]interface{}, len(v))
		for i, item := range v {
			items[i] = item
		}
		output.WriteItems(items)
	case map[string]interface{}:
		if results, ok := v["results"]; ok {
			if arr, ok := results.([]interface{}); ok {
				items := make([]interface{}, len(arr))
				copy(items, arr)
				output.WriteItems(items)
				return
			}
		}
		output.WriteItem(raw)
	default:
		output.WriteItem(raw)
	}
}

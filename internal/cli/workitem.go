package cli

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var wiCmd = &cobra.Command{
	Use:     "work-item",
	Aliases: []string{"wi", "work-items"},
	Short:   "Manage work items",
}

var wiListCmd = &cobra.Command{
	Use:   "list",
	Short: "List work items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		query := url.Values{}
		if v, _ := cmd.Flags().GetString("state"); v != "" {
			query.Set("state", v)
		}
		if v, _ := cmd.Flags().GetString("assignee"); v != "" {
			query.Set("assignees", v)
		}
		if v, _ := cmd.Flags().GetString("label"); v != "" {
			query.Set("labels", v)
		}
		if v, _ := cmd.Flags().GetString("cycle"); v != "" {
			query.Set("cycle_id", v)
		}
		if v, _ := cmd.Flags().GetString("search"); v != "" {
			query.Set("search", v)
		}

		body, code, err := client.Do("GET", api.EndpointWorkItems(), nil, query)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

var wiGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a work item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointWorkItem(args[0]), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

var wiCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a work item",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		payload := map[string]interface{}{}
		title, _ := cmd.Flags().GetString("title")
		if title == "" {
			output.WriteError("validation", "--title is required")
			return nil
		}
		payload["name"] = title
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			payload["description_html"] = v
		}
		if v, _ := cmd.Flags().GetString("state"); v != "" {
			payload["state"] = v
		}
		if v, _ := cmd.Flags().GetString("priority"); v != "" {
			payload["priority"] = v
		}
		if v, _ := cmd.Flags().GetString("assignee"); v != "" {
			payload["assignees"] = strings.Split(v, ",")
		}
		if v, _ := cmd.Flags().GetString("labels"); v != "" {
			payload["labels_list"] = strings.Split(v, ",")
		}

		body, code, err := client.Do("POST", api.EndpointWorkItems(), payload, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

var wiUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a work item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		payload := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("title"); v != "" {
			payload["name"] = v
		}
		if v, _ := cmd.Flags().GetString("state"); v != "" {
			payload["state"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			payload["description_html"] = v
		}
		if v, _ := cmd.Flags().GetString("assignee"); v != "" {
			payload["assignees"] = strings.Split(v, ",")
		}
		if v, _ := cmd.Flags().GetString("priority"); v != "" {
			payload["priority"] = v
		}
		if v, _ := cmd.Flags().GetString("labels"); v != "" {
			payload["labels_list"] = strings.Split(v, ",")
		}

		body, code, err := client.Do("PATCH", api.EndpointWorkItem(args[0]), payload, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

var wiDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a work item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("DELETE", api.EndpointWorkItem(args[0]), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

var wiCommentListCmd = &cobra.Command{
	Use:   "list <work-item-id>",
	Short: "List comments on a work item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		body, code, err := client.Do("GET", api.EndpointWorkItemComments(args[0]), nil, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

var wiCommentAddCmd = &cobra.Command{
	Use:   "add <work-item-id>",
	Short: "Add a comment to a work item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
			exitCode = 3
			output.WriteError("missing_context", err.Error())
			return nil
		}
		content, _ := cmd.Flags().GetString("content")
		if content == "" {
			output.WriteError("validation", "--content is required")
			return nil
		}
		payload := map[string]string{"comment_html": content}
		body, code, err := client.Do("POST", api.EndpointWorkItemComments(args[0]), payload, nil)
		if err != nil {
			handleAPIError(code, err, body)
			return nil
		}
		handleWorkItemsResponse(body)
		return nil
	},
}

func init() {
	wiListCmd.Flags().String("state", "", "Filter by state name")
	wiListCmd.Flags().String("assignee", "", "Filter by assignee email")
	wiListCmd.Flags().String("label", "", "Filter by label name")
	wiListCmd.Flags().String("cycle", "", "Filter by cycle ID")
	wiListCmd.Flags().String("search", "", "Text search query")

	wiCreateCmd.Flags().String("title", "", "Work item title (required)")
	wiCreateCmd.Flags().String("description", "", "Description (HTML)")
	wiCreateCmd.Flags().String("state", "", "State name")
	wiCreateCmd.Flags().String("priority", "", "Priority: urgent, high, medium, low, none")
	wiCreateCmd.Flags().String("assignee", "", "Assignee email")
	wiCreateCmd.Flags().String("labels", "", "Comma-separated label names")

	wiUpdateCmd.Flags().String("title", "", "Work item title")
	wiUpdateCmd.Flags().String("state", "", "State name")
	wiUpdateCmd.Flags().String("description", "", "Description (HTML)")
	wiUpdateCmd.Flags().String("assignee", "", "Assignee email")
	wiUpdateCmd.Flags().String("priority", "", "Priority: urgent, high, medium, low, none")
	wiUpdateCmd.Flags().String("labels", "", "Comma-separated label names")

	wiCommentAddCmd.Flags().String("content", "", "Comment content (HTML, required)")

	wiCmd.AddCommand(wiListCmd, wiGetCmd, wiCreateCmd, wiUpdateCmd, wiDeleteCmd)

	commentCmd := &cobra.Command{Use: "comment", Short: "Manage comments"}
	commentCmd.AddCommand(wiCommentListCmd, wiCommentAddCmd)
	wiCmd.AddCommand(commentCmd)

	rootCmd.AddCommand(wiCmd)
}

func handleWorkItemsResponse(body []byte) {
	var rawList []json.RawMessage
	if err := json.Unmarshal(body, &rawList); err != nil {
		var single json.RawMessage
		if err := json.Unmarshal(body, &single); err != nil {
			output.WriteError("parse_error", "Failed to parse API response")
			return
		}
		items := make([]interface{}, 1)
		items[0] = single
		output.WriteItems(items)
		return
	}
	items := make([]interface{}, len(rawList))
	for i, r := range rawList {
		items[i] = r
	}
	output.WriteItems(items)
}


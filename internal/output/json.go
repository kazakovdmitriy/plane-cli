package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/makeplane/plane-cli/internal/models"
)

func WriteItems(items []interface{}) {
	out := models.Collection{Items: items, Total: len(items)}
	write(out)
}

func WriteItem(item interface{}) {
	out := models.Collection{Items: []interface{}{item}, Total: 1}
	write(out)
}

func WriteError(errorCode, message string) {
	out := models.APIError{Error: errorCode, Message: message}
	data, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(data))
	fmt.Fprintf(os.Stderr, "plane: %s\n", message)
}

func write(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "plane: marshal error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

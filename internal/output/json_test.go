package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/kazakovdmitriy/plane-cli/internal/models"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestWriteItems(t *testing.T) {
	out := captureStdout(func() {
		WriteItems([]interface{}{map[string]string{"id": "1", "name": "test"}})
	})
	var result models.Collection
	json.Unmarshal([]byte(out), &result)
	if result.Total != 1 {
		t.Errorf("Expected 1, got %d", result.Total)
	}
}

func TestWriteItem(t *testing.T) {
	out := captureStdout(func() {
		WriteItem(map[string]string{"id": "1"})
	})
	var result models.Collection
	json.Unmarshal([]byte(out), &result)
	if result.Total != 1 {
		t.Errorf("Expected 1, got %d", result.Total)
	}
}

func TestWriteEmptyItems(t *testing.T) {
	out := captureStdout(func() {
		WriteItems([]interface{}{})
	})
	var result models.Collection
	json.Unmarshal([]byte(out), &result)
	if result.Total != 0 || len(result.Items) != 0 {
		t.Errorf("Expected empty, got total=%d", result.Total)
	}
}

func TestWriteError(t *testing.T) {
	out := captureStdout(func() {
		WriteError("not_found", "Work item 123 not found")
	})
	var apiErr models.APIError
	json.Unmarshal([]byte(out), &apiErr)
	if apiErr.Error != "not_found" {
		t.Errorf("Expected not_found, got %s", apiErr.Error)
	}
}

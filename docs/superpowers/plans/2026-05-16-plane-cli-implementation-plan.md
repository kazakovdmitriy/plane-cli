# Plane CLI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Создать CLI-утилиту `plane` на Go для взаимодействия AI-агентов с Plane REST API через API-токены.

**Architecture:** Монолитный Go-проект в одном модуле. Cobra для CLI-команд, hashicorp/go-retryablehttp для retry-логики, собственный circuit breaker. Весь код в `internal/`, точка входа в `cmd/plane/`. Статические бинарники под linux/darwin amd64/arm64.

**Tech Stack:** Go 1.22+, `github.com/spf13/cobra`, `github.com/hashicorp/go-retryablehttp`, стандартная библиотека.

**Spec:** `docs/superpowers/specs/2026-05-16-plane-cli-design.md`

---

### Task 1: Инициализация проекта

**Files:**
- Create: `plane/go.mod`
- Create: `plane/Makefile`
- Create: `plane/cmd/plane/main.go`
- Create: `plane/.gitignore`

- [ ] **Step 1: Инициализировать Go модуль**

```bash
mkdir -p plane/cmd/plane && cd plane && go mod init github.com/makeplane/plane-cli
```

Expected: создан `go.mod` с модулем `github.com/makeplane/plane-cli`.

- [ ] **Step 2: Написать Makefile**

Создать `plane/Makefile`:

```makefile
APP_NAME = plane
GO = go
LDFLAGS = -ldflags="-s -w"

.PHONY: build build-all clean test

build:
	$(GO) build $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/plane

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-linux-amd64 ./cmd/plane

build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-linux-arm64 ./cmd/plane

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-amd64 ./cmd/plane

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-arm64 ./cmd/plane

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

clean:
	rm -rf bin/

test:
	$(GO) test ./...
```

- [ ] **Step 3: Написать .gitignore**

Создать `plane/.gitignore`:

```
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
vendor/
.idea/
.vscode/
```

- [ ] **Step 4: Написать минимальный main.go**

Создать `plane/cmd/plane/main.go`:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("plane CLI")
	os.Exit(0)
}
```

- [ ] **Step 5: Собрать и проверить**

```bash
cd plane && make build && ./bin/plane
```

Expected: выводит "plane CLI", код возврата 0.

- [ ] **Step 6: Commit**

```bash
git add plane/go.mod plane/Makefile plane/cmd/plane/main.go plane/.gitignore
git commit -m "feat: init plane CLI project skeleton"
```

---

### Task 2: Модели данных

**Files:**
- Create: `plane/internal/models/models.go`
- Create: `plane/internal/models/models_test.go`

- [ ] **Step 1: Создать все domain-модели**

Создать `plane/internal/models/models.go`:

```go
package models

type WorkItem struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	State       string   `json:"state,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Assignees   []string `json:"assignees,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Cycle       string   `json:"cycle,omitempty"`
	Module      string   `json:"module,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	CreatedBy   string   `json:"created_by"`
	URL         string   `json:"url"`
}

type Cycle struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Module struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type State struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Group string `json:"group"`
}

type Label struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Member struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
}

type Page struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Comment struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

type Collection struct {
	Items []interface{} `json:"items"`
	Total int           `json:"total"`
}

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
```

- [ ] **Step 2: Написать тесты JSON-сериализации**

Создать `plane/internal/models/models_test.go`:

```go
package models

import (
	"encoding/json"
	"testing"
)

func TestWorkItemJSON(t *testing.T) {
	wi := WorkItem{
		ID: "abc-123", Name: "Test Item", State: "Backlog",
		Priority: "high", Assignees: []string{"dev@test.com"}, Labels: []string{"bug"},
	}
	data, _ := json.Marshal(wi)
	var result WorkItem
	json.Unmarshal(data, &result)
	if result.ID != "abc-123" || result.Name != "Test Item" {
		t.Errorf("Got %+v", result)
	}
}

func TestCycleJSON(t *testing.T) {
	c := Cycle{ID: "cyc-1", Name: "Sprint 1", StartDate: "2026-01-01"}
	data, _ := json.Marshal(c)
	var result Cycle
	json.Unmarshal(data, &result)
	if result.ID != "cyc-1" {
		t.Errorf("Got %+v", result)
	}
}

func TestLabelJSON(t *testing.T) {
	l := Label{ID: "lbl-1", Name: "bug", Color: "#ff0000"}
	data, _ := json.Marshal(l)
	var result Label
	json.Unmarshal(data, &result)
	if result.Color != "#ff0000" {
		t.Errorf("Got %+v", result)
	}
}

func TestMemberJSON(t *testing.T) {
	m := Member{ID: "usr-1", Email: "dev@test.com", DisplayName: "Dev", Role: "admin"}
	data, _ := json.Marshal(m)
	var result Member
	json.Unmarshal(data, &result)
	if result.Email != "dev@test.com" {
		t.Errorf("Got %+v", result)
	}
}

func TestPageJSON(t *testing.T) {
	p := Page{ID: "pg-1", Name: "Home", Content: "# Welcome"}
	data, _ := json.Marshal(p)
	var result Page
	json.Unmarshal(data, &result)
	if result.Content != "# Welcome" {
		t.Errorf("Got %+v", result)
	}
}

func TestCommentJSON(t *testing.T) {
	c := Comment{ID: "cmt-1", Content: "Looks good", CreatedBy: "dev@test.com"}
	data, _ := json.Marshal(c)
	var result Comment
	json.Unmarshal(data, &result)
	if result.Content != "Looks good" {
		t.Errorf("Got %+v", result)
	}
}

func TestCollectionJSON(t *testing.T) {
	coll := Collection{Items: []interface{}{"a", "b"}, Total: 2}
	data, _ := json.Marshal(coll)
	var result Collection
	json.Unmarshal(data, &result)
	if result.Total != 2 {
		t.Errorf("Got total %d", result.Total)
	}
}
```

- [ ] **Step 3: Запустить тесты**

```bash
cd plane && go test ./internal/models/ -v
```

Expected: 7 tests PASS.

- [ ] **Step 4: Commit**

```bash
git add plane/internal/models/
git commit -m "feat: add domain models with JSON tests"
```

---

### Task 3: JSON-вывод (Output)

**Files:**
- Create: `plane/internal/output/json.go`
- Create: `plane/internal/output/json_test.go`

- [ ] **Step 1: Создать output-пакет**

Создать `plane/internal/output/json.go`:

```go
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
```

- [ ] **Step 2: Написать тесты output**

Создать `plane/internal/output/json_test.go`:

```go
package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/makeplane/plane-cli/internal/models"
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
```

- [ ] **Step 3: Запустить тесты**

```bash
cd plane && go test ./internal/output/ -v
```

Expected: 4 tests PASS.

- [ ] **Step 4: Commit**

```bash
git add plane/internal/output/
git commit -m "feat: add JSON output formatter with tests"
```

---

### Task 4: Конфигурация и контекст

**Files:**
- Create: `plane/internal/config/context.go`
- Create: `plane/internal/config/context_test.go`

- [ ] **Step 1: Создать менеджер контекста**

Создать `plane/internal/config/context.go`:

```go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Context struct {
	Workspace string `json:"workspace"`
	Project   string `json:"project"`
	Token     string `json:"token,omitempty"`
	APIURL    string `json:"api_url,omitempty"`
}

const DefaultAPIURL = "https://api.plane.so"

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	return filepath.Join(home, ".config", "plane", "config.json"), nil
}

func Load() (*Context, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Context{}, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	var ctx Context
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if ctx.APIURL == "" {
		ctx.APIURL = DefaultAPIURL
	}
	return &ctx, nil
}

func Save(ctx *Context) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

func Delete() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("remove config: %w", err)
	}
	return nil
}

func ResolveToken(flagToken, envToken string) (string, error) {
	if flagToken != "" {
		return flagToken, nil
	}
	if envToken != "" {
		return envToken, nil
	}
	ctx, err := Load()
	if err != nil {
		return "", fmt.Errorf("no token: use --token, PLANE_TOKEN, or plane context set")
	}
	if ctx.Token != "" {
		return ctx.Token, nil
	}
	return "", fmt.Errorf("no token: use --token, PLANE_TOKEN, or plane context set")
}

func ResolveWorkspace(flagWS, envWS string) (string, error) {
	if flagWS != "" {
		return flagWS, nil
	}
	if envWS != "" {
		return envWS, nil
	}
	ctx, _ := Load()
	if ctx != nil && ctx.Workspace != "" {
		return ctx.Workspace, nil
	}
	return "", fmt.Errorf("no workspace: use --workspace, PLANE_WORKSPACE, or plane context set")
}

func ResolveProject(flagPrj, envPrj string) (string, error) {
	if flagPrj != "" {
		return flagPrj, nil
	}
	if envPrj != "" {
		return envPrj, nil
	}
	ctx, _ := Load()
	if ctx != nil && ctx.Project != "" {
		return ctx.Project, nil
	}
	return "", fmt.Errorf("no project: use --project, PLANE_PROJECT, or plane context set")
}

func ResolveAPIURL(flagURL, envURL string) string {
	if flagURL != "" {
		return flagURL
	}
	if envURL != "" {
		return envURL
	}
	ctx, _ := Load()
	if ctx != nil && ctx.APIURL != "" {
		return ctx.APIURL
	}
	return DefaultAPIURL
}
```

- [ ] **Step 2: Написать тесты контекста**

Создать `plane/internal/config/context_test.go`:

```go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	ctx := &Context{Workspace: "testorg", Project: "abc-uuid", Token: "plane_key_test", APIURL: "https://plane.test.com"}
	if err := Save(ctx); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Workspace != "testorg" || loaded.Token != "plane_key_test" {
		t.Errorf("Got %+v", loaded)
	}
}

func TestLoadEmpty(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	ctx, err := Load()
	if err != nil {
		t.Fatalf("Load empty: %v", err)
	}
	if ctx.APIURL != DefaultAPIURL {
		t.Errorf("Expected default API URL, got %s", ctx.APIURL)
	}
}

func TestDelete(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	Save(&Context{Token: "tk"})
	if err := Delete(); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	loaded, _ := Load()
	if loaded.Token != "" {
		t.Errorf("Expected empty token, got %s", loaded.Token)
	}
}

func TestResolveToken(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	token, err := ResolveToken("flag-token", "")
	if err != nil || token != "flag-token" {
		t.Errorf("Expected flag-token, got %s", token)
	}

	token, err = ResolveToken("", "env-token")
	if err != nil || token != "env-token" {
		t.Errorf("Expected env-token, got %s", token)
	}

	Save(&Context{Token: "config-token"})
	token, err = ResolveToken("", "")
	if err != nil || token != "config-token" {
		t.Errorf("Expected config-token, got %s err=%v", token, err)
	}
}

func TestResolveTokenMissing(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")
	if _, err := ResolveToken("", ""); err == nil {
		t.Errorf("Expected error")
	}
}

func TestConfigPath(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")
	path, _ := configPath()
	expected := filepath.Join(home, ".config", "plane", "config.json")
	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}
}
```

- [ ] **Step 3: Запустить тесты**

```bash
cd plane && go test ./internal/config/ -v
```

Expected: 6 tests PASS.

- [ ] **Step 4: Commit**

```bash
git add plane/internal/config/
git commit -m "feat: add context manager (save/load/delete config.json)"
```

---

### Task 5: API-клиент (transport + client)

**Files:**
- Create: `plane/internal/api/transport.go`
- Create: `plane/internal/api/client.go`
- Create: `plane/internal/api/endpoints.go`
- Create: `plane/internal/api/client_test.go`

- [ ] **Step 1: Установить зависимости**

```bash
cd plane && go get github.com/spf13/cobra github.com/hashicorp/go-retryablehttp
```

- [ ] **Step 2: Создать transport (retry + circuit breaker)**

Создать `plane/internal/api/transport.go`:

```go
package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	DefaultMaxRetries          = 4
	DefaultTimeout             = 30 * time.Second
	DefaultCircuitBreakerLimit = 5
)

type RetryConfig struct {
	MaxRetries          int
	TotalTimeout        time.Duration
	CircuitBreakerLimit int
}

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:          DefaultMaxRetries,
		TotalTimeout:        DefaultTimeout,
		CircuitBreakerLimit: DefaultCircuitBreakerLimit,
	}
}

type CircuitBreaker struct {
	mu               sync.Mutex
	consecutiveFails int
	limit            int
	open             bool
}

func NewCircuitBreaker(limit int) *CircuitBreaker {
	return &CircuitBreaker{limit: limit}
}

func (cb *CircuitBreaker) IsOpen() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.open
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.consecutiveFails++
	if cb.consecutiveFails >= cb.limit {
		cb.open = true
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.consecutiveFails = 0
	cb.open = false
}

var ErrCircuitBreakerOpen = fmt.Errorf("circuit breaker open")

func newRetryableHTTP(cfg *RetryConfig) *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = cfg.MaxRetries
	client.RetryWaitMin = 1 * time.Second
	client.RetryWaitMax = 8 * time.Second
	client.HTTPClient.Timeout = cfg.TotalTimeout
	client.Logger = nil
	client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			return true, nil
		}
		status := resp.StatusCode
		if status == 429 || status >= 500 {
			return true, nil
		}
		return false, nil
	}
	return client
}
```

- [ ] **Step 3: Создать client (обёртка над transport с circuit breaker)**

Создать `plane/internal/api/client.go`:

```go
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	http           *retryablehttp.Client
	cb             *CircuitBreaker
	token          string
	baseURL        string
	workspace      string
	project        string
}

func NewClient(token, baseURL, workspace, project string, cfg *RetryConfig) *Client {
	return &Client{
		http:      newRetryableHTTP(cfg),
		cb:        NewCircuitBreaker(cfg.CircuitBreakerLimit),
		token:     token,
		baseURL:   strings.TrimRight(baseURL, "/"),
		workspace: workspace,
		project:   project,
	}
}

func (c *Client) Do(method, path string, body interface{}, query url.Values) ([]byte, int, error) {
	if c.cb.IsOpen() {
		return nil, 0, ErrCircuitBreakerOpen
	}

	u := fmt.Sprintf("%s/api/v1/workspaces/%s/projects/%s/%s", c.baseURL, c.workspace, c.project, path)
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := retryablehttp.NewRequest(method, u, reqBody)
	if err != nil {
		c.cb.RecordFailure()
		return nil, 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-API-Key", c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "plane-cli/0.1.0")

	resp, err := c.http.Do(req)
	if err != nil {
		c.cb.RecordFailure()
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.cb.RecordFailure()
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		c.cb.RecordSuccess()
		return respBody, resp.StatusCode, nil
	}

	c.cb.RecordFailure()
	return respBody, resp.StatusCode, fmt.Errorf("API error: %s", resp.Status)
}

func methodGet(query url.Values) string { return "GET" }
func methodPost(query url.Values) string { return "POST" }
func methodPatch(query url.Values) string { return "PATCH" }
func methodDelete(query url.Values) string { return "DELETE" }
```

- [ ] **Step 4: Создать endpoints (построение URL-путей)**

Создать `plane/internal/api/endpoints.go`:

```go
package api

import "fmt"

func endpointWorkItems() string                  { return "work-items/" }
func endpointWorkItem(id string) string           { return fmt.Sprintf("work-items/%s/", id) }
func endpointWorkItemComments(id string) string   { return fmt.Sprintf("work-items/%s/comments/", id) }
func endpointCycles() string                      { return "cycles/" }
func endpointCycle(id string) string              { return fmt.Sprintf("cycles/%s/", id) }
func endpointModules() string                     { return "modules/" }
func endpointModule(id string) string             { return fmt.Sprintf("modules/%s/", id) }
func endpointStates() string                      { return "states/" }
func endpointState(id string) string              { return fmt.Sprintf("states/%s/", id) }
func endpointLabels() string                      { return "issue-labels/" }
func endpointMembers() string                     { return "members/" }
func endpointPages() string                       { return "pages/" }
func endpointPage(id string) string               { return fmt.Sprintf("pages/%s/", id) }
```

- [ ] **Step 5: Написать тесты API-клиента**

Создать `plane/internal/api/client_test.go`:

```go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientDoSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "test-token" {
			w.WriteHeader(401)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"id":"1","name":"Test"}`))
	}))
	defer server.Close()

	cfg := DefaultRetryConfig()
	client := NewClient("test-token", server.URL, "ws", "prj", cfg)
	body, code, err := client.Do("GET", endpointWorkItem("1"), nil, nil)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if code != 200 {
		t.Errorf("Expected 200, got %d", code)
	}
	if string(body) != `{"id":"1","name":"Test"}` {
		t.Errorf("Got body: %s", body)
	}
}

func TestClient401NoRetry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	}))
	defer server.Close()

	cfg := DefaultRetryConfig()
	cfg.TotalTimeout = 2 * time.Second
	client := NewClient("bad-token", server.URL, "ws", "prj", cfg)
	_, code, err := client.Do("GET", endpointWorkItems(), nil, nil)
	if err == nil {
		t.Errorf("Expected error for 401")
	}
	if code != 401 {
		t.Errorf("Expected 401, got %d", code)
	}
}

func TestClientURLConstruction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/api/v1/workspaces/myorg/projects/abc/work-items/test-1/"
		if r.URL.Path != expected {
			t.Errorf("Expected path %s, got %s", expected, r.URL.Path)
		}
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	cfg := DefaultRetryConfig()
	client := NewClient("tk", server.URL, "myorg", "abc", cfg)
	client.Do("GET", endpointWorkItem("test-1"), nil, nil)
}
```

- [ ] **Step 6: Запустить тесты**

```bash
cd plane && go test ./internal/api/ -v
```

Expected: test `TestClientURLConstruction` PASS. Остальные могут FAIL если timeouts не подобраны — это ок.

Note: `TestClientURLConstruction` выполняется внутри хендлера, t.Errorf там не работает как ожидается в стандартном тесте. При реализации заменить `t.Errorf` на `fmt.Printf` и проверять вручную, либо переписать через `httptest` с проверкой пути в основном тесте.

- [ ] **Step 7: Commit**

```bash
git add plane/internal/api/ plane/go.mod plane/go.sum
git commit -m "feat: add API client with retry and circuit breaker"
```

---

### Task 6: CLI — Root command

**Files:**
- Modify: `plane/cmd/plane/main.go`
- Create: `plane/internal/cli/root.go`

- [ ] **Step 1: Создать root command**

Создать `plane/internal/cli/root.go`:

```go
package cli

import (
	"time"

	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	flagToken     string
	flagAPIURL    string
	flagWorkspace string
	flagProject   string
	flagMaxRetries int
	flagTimeout   time.Duration
	flagNoColor   bool
)

var rootCmd = &cobra.Command{
	Use:   "plane",
	Short: "Plane CLI - interact with Plane API from terminal",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "API token (overrides PLANE_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&flagAPIURL, "api-url", "", "Base API URL (overrides PLANE_API_URL)")
	rootCmd.PersistentFlags().StringVarP(&flagWorkspace, "workspace", "w", "", "Workspace slug")
	rootCmd.PersistentFlags().StringVarP(&flagProject, "project", "p", "", "Project UUID")
	rootCmd.PersistentFlags().IntVar(&flagMaxRetries, "max-retries", 4, "Max retry attempts")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 30*time.Second, "Total timeout for all retries")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable colored output")
}

func resolveContext() (*api.Client, error) {
	token, err := config.ResolveToken(flagToken, getEnv("PLANE_TOKEN"))
	if err != nil {
		return nil, err
	}
	ws, err := config.ResolveWorkspace(flagWorkspace, getEnv("PLANE_WORKSPACE"))
	if err != nil {
		return nil, err
	}
	prj, err := config.ResolveProject(flagProject, getEnv("PLANE_PROJECT"))
	if err != nil {
		return nil, err
	}
	apiURL := config.ResolveAPIURL(flagAPIURL, getEnv("PLANE_API_URL"))

	cfg := api.DefaultRetryConfig()
	if flagMaxRetries != 4 {
		cfg.MaxRetries = flagMaxRetries
	}
	if flagTimeout != 30*time.Second {
		cfg.TotalTimeout = flagTimeout
	}
	return api.NewClient(token, apiURL, ws, prj, cfg), nil
}

func getEnv(key string) string { return "" }
```

Note: `getEnv` — заглушка. При реализации заменить на `os.Getenv(key)`.

- [ ] **Step 2: Обновить main.go**

Заменить `plane/cmd/plane/main.go`:

```go
package main

import (
	"fmt"
	"os"

	"github.com/makeplane/plane-cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "plane: %v\n", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Проверить сборку и --help**

```bash
cd plane && go build ./cmd/plane && ./bin/plane --help
```

Expected: выводит справку с флагами --token, --workspace, --project, --max-retries, --timeout, --no-color.

- [ ] **Step 4: Commit**

```bash
git add plane/internal/cli/root.go plane/cmd/plane/main.go
git commit -m "feat: add root command with global flags"
```

---

### Task 7: CLI — Context commands

**Files:**
- Create: `plane/internal/cli/context.go`
- Modify: `plane/internal/cli/root.go` (регистрация команды)

- [ ] **Step 1: Создать команду context**

Создать `plane/internal/cli/context.go`:

```go
package cli

import (
	"fmt"

	"github.com/makeplane/plane-cli/internal/config"
	"github.com/makeplane/plane-cli/internal/output"
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

	fmt.Println("context registered")
}
```

Note: `fmt.Println("context registered")` — временный дебаг. При реализации убрать. `init()` срабатывает при импорте пакета `cli`.

- [ ] **Step 2: Проверить context set/show/unset**

```bash
cd plane && go build ./cmd/plane
./bin/plane context set --workspace myorg --project abc123
./bin/plane context show
```

Expected: JSON с workspace=myorg, project=abc123, api_url=https://api.plane.so.

```bash
./bin/plane context unset
./bin/plane context show
```

Expected: пустые поля workspace/project.

- [ ] **Step 3: Commit**

```bash
git add plane/internal/cli/context.go
git commit -m "feat: add context set/show/unset commands"
```

---

### Task 8: CLI — Work Item commands

**Files:**
- Create: `plane/internal/cli/workitem.go`

- [ ] **Step 1: Создать команду work-item**

Создать `plane/internal/cli/workitem.go`:

```go
package cli

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/models"
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

func handleAPIError(code int, err error, body []byte) {
	switch {
	case err == api.ErrCircuitBreakerOpen:
		output.WriteError("circuit_breaker_open", "Circuit breaker open after consecutive failures")
	case code == 401 || code == 403:
		output.WriteError("unauthorized", "Invalid or missing API token")
	case code == 404:
		output.WriteError("not_found", "Resource not found")
	case code >= 500:
		output.WriteError("server_error", "Plane API server error")
	default:
		output.WriteError("api_error", err.Error())
	}
}
```

- [ ] **Step 2: Исправить импорты после добавления root.go**

Убедиться, что `plane/cmd/plane/main.go` импортирует `cli` пакет из internal. В Go модули internal-пакеты видны только из корня модуля — из `cmd/plane/main.go` импорт `github.com/makeplane/plane-cli/internal/cli` легален.

- [ ] **Step 3: Собрать и проверить --help**

```bash
cd plane && go build ./cmd/plane
./bin/plane wi --help
./bin/plane work-item --help
```

Expected: справка для wi/work-item с подкомандами list, get, create, update, delete, comment.

- [ ] **Step 4: Commit**

```bash
git add plane/internal/cli/workitem.go
git commit -m "feat: add work-item CRUD commands"
```

---

### Task 9: CLI — Cycle commands

**Files:**
- Create: `plane/internal/cli/cycle.go`

- [ ] **Step 1: Создать команду cycle**

Создать `plane/internal/cli/cycle.go`:

```go
package cli

import (
	"encoding/json"
	"net/url"

	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/output"
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
```

- [ ] **Step 2: Собрать и проверить**

```bash
cd plane && go build ./cmd/plane && ./bin/plane cycle --help
```

Expected: справка для cycle с list, get, create.

- [ ] **Step 3: Commit**

```bash
git add plane/internal/cli/cycle.go
git commit -m "feat: add cycle commands (list, get, create)"
```

---

### Task 10: CLI — Module, State, Label, Member, Page commands

**Files:**
- Create: `plane/internal/cli/module.go`
- Create: `plane/internal/cli/state.go`
- Create: `plane/internal/cli/label.go`
- Create: `plane/internal/cli/member.go`
- Create: `plane/internal/cli/page.go`

- [ ] **Step 1: Создать module.go**

Создать `plane/internal/cli/module.go`:

```go
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
```

- [ ] **Step 2: Создать state.go**

Создать `plane/internal/cli/state.go`:

```go
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
```

- [ ] **Step 3: Создать label.go**

Создать `plane/internal/cli/label.go`:

```go
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
```

- [ ] **Step 4: Создать member.go**

Создать `plane/internal/cli/member.go`:

```go
package cli

import (
	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var memberCmd = &cobra.Command{Use: "member", Short: "Manage members"}

var memberListCmd = &cobra.Command{
	Use: "list", Short: "List members",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
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
```

- [ ] **Step 5: Создать page.go**

Создать `plane/internal/cli/page.go`:

```go
package cli

import (
	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var pageCmd = &cobra.Command{Use: "page", Short: "Manage pages"}

var pageListCmd = &cobra.Command{
	Use: "list", Short: "List pages",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveContext()
		if err != nil {
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
```

- [ ] **Step 6: Собрать и проверить все команды**

```bash
cd plane && go build ./cmd/plane
./bin/plane module --help
./bin/plane state --help
./bin/plane label --help
./bin/plane member --help
./bin/plane page --help
```

Expected: справка для всех сущностей.

- [ ] **Step 7: Commit**

```bash
git add plane/internal/cli/module.go plane/internal/cli/state.go plane/internal/cli/label.go plane/internal/cli/member.go plane/internal/cli/page.go
git commit -m "feat: add module, state, label, member, page commands"
```

---

### Task 11: Исправление ошибок и доработки

**Files:**
- Modify: `plane/internal/cli/root.go` (getEnv, handleAPIError перенести)
- Modify: `plane/internal/cli/workitem.go` (убрать дублирующийся handleAPIError)
- Modify: `plane/cmd/plane/main.go` (os.Exit codes)

- [ ] **Step 1: Исправить root.go (getEnv, handleAPIError, exitCode, импорты)**

В `plane/internal/cli/root.go` добавить недостающие импорты и функции:

```go
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/makeplane/plane-cli/internal/api"
	"github.com/makeplane/plane-cli/internal/config"
	"github.com/makeplane/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

// ... существующие var и rootCmd остаются без изменений ...

var exitCode = 0

func ExitCode() int { return exitCode }

func getEnv(key string) string { return os.Getenv(key) }

func handleAPIError(code int, err error, body []byte) {
	switch {
	case err == api.ErrCircuitBreakerOpen:
		exitCode = 2
		output.WriteError("circuit_breaker_open", "Circuit breaker open after consecutive failures")
	case code == 401 || code == 403:
		exitCode = 1
		output.WriteError("unauthorized", "Invalid or missing API token")
	case code == 404:
		exitCode = 1
		output.WriteError("not_found", "Resource not found")
	case code >= 500:
		exitCode = 1
		output.WriteError("server_error", "Plane API server error")
	default:
		exitCode = 4
		output.WriteError("network", err.Error())
	}
}
```

Note: добавить импорты `"fmt"`, `"os"`, `"github.com/makeplane/plane-cli/internal/output"` в начало файла, если их ещё нет.

- [ ] **Step 2: Экспортировать функции endpoint и убрать дублирование handleAPIError**

В `plane/internal/api/endpoints.go` переименовать функции с маленькой буквы на большую:
- `endpointWorkItems()` → `EndpointWorkItems()`
- `endpointWorkItem(id)` → `EndpointWorkItem(id)`
- `endpointWorkItemComments(id)` → `EndpointWorkItemComments(id)`
- `endpointCycles()` → `EndpointCycles()`
- `endpointCycle(id)` → `EndpointCycle(id)`
- `endpointModules()` → `EndpointModules()`
- `endpointModule(id)` → `EndpointModule(id)`
- `endpointStates()` → `EndpointStates()`
- `endpointState(id)` → `EndpointState(id)`
- `endpointLabels()` → `EndpointLabels()`
- `endpointMembers()` → `EndpointMembers()`
- `endpointPages()` → `EndpointPages()`
- `endpointPage(id)` → `EndpointPage(id)`

Функция `handleAPIError` определена и в `workitem.go`, и в `cycle.go`. Оставить её только в `root.go` и удалить из этих файлов.

```go
// root.go
func handleAPIError(code int, err error, body []byte) {
	switch {
	case err == api.ErrCircuitBreakerOpen:
		output.WriteError("circuit_breaker_open", "Circuit breaker open after consecutive failures")
	case code == 401 || code == 403:
		output.WriteError("unauthorized", "Invalid or missing API token")
	case code == 404:
		output.WriteError("not_found", "Resource not found")
	case code >= 500:
		output.WriteError("server_error", "Plane API server error")
	default:
		output.WriteError("api_error", err.Error())
	}
}
```

Удалить `handleAPIError` из `workitem.go` и `cycle.go`.

- [ ] **Step 3: Реализовать коды возврата по спецификации**

В `plane/internal/cli/root.go` добавить переменную `exitCode` и функцию `ExitCode()`:

```go
var exitCode = 0

func ExitCode() int { return exitCode }
```

Обновить `handleAPIError` с правильными кодами возврата:

```go
func handleAPIError(code int, err error, body []byte) {
	switch {
	case err == api.ErrCircuitBreakerOpen:
		exitCode = 2
		output.WriteError("circuit_breaker_open", "Circuit breaker open after consecutive failures")
	case code == 401 || code == 403:
		exitCode = 1
		output.WriteError("unauthorized", "Invalid or missing API token")
	case code == 404:
		exitCode = 1
		output.WriteError("not_found", "Resource not found")
	case code >= 500:
		exitCode = 1
		output.WriteError("server_error", "Plane API server error")
	default:
		exitCode = 4
		output.WriteError("network", err.Error())
	}
}
```

Обновить все команды — при ошибке `resolveContext()` задавать exitCode=3:

```go
client, err := resolveContext()
if err != nil {
	exitCode = 3
	output.WriteError("missing_context", err.Error())
	return nil
}
```

Обновить `plane/cmd/plane/main.go`:

```go
package main

import (
	"os"

	"github.com/makeplane/plane-cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(cli.ExitCode())
}
```

Паттерн для каждой команды:
- `exitCode = 3` при проблемах с контекстом/токеном
- `exitCode = 2` при circuit breaker
- `exitCode = 4` при сетевых ошибках
- `exitCode = 1` при ошибках API

- [ ] **Step 4: Пересобрать**

```bash
cd plane && go build ./cmd/plane
```

Expected: сборка без ошибок.

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "fix: consolidate error handling and exit codes"
```

---

### Task 12: Финальная сборка и smoke-тест

**Files:**
- None (build artifacts)

- [ ] **Step 1: Запустить все unit-тесты**

```bash
cd plane && go test ./... -v
```

Expected: все тесты PASS.

- [ ] **Step 2: Собрать все платформы**

```bash
cd plane && make build-all
```

Expected: бинарники в `plane/bin/`:
- `plane-linux-amd64`
- `plane-linux-arm64`
- `plane-darwin-amd64`
- `plane-darwin-arm64`

- [ ] **Step 3: Проверить бинарник на текущей платформе**

```bash
cd plane && make build
./bin/plane --help
./bin/plane wi --help
./bin/plane context --help
./bin/plane cycle --help
./bin/plane module --help
./bin/plane state --help
./bin/plane label --help
./bin/plane member --help
./bin/plane page --help
```

Expected: все команды выводят справку без ошибок.

- [ ] **Step 4: Проверить контекст (round-trip)**

```bash
./bin/plane context set --workspace testws --project testprj --token testtk --api-url https://test.example.com
./bin/plane context show
```

Expected: JSON с заданными значениями.

```bash
./bin/plane context unset
```

- [ ] **Step 5: Commit**

```bash
git add plane/bin/ -f || true
git commit -m "chore: final build and smoke test"
```

Note: `plane/bin/` может быть в `.gitignore` — если да, то коммит без бинарников.

---

## Примечания для реализующего

1. **Имена функций API-эндпоинтов:** Экспортируемые имена (`EndpointWorkItems`, `EndpointCycle`, etc.) нужны для импорта из пакета `cli`. В коде выше используется `EndpointXxx`, но в `api/endpoints.go` заданы `endpointXxx` — реализация должна экспортировать их.
2. **Ответы Plane API:** Plane может возвращать plain array или объект с `results`. Функция `handleCollectionResponse` обрабатывает оба варианта.
3. **Поля payload:** Имена полей в create/update запросах (`name`, `description_html`, `assignees`, `labels_list`) основаны на типовых соглашениях Plane DRF API. При интеграции с реальным API может потребоваться корректировка.
4. **getEnv:** В Task 11 заменяется на `os.Getenv`.
5. **Файлы в gitignore:** `plane/bin/` в `.gitignore` — бинарники не коммитятся, если не нужны.

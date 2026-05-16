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

func TestReadAgentTokens(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".plane-tokens")

	if err := os.WriteFile(path, []byte("pm=token-pm\ntl=token-tl\n# comment\ndev=token-dev\n\nqa=token-qa\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	tokens, err := ReadAgentTokens(path)
	if err != nil {
		t.Fatalf("ReadAgentTokens: %v", err)
	}
	if tokens["pm"] != "token-pm" {
		t.Errorf("Expected token-pm, got %s", tokens["pm"])
	}
	if tokens["tl"] != "token-tl" {
		t.Errorf("Expected token-tl, got %s", tokens["tl"])
	}
	if tokens["qa"] != "token-qa" {
		t.Errorf("Expected token-qa, got %s", tokens["qa"])
	}
	if _, ok := tokens["cr"]; ok {
		t.Errorf("Expected no cr token")
	}
}

func TestReadAgentTokensMissing(t *testing.T) {
	_, err := ReadAgentTokens("/nonexistent/path/.plane-tokens")
	if err == nil {
		t.Errorf("Expected error for missing file")
	}
}

func TestReadAgentTokensEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".plane-tokens")

	if err := os.WriteFile(path, []byte("# only comment\n\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := ReadAgentTokens(path)
	if err == nil {
		t.Errorf("Expected error for empty tokens file")
	}
}

func TestResolveTokenWithAgent(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	dir := t.TempDir()
	tokensPath := filepath.Join(dir, ".plane-tokens")
	if err := os.WriteFile(tokensPath, []byte("pm=token-pm\ndev=token-dev\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	os.Setenv("PLANE_AGENT", "pm")
	defer os.Unsetenv("PLANE_AGENT")

	token, err := resolveToken("", "", tokensPath)
	if err != nil || token != "token-pm" {
		t.Errorf("Expected token-pm, got %s err=%v", token, err)
	}
}

func TestResolveTokenAgentRoleNotFound(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	dir := t.TempDir()
	tokensPath := filepath.Join(dir, ".plane-tokens")
	if err := os.WriteFile(tokensPath, []byte("pm=token-pm\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	os.Setenv("PLANE_AGENT", "qa")
	defer os.Unsetenv("PLANE_AGENT")

	Save(&Context{Token: "config-token"})
	token, err := resolveToken("", "", tokensPath)
	if err != nil || token != "config-token" {
		t.Errorf("Expected fallback to config-token, got %s err=%v", token, err)
	}
}

func TestResolveTokenAgentFileMissing(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	os.Setenv("PLANE_AGENT", "pm")
	defer os.Unsetenv("PLANE_AGENT")

	Save(&Context{Token: "config-token"})
	token, err := resolveToken("", "", "/nonexistent/path/.plane-tokens")
	if err != nil || token != "config-token" {
		t.Errorf("Expected fallback to config-token, got %s err=%v", token, err)
	}
}

func TestResolveTokenAgentPriority(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	dir := t.TempDir()
	tokensPath := filepath.Join(dir, ".plane-tokens")
	if err := os.WriteFile(tokensPath, []byte("pm=token-pm\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	os.Setenv("PLANE_AGENT", "pm")
	defer os.Unsetenv("PLANE_AGENT")

	Save(&Context{Token: "config-token"})

	token, _ := resolveToken("flag-token", "", tokensPath)
	if token != "flag-token" {
		t.Errorf("Expected flag-token, got %s", token)
	}

	token, _ = resolveToken("", "env-token", tokensPath)
	if token != "env-token" {
		t.Errorf("Expected env-token, got %s", token)
	}

	token, _ = resolveToken("", "", tokensPath)
	if token != "token-pm" {
		t.Errorf("Expected token-pm, got %s", token)
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

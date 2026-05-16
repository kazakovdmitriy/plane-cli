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

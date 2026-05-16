package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Context struct {
	Workspace string `json:"workspace"`
	Project   string `json:"project"`
	Token     string `json:"token,omitempty"`
	APIURL    string `json:"api_url,omitempty"`
}

const DefaultAPIURL = "https://api.plane.so"

const AgentTokensFile = "/workspace/config/.plane-tokens"

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
			return &Context{APIURL: DefaultAPIURL}, nil
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

func ReadAgentTokens(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tokens := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		role := strings.TrimSpace(parts[0])
		token := strings.TrimSpace(parts[1])
		if role != "" && token != "" {
			tokens[role] = token
		}
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no valid tokens in %s", path)
	}
	return tokens, nil
}

func ResolveToken(flagToken, envToken string) (string, error) {
	return resolveToken(flagToken, envToken, AgentTokensFile)
}

func resolveToken(flagToken, envToken, agentTokensFile string) (string, error) {
	if flagToken != "" {
		return flagToken, nil
	}
	if envToken != "" {
		return envToken, nil
	}

	agent := os.Getenv("PLANE_AGENT")
	if agent != "" {
		if tokens, err := ReadAgentTokens(agentTokensFile); err == nil {
			if tok, ok := tokens[agent]; ok {
				return tok, nil
			}
		}
	}

	ctx, err := Load()
	if err != nil {
		return "", fmt.Errorf("no token: use --token, PLANE_TOKEN, PLANE_AGENT, or plane context set")
	}
	if ctx.Token != "" {
		return ctx.Token, nil
	}
	return "", fmt.Errorf("no token: use --token, PLANE_TOKEN, PLANE_AGENT, or plane context set")
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

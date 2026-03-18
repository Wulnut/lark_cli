/*
 * @Author: wulnut carepdime@gmail.com
 * @Date: 2026-03-19 00:16:53
 * @LastEditors: wulnut carepdime@gmail.com
 * @LastEditTime: 2026-03-19 01:39:37
 * @FilePath: /lark_cli/internal/config/config.go
 * @Description: 配置文件
 */
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const DefaultBaseURL = "https://project.feishu.cn"

type Config struct {
	BaseURL       string
	PluginID      string
	PluginSecret  string
	SessionPath   string
	HTTPTimeout   time.Duration
	RefreshLeeway time.Duration
}

type fileConfig struct {
	BaseURL      string `json:"base_url"`
	PluginID     string `json:"plugin_id"`
	PluginSecret string `json:"plugin_secret"`
	SessionPath  string `json:"session_path"`
}

func Load() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		BaseURL:       DefaultBaseURL,
		SessionPath:   filepath.Join(home, ".lark", "session.json"),
		HTTPTimeout:   15 * time.Second,
		RefreshLeeway: 10 * time.Minute,
	}

	filePath := filepath.Join(home, ".lark", "config.json")
	fc, err := loadFileConfig(filePath)
	if err != nil {
		return Config{}, err
	}

	if fc.BaseURL != "" {
		cfg.BaseURL = fc.BaseURL
	}
	if fc.PluginID != "" {
		cfg.PluginID = fc.PluginID
	}
	if fc.PluginSecret != "" {
		cfg.PluginSecret = fc.PluginSecret
	}
	if fc.SessionPath != "" {
		cfg.SessionPath = fc.SessionPath
	}

	applyEnvOverrides(&cfg)

	return cfg, nil
}

// Deprecated: use Load for defaults + config file + env overrides.
func LoadFromEnv() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		BaseURL:       DefaultBaseURL,
		SessionPath:   filepath.Join(home, ".lark", "session.json"),
		HTTPTimeout:   15 * time.Second,
		RefreshLeeway: 10 * time.Minute,
	}

	applyEnvOverrides(&cfg)

	return cfg, nil
}

func loadFileConfig(path string) (fileConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fileConfig{}, nil
		}
		return fileConfig{}, fmt.Errorf("read config file %s: %w", path, err)
	}

	var cfg fileConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("parse config file %s: %w", path, err)
	}

	return cfg, nil
}

func applyEnvOverrides(c *Config) {
	if v := os.Getenv("LARK_BASE_URL"); v != "" {
		c.BaseURL = v
	}
	if v := os.Getenv("LARK_PLUGIN_ID"); v != "" {
		c.PluginID = v
	}
	if v := os.Getenv("LARK_PLUGIN_SECRET"); v != "" {
		c.PluginSecret = v
	}
	if v := os.Getenv("LARK_SESSION_PATH"); v != "" {
		c.SessionPath = v
	}
}

func (c Config) ValidateForPluginToken() error {
	if c.BaseURL == "" {
		return errors.New("LARK_BASE_URL is required")
	}
	if c.PluginID == "" {
		return errors.New("LARK_PLUGIN_ID is required")
	}
	if c.PluginSecret == "" {
		return errors.New("LARK_PLUGIN_SECRET is required")
	}
	return nil
}

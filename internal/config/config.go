/*
 * @Author: wulnut carepdime@gmail.com
 * @Date: 2026-03-19 00:16:53
 * @LastEditors: wulnut carepdime@gmail.com
 * @LastEditTime: 2026-03-19 00:47:45
 * @FilePath: /lark_cli/internal/config/config.go
 * @Description: 配置文件
 */
package config

import (
	"errors"
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

func LoadFromEnv() (Config, error) {
	c := Config{
		BaseURL:       os.Getenv("LARK_BASE_URL"),
		PluginID:      os.Getenv("LARK_PLUGIN_ID"),
		PluginSecret:  os.Getenv("LARK_PLUGIN_SECRET"),
		SessionPath:   os.Getenv("LARK_SESSION_PATH"),
		HTTPTimeout:   15 * time.Second,
		RefreshLeeway: 10 * time.Minute,
	}
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.SessionPath == "" {
		d, err := os.UserConfigDir()
		if err != nil {
			return Config{}, err
		}
		c.SessionPath = filepath.Join(d, "lark", "session.json")
	}
	return c, nil
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
package config

import (
	"os"
	"path/filepath"
)

// SubmoduleConfig 定义单个 submodule 的配置
type SubmoduleConfig struct {
	Name    string `json:"name"`
	Repo    string `json:"repo"`
	Type    string `json:"type"`    // service, client, specs, tools
	Product string `json:"product"` // lingbo, inspirai, independent
}

// Config 定义 sm 工具的配置
type Config struct {
	SubmodulesDir string            `json:"submodules_dir"`
	Submodules    []SubmoduleConfig `json:"submodules"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		SubmodulesDir: ".submodules",
		Submodules: []SubmoduleConfig{
			{Name: "lingbo-desktop", Repo: "git@github.com:inspirai-store/lingbo-desktop.git", Type: "client", Product: "lingbo"},
			{Name: "lingbo-web", Repo: "git@github.com:inspirai-store/lingbo-web.git", Type: "client", Product: "lingbo"},
			{Name: "inspirai-user", Repo: "git@github.com:inspirai-store/inspirai-user.git", Type: "service", Product: "inspirai"},
			{Name: "inspirai-admin", Repo: "git@github.com:inspirai-store/inspirai-admin.git", Type: "service", Product: "inspirai"},
			{Name: "inspirai-api-specs", Repo: "git@github.com:inspirai-store/inspirai-api-specs.git", Type: "specs", Product: "inspirai"},
			{Name: "inspirai-devkit", Repo: "git@github.com:inspirai-store/inspirai-devkit.git", Type: "tools", Product: "inspirai"},
			{Name: "skill-market", Repo: "git@github.com:inspirai-store/skill-market.git", Type: "tools", Product: "independent"},
		},
	}
}

// GetProjectRoot 获取项目根目录（向上查找 .git）
func GetProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

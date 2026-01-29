package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
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
			// lingbo 产品线
			{Name: "lingbo-desktop", Repo: "git@github.com:inspirai-store/lingbo-desktop.git", Type: "client", Product: "lingbo"},
			{Name: "lingbo-web", Repo: "git@github.com:inspirai-store/lingbo-web.git", Type: "client", Product: "lingbo"},
			{Name: "lingbo-plugin", Repo: "git@github.com:inspirai-store/lingbo-plugin.git", Type: "tools", Product: "lingbo"},
			// inspirai 平台
			{Name: "inspirai-user", Repo: "git@github.com:inspirai-store/inspirai-user.git", Type: "service", Product: "inspirai"},
			{Name: "inspirai-ai-gateway", Repo: "git@github.com:inspirai-store/inspirai-ai-gateway.git", Type: "service", Product: "inspirai"},
			{Name: "inspirai-admin", Repo: "git@github.com:inspirai-store/inspirai-admin.git", Type: "client", Product: "inspirai"},
			{Name: "inspirai-web", Repo: "git@github.com:inspirai-store/inspirai-web.git", Type: "client", Product: "inspirai"},
			{Name: "inspirai-api-specs", Repo: "git@github.com:inspirai-store/inspirai-api-specs.git", Type: "specs", Product: "inspirai"},
			{Name: "inspirai-devkit", Repo: "git@github.com:inspirai-store/inspirai-devkit.git", Type: "tools", Product: "inspirai"},
			// magicbook 产品线
			{Name: "magicbook-service", Repo: "git@github.com:inspirai-store/magicbook-service.git", Type: "service", Product: "magicbook"},
			{Name: "magicbook-h5", Repo: "git@github.com:inspirai-store/magicbook-h5.git", Type: "client", Product: "magicbook"},
			{Name: "magicbook-admin", Repo: "git@github.com:inspirai-store/magicbook-admin.git", Type: "client", Product: "magicbook"},
			// zenix 产品线
			{Name: "zeni-x-desktop", Repo: "git@github.com:inspirai-store/zeni-x-desktop.git", Type: "tools", Product: "zenix"},
			// 独立项目
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

// GetGitCloneMethod 读取 .bootstrap.conf 获取 git clone 方式
func GetGitCloneMethod(root string) string {
	configPath := filepath.Join(root, ".bootstrap.conf")
	file, err := os.Open(configPath)
	if err != nil {
		return "ssh" // 默认 SSH
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 移除可能的 UTF-8 BOM
		line = strings.TrimPrefix(line, "\ufeff")
		if strings.HasPrefix(line, "GIT_CLONE_METHOD=") {
			return strings.TrimPrefix(line, "GIT_CLONE_METHOD=")
		}
	}
	return "ssh"
}

// ConvertRepoURL 根据配置转换 repo URL
func ConvertRepoURL(repo, method string) string {
	if method != "https" {
		return repo
	}
	// git@github.com:org/repo.git -> https://github.com/org/repo.git
	if strings.HasPrefix(repo, "git@github.com:") {
		return "https://github.com/" + strings.TrimPrefix(repo, "git@github.com:")
	}
	return repo
}

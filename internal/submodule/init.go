package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/inspirai-store/inspirai-devkit/internal/config"
)

// Init 初始化所有 submodule 并创建软链
func Init(cfg *config.Config, root string) error {
	submodulesDir := filepath.Join(root, cfg.SubmodulesDir)

	// 确保 .submodules 目录存在
	if err := os.MkdirAll(submodulesDir, 0755); err != nil {
		return fmt.Errorf("failed to create submodules dir: %w", err)
	}

	// 克隆每个 submodule
	for _, sm := range cfg.Submodules {
		smPath := filepath.Join(submodulesDir, sm.Name)
		if _, err := os.Stat(smPath); err == nil {
			color.Yellow("  [skip] %s already exists", sm.Name)
			continue
		}

		color.Cyan("  [clone] %s", sm.Name)
		cmd := exec.Command("git", "clone", sm.Repo, smPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			color.Red("  [error] failed to clone %s: %v", sm.Name, err)
			continue
		}
		color.Green("  [done] %s", sm.Name)
	}

	// 创建软链
	if err := CreateLinks(cfg, root); err != nil {
		return err
	}

	return nil
}

// CreateLinks 创建两种视图的软链
func CreateLinks(cfg *config.Config, root string) error {
	color.Cyan("\nCreating symlinks...")

	// by-type 视图
	typeMap := map[string][]string{
		"services": {},
		"clients":  {},
		"specs":    {},
		"tools":    {},
	}
	typeMapping := map[string]string{
		"service": "services",
		"client":  "clients",
		"specs":   "specs",
		"tools":   "tools",
	}

	for _, sm := range cfg.Submodules {
		if dir, ok := typeMapping[sm.Type]; ok {
			typeMap[dir] = append(typeMap[dir], sm.Name)
		}
	}

	for dir, names := range typeMap {
		dirPath := filepath.Join(root, "by-type", dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
		for _, name := range names {
			linkPath := filepath.Join(dirPath, name)
			target := filepath.Join("..", "..", cfg.SubmodulesDir, name)
			createSymlink(linkPath, target)
		}
	}

	// by-product 视图
	productMap := map[string]map[string]string{}
	for _, sm := range cfg.Submodules {
		if productMap[sm.Product] == nil {
			productMap[sm.Product] = map[string]string{}
		}
		// 使用短名称（去掉产品前缀）
		shortName := getShortName(sm.Name, sm.Product)
		productMap[sm.Product][shortName] = sm.Name
	}

	for product, items := range productMap {
		dirPath := filepath.Join(root, "by-product", product)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
		for shortName, fullName := range items {
			linkPath := filepath.Join(dirPath, shortName)
			target := filepath.Join("..", "..", cfg.SubmodulesDir, fullName)
			createSymlink(linkPath, target)
		}
	}

	color.Green("Symlinks created successfully")
	return nil
}

func createSymlink(linkPath, target string) {
	// 如果已存在，先删除
	if _, err := os.Lstat(linkPath); err == nil {
		os.Remove(linkPath)
	}

	if err := os.Symlink(target, linkPath); err != nil {
		color.Red("  [error] %s: %v", linkPath, err)
	} else {
		color.Green("  [link] %s -> %s", linkPath, target)
	}
}

func getShortName(name, product string) string {
	// lingbo-desktop -> desktop
	// inspirai-user -> user
	// skill-market -> skill-market (independent)
	prefixes := map[string]string{
		"lingbo":   "lingbo-",
		"inspirai": "inspirai-",
	}
	if prefix, ok := prefixes[product]; ok {
		if len(name) > len(prefix) && name[:len(prefix)] == prefix {
			return name[len(prefix):]
		}
	}
	return name
}

package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/inspirai-store/inspirai-devkit/internal/config"
)

// RunnerType 定义项目运行器类型
type RunnerType string

const (
	RunnerJust    RunnerType = "just"
	RunnerNpm     RunnerType = "npm"
	RunnerMake    RunnerType = "make"
	RunnerUnknown RunnerType = "unknown"
)

// Run 在指定项目中执行命令
func Run(cfg *config.Config, root string, projectName string, command string) error {
	submodulesDir := filepath.Join(root, cfg.SubmodulesDir)
	projectPath := filepath.Join(submodulesDir, projectName)

	// 检查项目是否存在
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project '%s' not found in %s", projectName, submodulesDir)
	}

	// 检测运行器类型
	runner := detectRunner(projectPath)
	if runner == RunnerUnknown {
		return fmt.Errorf("no supported build tool found in '%s' (justfile, package.json, or Makefile)", projectName)
	}

	// 构建命令
	var cmd *exec.Cmd
	switch runner {
	case RunnerJust:
		color.Cyan("  [just] %s %s", projectName, command)
		cmd = exec.Command("just", command)
	case RunnerNpm:
		color.Cyan("  [npm] %s run %s", projectName, command)
		cmd = exec.Command("npm", "run", command)
	case RunnerMake:
		color.Cyan("  [make] %s %s", projectName, command)
		cmd = exec.Command("make", command)
	}

	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// RunProduct 运行指定产品线的所有项目
func RunProduct(cfg *config.Config, root string, product string, command string) error {
	found := false
	for _, sm := range cfg.Submodules {
		if sm.Product == product {
			found = true
			color.Cyan("\n=== %s ===", sm.Name)
			if err := Run(cfg, root, sm.Name, command); err != nil {
				color.Red("  [error] %s: %v", sm.Name, err)
				// 继续执行其他项目，不中断
			}
		}
	}

	if !found {
		return fmt.Errorf("no projects found for product '%s'", product)
	}
	return nil
}

// ListRunnable 列出所有可运行的项目及其运行器类型
func ListRunnable(cfg *config.Config, root string) {
	submodulesDir := filepath.Join(root, cfg.SubmodulesDir)

	fmt.Printf("%-20s %-10s %-10s\n", "PROJECT", "PRODUCT", "RUNNER")
	fmt.Println("----------------------------------------")

	for _, sm := range cfg.Submodules {
		projectPath := filepath.Join(submodulesDir, sm.Name)
		runner := detectRunner(projectPath)
		runnerStr := string(runner)
		if runner == RunnerUnknown {
			runnerStr = "-"
		}
		fmt.Printf("%-20s %-10s %-10s\n", sm.Name, sm.Product, runnerStr)
	}
}

func detectRunner(projectPath string) RunnerType {
	// 优先级：justfile > package.json > Makefile
	if _, err := os.Stat(filepath.Join(projectPath, "justfile")); err == nil {
		return RunnerJust
	}
	if _, err := os.Stat(filepath.Join(projectPath, "package.json")); err == nil {
		return RunnerNpm
	}
	if _, err := os.Stat(filepath.Join(projectPath, "Makefile")); err == nil {
		return RunnerMake
	}
	return RunnerUnknown
}

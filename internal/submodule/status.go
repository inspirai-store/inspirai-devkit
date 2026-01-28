package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/inspirai-store/inspirai-devkit/internal/config"
)

// Status 显示所有 submodule 的状态
func Status(cfg *config.Config, root string) error {
	submodulesDir := filepath.Join(root, cfg.SubmodulesDir)

	fmt.Printf("%-20s %-15s %-10s %s\n", "NAME", "BRANCH", "STATUS", "COMMIT")
	fmt.Println(strings.Repeat("-", 70))

	for _, sm := range cfg.Submodules {
		smPath := filepath.Join(submodulesDir, sm.Name)

		if _, err := os.Stat(smPath); os.IsNotExist(err) {
			color.Red("%-20s %-15s %-10s %s\n", sm.Name, "-", "missing", "-")
			continue
		}

		branch := getGitBranch(smPath)
		status := getGitStatus(smPath)
		commit := getGitCommit(smPath)

		statusColor := color.New(color.FgGreen)
		if status != "clean" {
			statusColor = color.New(color.FgYellow)
		}

		fmt.Printf("%-20s %-15s ", sm.Name, branch)
		statusColor.Printf("%-10s ", status)
		fmt.Printf("%s\n", commit)
	}

	return nil
}

func getGitBranch(path string) string {
	cmd := exec.Command("git", "-C", path, "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func getGitStatus(path string) string {
	cmd := exec.Command("git", "-C", path, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return "error"
	}
	if len(strings.TrimSpace(string(out))) > 0 {
		return "modified"
	}
	return "clean"
}

func getGitCommit(path string) string {
	cmd := exec.Command("git", "-C", path, "log", "-1", "--format=%h %s")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	result := strings.TrimSpace(string(out))
	if len(result) > 50 {
		result = result[:47] + "..."
	}
	return result
}

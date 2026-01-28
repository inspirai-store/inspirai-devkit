package submodule

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/inspirai-store/inspirai-devkit/internal/config"
)

// Sync 同步所有 submodule
func Sync(cfg *config.Config, root string) error {
	submodulesDir := filepath.Join(root, cfg.SubmodulesDir)

	for _, sm := range cfg.Submodules {
		smPath := filepath.Join(submodulesDir, sm.Name)

		if _, err := os.Stat(smPath); os.IsNotExist(err) {
			color.Yellow("  [skip] %s not found", sm.Name)
			continue
		}

		color.Cyan("  [sync] %s", sm.Name)

		cmd := exec.Command("git", "-C", smPath, "pull", "--rebase")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			color.Red("  [error] %s: %v", sm.Name, err)
			continue
		}

		color.Green("  [done] %s", sm.Name)
	}

	return nil
}

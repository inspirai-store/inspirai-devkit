package main

import (
	"fmt"
	"os"

	"github.com/inspirai-store/inspirai-devkit/internal/config"
	"github.com/inspirai-store/inspirai-devkit/internal/submodule"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:     "sm",
		Short:   "Submodule Manager for inspirai projects",
		Version: version,
	}

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(linksCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize all submodules and create symlinks",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.GetProjectRoot()
			if err != nil {
				return fmt.Errorf("not in a git repository: %w", err)
			}

			cfg := config.DefaultConfig()
			fmt.Println("Initializing submodules...")
			return submodule.Init(cfg, root)
		},
	}
}

func syncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync all submodules (git pull --rebase)",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.GetProjectRoot()
			if err != nil {
				return fmt.Errorf("not in a git repository: %w", err)
			}

			cfg := config.DefaultConfig()
			fmt.Println("Syncing submodules...")
			return submodule.Sync(cfg, root)
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show status of all submodules",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.GetProjectRoot()
			if err != nil {
				return fmt.Errorf("not in a git repository: %w", err)
			}

			cfg := config.DefaultConfig()
			return submodule.Status(cfg, root)
		},
	}
}

func linksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "links",
		Short: "Rebuild all symlinks",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("sm links: Not implemented yet")
			return nil
		},
	}
}

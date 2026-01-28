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
	rootCmd.AddCommand(runCmd())

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
			root, err := config.GetProjectRoot()
			if err != nil {
				return fmt.Errorf("not in a git repository: %w", err)
			}

			cfg := config.DefaultConfig()
			return submodule.CreateLinks(cfg, root)
		},
	}
}

func runCmd() *cobra.Command {
	var listFlag bool
	var productFlag string

	cmd := &cobra.Command{
		Use:   "run <project> <command>",
		Short: "Run a command in a project (auto-detects just/npm/make)",
		Long: `Run a command in a project directory.

Automatically detects the build tool:
  - justfile  -> just <command>
  - package.json -> npm run <command>
  - Makefile  -> make <command>

Examples:
  sm run lingbo-desktop dev     # Run 'just dev' in lingbo-desktop
  sm run lingbo-web dev         # Run 'npm run dev' in lingbo-web
  sm run --product lingbo dev   # Run 'dev' in all lingbo projects
  sm run --list                 # List all projects and their runners`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.GetProjectRoot()
			if err != nil {
				return fmt.Errorf("not in a git repository: %w", err)
			}

			cfg := config.DefaultConfig()

			// List mode
			if listFlag {
				submodule.ListRunnable(cfg, root)
				return nil
			}

			// Product mode
			if productFlag != "" {
				if len(args) < 1 {
					return fmt.Errorf("command required: sm run --product <product> <command>")
				}
				return submodule.RunProduct(cfg, root, productFlag, args[0])
			}

			// Project mode
			if len(args) < 2 {
				return fmt.Errorf("usage: sm run <project> <command>")
			}

			return submodule.Run(cfg, root, args[0], args[1])
		},
	}

	cmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all projects and their runners")
	cmd.Flags().StringVarP(&productFlag, "product", "p", "", "Run command in all projects of a product")

	return cmd
}

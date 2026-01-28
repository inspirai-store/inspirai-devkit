package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/inspirai-store/inspirai-devkit/internal/codegen"
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
	rootCmd.AddCommand(codegenCmd())

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

func codegenCmd() *cobra.Command {
	var langFlag string
	var outputFlag string
	var listFlag bool

	cmd := &cobra.Command{
		Use:   "codegen <service>",
		Short: "Generate code from API specifications",
		Long: `Generate client code from API specifications.

Reads YAML specs from .submodules/inspirai-api-specs/<service>/ and generates
type definitions and API client code.

Supported languages:
  - go         Generate Go structs
  - typescript Generate TypeScript interfaces

Examples:
  sm codegen --list                              # List available services
  sm codegen inspirai-user --lang go -o ./gen    # Generate Go code
  sm codegen inspirai-user --lang ts -o ./types  # Generate TypeScript`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.GetProjectRoot()
			if err != nil {
				return fmt.Errorf("not in a git repository: %w", err)
			}

			specDir := filepath.Join(root, ".submodules", "inspirai-api-specs")

			// List mode
			if listFlag {
				services, err := codegen.ListServices(specDir)
				if err != nil {
					return err
				}
				fmt.Println("Available services:")
				for _, svc := range services {
					fmt.Printf("  - %s\n", svc)
				}
				return nil
			}

			if len(args) < 1 {
				return fmt.Errorf("service name required: sm codegen <service>")
			}

			if langFlag == "" {
				return fmt.Errorf("language required: use --lang go or --lang typescript")
			}

			if outputFlag == "" {
				outputFlag = filepath.Join(".", "generated", args[0])
			}

			gen := &codegen.Generator{
				SpecDir: specDir,
				Lang:    langFlag,
				Output:  outputFlag,
			}

			fmt.Printf("Generating %s code for %s...\n", langFlag, args[0])
			return gen.Generate(args[0])
		},
	}

	cmd.Flags().StringVarP(&langFlag, "lang", "l", "", "Target language (go, typescript)")
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output directory")
	cmd.Flags().BoolVar(&listFlag, "list", false, "List available services")

	return cmd
}

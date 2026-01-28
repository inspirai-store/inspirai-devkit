package main

import (
	"fmt"
	"os"

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
			fmt.Println("sm init: Not implemented yet")
			return nil
		},
	}
}

func syncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync all submodules (git pull)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("sm sync: Not implemented yet")
			return nil
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show status of all submodules",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("sm status: Not implemented yet")
			return nil
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

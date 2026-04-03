package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmashuda/wormtongue/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize wormtongue config and library directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		cfgPath := cfgFile
		if cfgPath == "" {
			cfgPath = config.DefaultPath()
		}

		if !force {
			if _, err := os.Stat(cfgPath); err == nil {
				return fmt.Errorf("config already exists at %s (use --force to overwrite)", cfgPath)
			}
		}

		examplesDir := filepath.Join(filepath.Dir(cfgPath), "examples")

		if err := os.MkdirAll(examplesDir, 0o755); err != nil {
			return fmt.Errorf("creating examples directory: %w", err)
		}

		cfg := &config.Config{
			Sources: []config.Source{
				{Name: "default", Path: examplesDir},
			},
		}
		if err := config.Save(cfgPath, cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		cmd.Printf("Created config: %s\n", cfgPath)
		cmd.Printf("Created library: %s\n", examplesDir)
		return nil
	},
}

func init() {
	initCmd.Flags().Bool("force", false, "Overwrite existing config")
	rootCmd.AddCommand(initCmd)
}

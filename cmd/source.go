package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/dmashuda/wormtongue/internal/config"
	"github.com/spf13/cobra"
)

var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Manage example sources",
}

var sourceAddCmd = &cobra.Command{
	Use:   "add <name> <path>",
	Short: "Register an external example source",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path, err := filepath.Abs(args[1])
		if err != nil {
			return fmt.Errorf("resolving path: %w", err)
		}

		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}

		for _, s := range cfg.Sources {
			if s.Name == name {
				return fmt.Errorf("source %q already exists", name)
			}
		}

		cfg.Sources = append(cfg.Sources, config.Source{Name: name, Path: path})
		if err := config.Save(cfgFile, cfg); err != nil {
			return err
		}

		cmd.Printf("Added source %q at %s\n", name, path)
		return nil
	},
}

var sourceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered example sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}

		if len(cfg.Sources) == 0 {
			cmd.Println("No external sources registered.")
			return nil
		}

		cmd.Printf("%-20s %s\n", "NAME", "PATH")
		for _, s := range cfg.Sources {
			cmd.Printf("%-20s %s\n", s.Name, s.Path)
		}
		return nil
	},
}

var sourceRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a registered example source",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}

		found := false
		filtered := make([]config.Source, 0, len(cfg.Sources))
		for _, s := range cfg.Sources {
			if s.Name == name {
				found = true
				continue
			}
			filtered = append(filtered, s)
		}

		if !found {
			return fmt.Errorf("source %q not found", name)
		}

		cfg.Sources = filtered
		if err := config.Save(cfgFile, cfg); err != nil {
			return err
		}

		cmd.Printf("Removed source %q\n", name)
		return nil
	},
}

func init() {
	sourceCmd.AddCommand(sourceAddCmd, sourceListCmd, sourceRemoveCmd)
	rootCmd.AddCommand(sourceCmd)
}

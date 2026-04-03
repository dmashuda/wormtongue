package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmashuda/wormtongue/internal/config"
	"github.com/dmashuda/wormtongue/internal/examples"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	store   *examples.ExampleStore
)

var rootCmd = &cobra.Command{
	Use:   "wormtongue",
	Short: "A code style example library for LLMs",
	Long:  "Wormtongue provides curated code style and pattern examples that LLMs can query via MCP or CLI.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Name() == "init" {
			return nil
		}
		return initStore()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default ~/.wormtongue/config.yaml)")
}

func initStore() error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	sources := builtinSources()
	for _, s := range cfg.Sources {
		sources = append(sources, s.Path)
	}

	store = examples.NewStore(sources)
	return nil
}

func builtinSources() []string {
	// Check WORMTONGUE_EXAMPLES env var first (useful for development)
	if envDir := os.Getenv("WORMTONGUE_EXAMPLES"); envDir != "" {
		abs, err := filepath.Abs(envDir)
		if err == nil {
			return []string{abs}
		}
		return []string{envDir}
	}

	// Look for examples/ relative to the binary
	exe, err := os.Executable()
	if err != nil {
		return nil
	}
	dir := filepath.Join(filepath.Dir(exe), "examples")
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return []string{dir}
	}
	return nil
}

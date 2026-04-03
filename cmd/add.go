package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/dmashuda/wormtongue/internal/examples"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <language> <category> <name>",
	Short: "Add a new example",
	Long:  "Add a new code example. Content is read from stdin or provided via the --content flag.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		language := args[0]
		category := args[1]
		name := args[2]

		content, _ := cmd.Flags().GetString("content")
		force, _ := cmd.Flags().GetBool("force")

		if content == "" {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
			content = string(data)
		}

		ex, err := store.Add(language, category, name, content, examples.AddOptions{
			Force: force,
		})
		if err != nil {
			return err
		}

		cmd.Printf("Added example: %s\n", ex.Path)
		return nil
	},
}

func init() {
	addCmd.Flags().StringP("content", "m", "", "Example content (reads from stdin if not provided)")
	addCmd.Flags().BoolP("force", "f", false, "Overwrite existing example")
	rootCmd.AddCommand(addCmd)
}

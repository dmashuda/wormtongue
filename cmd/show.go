package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <path>",
	Short: "Display a specific example",
	Long:  "Display the content of an example by its path (e.g. go/concurrency/worker-pool).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, content, err := store.Get(args[0])
		if err != nil {
			return err
		}
		fmt.Print(content)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

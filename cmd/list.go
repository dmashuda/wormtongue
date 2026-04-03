package cmd

import (
	"github.com/dmashuda/wormtongue/internal/examples"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available examples",
	RunE: func(cmd *cobra.Command, args []string) error {
		language, _ := cmd.Flags().GetString("language")
		category, _ := cmd.Flags().GetString("category")

		filter := examples.Filter{
			Language: language,
			Category: category,
		}

		results := store.List(filter)
		if len(results) == 0 {
			cmd.Println("No examples found.")
			return nil
		}

		cmd.Printf("%-12s %-20s %s\n", "LANGUAGE", "CATEGORY", "NAME")
		for _, ex := range results {
			cmd.Printf("%-12s %-20s %s\n", ex.Language, ex.Category, ex.Name)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().StringP("language", "l", "", "Filter by language")
	listCmd.Flags().StringP("category", "c", "", "Filter by category")
	rootCmd.AddCommand(listCmd)
}

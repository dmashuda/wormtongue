package cmd

import (
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search examples by keyword",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")

		results := store.Search(args[0], limit)
		if len(results) == 0 {
			cmd.Println("No matching examples found.")
			return nil
		}

		for _, r := range results {
			cmd.Printf("%s/%s/%s\n", r.Example.Language, r.Example.Category, r.Example.Name)
			if r.MatchLine != "" {
				cmd.Printf("  ...%s...\n", r.MatchLine)
			}
			cmd.Println()
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().IntP("limit", "n", 10, "Maximum number of results")
	rootCmd.AddCommand(searchCmd)
}

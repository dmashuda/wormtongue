package cmd

import (
	"github.com/spf13/cobra"
)

var languagesCmd = &cobra.Command{
	Use:   "languages",
	Short: "List available languages",
	RunE: func(cmd *cobra.Command, args []string) error {
		langs := store.Languages()
		if len(langs) == 0 {
			cmd.Println("No languages found.")
			return nil
		}

		for _, lang := range langs {
			cmd.Println(lang)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(languagesCmd)
}

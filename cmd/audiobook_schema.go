package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/deegital/elevencli/internal/audiobook"
)

var audiobookSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Print the JSON Schema for audiobook script files",
	Annotations: map[string]string{"noAuth": "true"},
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(audiobook.Schema)
		return nil
	},
}

func init() {
	audiobookCmd.AddCommand(audiobookSchemaCmd)
}

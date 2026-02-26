package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var voiceSearch string

var voicesCmd = &cobra.Command{
	Use:   "voices",
	Short: "List available voices",
	RunE: func(cmd *cobra.Command, args []string) error {
		voices, err := client.GetVoices()
		if err != nil {
			return fmt.Errorf("failed to list voices: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tCATEGORY\tLABELS")

		for _, v := range voices {
			if voiceSearch != "" && !strings.Contains(strings.ToLower(v.Name), strings.ToLower(voiceSearch)) {
				continue
			}
			labels := formatLabels(v.Labels)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", v.VoiceId, v.Name, v.Category, labels)
		}
		return w.Flush()
	},
}

func formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	parts := make([]string, 0, len(labels))
	for k, v := range labels {
		parts = append(parts, k+":"+v)
	}
	return strings.Join(parts, ", ")
}

func init() {
	voicesCmd.Flags().StringVarP(&voiceSearch, "search", "s", "", "Filter voices by name")
	rootCmd.AddCommand(voicesCmd)
}

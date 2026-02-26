package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	sfxOutput   string
	sfxDuration float64
	sfxFormat   string
)

type soundGenRequest struct {
	Text            string  `json:"text"`
	DurationSeconds float64 `json:"duration_seconds,omitempty"`
}

var sfxCmd = &cobra.Command{
	Use:   "sfx <prompt>",
	Short: "Generate a sound effect from a text description",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiFormat, err := resolveFormat(sfxFormat)
		if err != nil {
			return err
		}

		key, err := resolveAPIKeyValue()
		if err != nil {
			return err
		}

		req := soundGenRequest{Text: args[0]}
		if sfxDuration > 0 {
			req.DurationSeconds = sfxDuration
		}

		body, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("failed to build request: %w", err)
		}

		url := fmt.Sprintf("https://api.elevenlabs.io/v1/sound-generation?output_format=%s", apiFormat)
		httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("xi-api-key", key)

		fmt.Fprintf(os.Stderr, "Generating sound effect...\n")

		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			return fmt.Errorf("SFX request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
		}

		audio, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if err := os.WriteFile(sfxOutput, audio, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", sfxOutput, err)
		}

		fmt.Println(sfxOutput)
		return nil
	},
}

func init() {
	sfxCmd.Flags().StringVarP(&sfxOutput, "output", "o", "output.mp3", "Output file path")
	sfxCmd.Flags().Float64VarP(&sfxDuration, "duration", "d", 0, "Duration in seconds (0.5-30)")
	sfxCmd.Flags().StringVarP(&sfxFormat, "format", "f", "mp3", "Output format: mp3, pcm, ulaw")
	rootCmd.AddCommand(sfxCmd)
}

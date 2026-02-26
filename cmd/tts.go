package cmd

import (
	"fmt"
	"os"

	elevenlabs "github.com/haguro/elevenlabs-go"
	"github.com/spf13/cobra"
)

var (
	ttsVoice  string
	ttsOutput string
	ttsFormat string
	ttsModel  string
)

// formatMap maps user-friendly format names to ElevenLabs API format strings.
var formatMap = map[string]string{
	"mp3":  "mp3_44100_128",
	"pcm":  "pcm_44100",
	"ulaw": "ulaw_8000",
}

func resolveFormat(f string) (string, error) {
	if mapped, ok := formatMap[f]; ok {
		return mapped, nil
	}
	return "", fmt.Errorf("unsupported format %q (supported: mp3, pcm, ulaw)", f)
}

var ttsCmd = &cobra.Command{
	Use:   "tts <text>",
	Short: "Generate speech from text",
	Long:  `Generate speech from text using ElevenLabs text-to-speech API.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ttsVoice == "" {
			return fmt.Errorf("--voice is required. Use 'elevencli voices' to list available voices")
		}

		apiFormat, err := resolveFormat(ttsFormat)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Generating speech...\n")

		audio, err := client.TextToSpeech(ttsVoice, elevenlabs.TextToSpeechRequest{
			Text:    args[0],
			ModelID: ttsModel,
		}, elevenlabs.OutputFormat(apiFormat))
		if err != nil {
			return fmt.Errorf("TTS request failed: %w", err)
		}

		if err := os.WriteFile(ttsOutput, audio, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", ttsOutput, err)
		}

		fmt.Println(ttsOutput)
		return nil
	},
}

func init() {
	ttsCmd.Flags().StringVarP(&ttsVoice, "voice", "v", "", "Voice ID (required)")
	ttsCmd.Flags().StringVarP(&ttsOutput, "output", "o", "output.mp3", "Output file path")
	ttsCmd.Flags().StringVarP(&ttsFormat, "format", "f", "mp3", "Output format: mp3, pcm, ulaw")
	ttsCmd.Flags().StringVarP(&ttsModel, "model", "m", "eleven_multilingual_v2", "Model ID")
	_ = ttsCmd.MarkFlagRequired("voice")
	rootCmd.AddCommand(ttsCmd)
}

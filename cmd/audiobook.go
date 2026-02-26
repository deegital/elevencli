package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/deegital/elevencli/internal/audio"
	"github.com/deegital/elevencli/internal/audiobook"
)

var (
	audiobookOutput     string
	audiobookKeepBlocks bool
	audiobookStdin      bool
	audiobookStdout     bool
)

var audiobookCmd = &cobra.Command{
	Use:   "audiobook [script.json]",
	Short: "Generate an audiobook from a JSON script",
	Long: `Generate an audiobook by processing a JSON script that defines a sequence
of TTS narration, sound effects, and silence blocks. The blocks are rendered
via the ElevenLabs API and merged into a single MP3 file.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateStdinArgs(cmd, args, audiobookStdin, audiobookStdout); err != nil {
			return err
		}

		var data []byte
		var err error
		if audiobookStdin {
			data, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
		} else {
			data, err = os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read script: %w", err)
			}
		}

		var script audiobook.Script
		if err := json.Unmarshal(data, &script); err != nil {
			return fmt.Errorf("failed to parse script: %w", err)
		}

		if err := script.Validate(); err != nil {
			return fmt.Errorf("invalid script: %w", err)
		}

		key, err := resolveAPIKeyValue()
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Generating audiobook (%d blocks)...\n", len(script.Blocks))

		result, err := audiobook.Generate(&script, key, client)
		if err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}

		mp3Data, err := audio.EncodePCMToMP3(result.MergedPCM)
		if err != nil {
			return fmt.Errorf("MP3 encoding failed: %w", err)
		}

		if audiobookKeepBlocks {
			dir := "."
			if !audiobookStdout {
				dir = filepath.Dir(audiobookOutput)
			}
			for i, pcm := range result.BlockPCMs {
				blockMP3, err := audio.EncodePCMToMP3(pcm)
				if err != nil {
					return fmt.Errorf("failed to encode block %d: %w", i, err)
				}
				blockPath := filepath.Join(dir, fmt.Sprintf("block_%03d.mp3", i+1))
				if err := os.WriteFile(blockPath, blockMP3, 0644); err != nil {
					return fmt.Errorf("failed to write %s: %w", blockPath, err)
				}
				fmt.Fprintf(os.Stderr, "Wrote %s\n", blockPath)
			}
		}

		return writeOutput(mp3Data, audiobookOutput, audiobookStdout)
	},
}

func init() {
	audiobookCmd.Flags().StringVarP(&audiobookOutput, "output", "o", "audiobook.mp3", "Output file path")
	audiobookCmd.Flags().BoolVar(&audiobookKeepBlocks, "keep-blocks", false, "Keep individual block audio files")
	audiobookCmd.Flags().BoolVar(&audiobookStdin, "stdin", false, "Read script JSON from stdin")
	audiobookCmd.Flags().BoolVar(&audiobookStdout, "stdout", false, "Write audio to stdout")
	rootCmd.AddCommand(audiobookCmd)
}

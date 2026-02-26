package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func validateStdinArgs(cmd *cobra.Command, args []string, useStdin, useStdout bool) error {
	if useStdin && len(args) > 0 {
		return fmt.Errorf("cannot use --stdin with a positional argument")
	}
	if !useStdin && len(args) == 0 {
		return fmt.Errorf("requires a positional argument or --stdin")
	}
	if useStdout && cmd.Flags().Changed("output") {
		return fmt.Errorf("cannot use --stdout with --output")
	}
	return nil
}

func readTextFromStdinOrArg(useStdin bool, args []string) (string, error) {
	if useStdin {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read stdin: %w", err)
		}
		text := strings.TrimSpace(string(data))
		if text == "" {
			return "", fmt.Errorf("stdin was empty")
		}
		return text, nil
	}
	return args[0], nil
}

func writeOutput(audioData []byte, outputPath string, useStdout bool) error {
	if useStdout {
		_, err := os.Stdout.Write(audioData)
		return err
	}
	if err := os.WriteFile(outputPath, audioData, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}
	fmt.Println(outputPath)
	return nil
}

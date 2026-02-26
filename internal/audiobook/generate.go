package audiobook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	elevenlabs "github.com/haguro/elevenlabs-go"

	"github.com/deegital/elevencli/internal/audio"
)

// GenerateResult holds the output of audiobook generation.
type GenerateResult struct {
	// MergedPCM is the final concatenated/mixed PCM audio.
	MergedPCM []byte
	// BlockPCMs holds individual block PCM data (indexed by block position).
	BlockPCMs [][]byte
}

type sfxRequest struct {
	Text            string  `json:"text"`
	DurationSeconds float64 `json:"duration_seconds,omitempty"`
}

// Generate processes an audiobook script and returns PCM audio data.
func Generate(script *Script, apiKey string, client *elevenlabs.Client) (*GenerateResult, error) {
	var (
		segments  [][]byte // sequential PCM segments to concatenate
		blockPCMs [][]byte // one per block for --keep-blocks
		pendingBG []byte   // background SFX PCM waiting to be mixed into next TTS
	)

	for i, block := range script.Blocks {
		fmt.Fprintf(os.Stderr, "Processing block %d/%d (%s)...\n", i+1, len(script.Blocks), block.Type)

		switch block.Type {
		case "tts":
			pcm, err := generateTTS(block, client)
			if err != nil {
				return nil, fmt.Errorf("block %d (tts): %w", i, err)
			}

			if pendingBG != nil {
				pcm = audio.Mix(pcm, pendingBG)
				pendingBG = nil
			}

			segments = append(segments, pcm)
			blockPCMs = append(blockPCMs, pcm)

		case "sfx":
			pcm, err := generateSFX(block, apiKey)
			if err != nil {
				return nil, fmt.Errorf("block %d (sfx): %w", i, err)
			}

			if block.Background {
				pendingBG = pcm
				blockPCMs = append(blockPCMs, pcm)
			} else {
				segments = append(segments, pcm)
				blockPCMs = append(blockPCMs, pcm)
			}

		case "silence":
			pcm := audio.Silence(block.Duration)
			segments = append(segments, pcm)
			blockPCMs = append(blockPCMs, pcm)
		}
	}

	// If there's a trailing background SFX with no following TTS, append it.
	if pendingBG != nil {
		segments = append(segments, pendingBG)
	}

	merged := audio.Concat(segments...)

	return &GenerateResult{
		MergedPCM: merged,
		BlockPCMs: blockPCMs,
	}, nil
}

func generateTTS(block Block, client *elevenlabs.Client) ([]byte, error) {
	req := elevenlabs.TextToSpeechRequest{
		Text:    block.Text,
		ModelID: block.Model,
	}

	if block.Model == "" {
		req.ModelID = "eleven_multilingual_v2"
	}

	if block.Stability != 0 || block.SimilarityBoost != 0 || block.Style != 0 {
		req.VoiceSettings = &elevenlabs.VoiceSettings{
			Stability:       block.Stability,
			SimilarityBoost: block.SimilarityBoost,
			Style:           block.Style,
		}
	}

	pcm, err := client.TextToSpeech(block.Voice, req, elevenlabs.OutputFormat("pcm_44100"))
	if err != nil {
		return nil, fmt.Errorf("TTS API request failed: %w", err)
	}

	return pcm, nil
}

func generateSFX(block Block, apiKey string) ([]byte, error) {
	req := sfxRequest{Text: block.Text}
	if block.Duration > 0 {
		req.DurationSeconds = block.Duration
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build SFX request: %w", err)
	}

	url := "https://api.elevenlabs.io/v1/sound-generation?output_format=pcm_44100"
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create SFX request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("xi-api-key", apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("SFX API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("SFX API error (%d): %s", resp.StatusCode, string(respBody))
	}

	pcm, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read SFX response: %w", err)
	}

	return pcm, nil
}

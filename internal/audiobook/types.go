package audiobook

import "fmt"

// Script represents an audiobook script containing a sequence of blocks.
type Script struct {
	Blocks []Block `json:"blocks"`
}

// Block represents a single segment in the audiobook script.
type Block struct {
	Type            string  `json:"type"`
	Voice           string  `json:"voice,omitempty"`
	Text            string  `json:"text,omitempty"`
	Model           string  `json:"model,omitempty"`
	Stability       float32 `json:"stability,omitempty"`
	SimilarityBoost float32 `json:"similarity_boost,omitempty"`
	Style           float32 `json:"style,omitempty"`
	Speed           float64 `json:"speed,omitempty"`
	Background      bool    `json:"background,omitempty"`
	Duration        float64 `json:"duration,omitempty"`
}

// Validate checks the script for structural correctness.
func (s *Script) Validate() error {
	if len(s.Blocks) == 0 {
		return fmt.Errorf("script has no blocks")
	}
	for i, b := range s.Blocks {
		if err := b.validate(); err != nil {
			return fmt.Errorf("block %d: %w", i, err)
		}
	}
	return nil
}

func (b *Block) validate() error {
	switch b.Type {
	case "tts":
		if b.Voice == "" {
			return fmt.Errorf("tts block requires 'voice'")
		}
		if b.Text == "" {
			return fmt.Errorf("tts block requires 'text'")
		}
		if b.Stability < 0 || b.Stability > 1 {
			return fmt.Errorf("tts 'stability' must be between 0.0 and 1.0")
		}
		if b.SimilarityBoost < 0 || b.SimilarityBoost > 1 {
			return fmt.Errorf("tts 'similarity_boost' must be between 0.0 and 1.0")
		}
		if b.Style < 0 || b.Style > 1 {
			return fmt.Errorf("tts 'style' must be between 0.0 and 1.0")
		}
		if b.Speed != 0 && (b.Speed < 0.5 || b.Speed > 2.0) {
			return fmt.Errorf("tts 'speed' must be between 0.5 and 2.0")
		}
	case "sfx":
		if b.Text == "" {
			return fmt.Errorf("sfx block requires 'text'")
		}
		if b.Duration != 0 && (b.Duration < 0.5 || b.Duration > 22) {
			return fmt.Errorf("sfx 'duration' must be between 0.5 and 22.0 seconds")
		}
	case "silence":
		if b.Duration <= 0 {
			return fmt.Errorf("silence block requires positive 'duration'")
		}
	default:
		return fmt.Errorf("unknown block type %q", b.Type)
	}
	return nil
}

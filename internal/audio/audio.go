package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/braheezy/shine-mp3/pkg/mp3"
)

const (
	SampleRate = 44100
	Channels   = 1
	BitDepth   = 16
)

// Concat appends PCM byte slices sequentially.
func Concat(segments ...[]byte) []byte {
	total := 0
	for _, s := range segments {
		total += len(s)
	}
	out := make([]byte, 0, total)
	for _, s := range segments {
		out = append(out, s...)
	}
	return out
}

// Mix overlays two PCM byte slices by adding int16 samples with clamping.
// The shorter slice is zero-padded to match the longer one.
func Mix(base, overlay []byte) []byte {
	baseLen := len(base) / 2
	overlayLen := len(overlay) / 2
	outLen := baseLen
	if overlayLen > outLen {
		outLen = overlayLen
	}
	out := make([]byte, outLen*2)

	for i := 0; i < outLen; i++ {
		var a, b int32
		if i < baseLen {
			a = int32(int16(binary.LittleEndian.Uint16(base[i*2:])))
		}
		if i < overlayLen {
			b = int32(int16(binary.LittleEndian.Uint16(overlay[i*2:])))
		}
		mixed := a + b
		if mixed > 32767 {
			mixed = 32767
		}
		if mixed < -32768 {
			mixed = -32768
		}
		binary.LittleEndian.PutUint16(out[i*2:], uint16(int16(mixed)))
	}
	return out
}

// Silence generates zero-filled PCM bytes for the given duration.
func Silence(duration float64) []byte {
	numSamples := int(duration * float64(SampleRate))
	return make([]byte, numSamples*2)
}

// EncodePCMToMP3 encodes 16-bit mono PCM data to MP3.
func EncodePCMToMP3(pcm []byte) ([]byte, error) {
	numSamples := len(pcm) / 2
	samples := make([]int16, numSamples)
	for i := 0; i < numSamples; i++ {
		samples[i] = int16(binary.LittleEndian.Uint16(pcm[i*2:]))
	}

	encoder := mp3.NewEncoder(Channels, SampleRate)

	var buf bytes.Buffer
	if err := encoder.Write(&buf, samples); err != nil {
		return nil, fmt.Errorf("MP3 encoding failed: %w", err)
	}

	return buf.Bytes(), nil
}

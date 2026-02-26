```
█▀▀ █   █▀▀ █  █ █▀▀ █▄  █ █▀▀ █   ▀█▀
█▀▀ █   █▀▀ ▀▄▄▀ █▀▀ █ █ █ █   █    █
▀▀▀ ▀▀▀ ▀▀▀  ▀▀  ▀▀▀ ▀  ▀▀ ▀▀▀ ▀▀▀ ▀▀▀
```

ElevenLabs CLI by — text-to-speech and sound effects from the command line.

## Features

- **tts** — Convert text to speech with any ElevenLabs voice
- **sfx** — Generate sound effects from a text prompt
- **voices** — List and search available voices
- **audiobook** — Stitch narration, sound effects, and silence into a single audio file

## Installation

```sh
go install github.com/deegital/elevencli@latest
```

Or build from source:

```sh
git clone https://github.com/deegital/elevencli.git
cd elevencli
make build
```

## Configuration

An [ElevenLabs API key](https://elevenlabs.io) is required. Provide it in one of three ways (highest priority first):

| Method | Example |
|--------|---------|
| Flag | `--api-key sk-...` |
| Environment variable | `export ELEVENLABS_API_KEY=sk-...` |
| Config file | `~/.elevencli.yaml` with `api_key: sk-...` |

## Usage

### Text-to-Speech

```sh
elevencli tts "Hello, world!" --voice JBFqnCBsd6RMkjVDRZzb --output hello.mp3
```

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-v, --voice` | *(required)* | Voice ID |
| `-o, --output` | `output.mp3` | Output file path |
| `-f, --format` | `mp3` | Audio format: `mp3`, `pcm`, `ulaw` |
| `-m, --model` | `eleven_multilingual_v2` | Model ID |

### Sound Effects

```sh
elevencli sfx "heavy rain on a tin roof" --duration 5 --output rain.mp3
```

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output` | `output.mp3` | Output file path |
| `-d, --duration` | auto | Duration in seconds (0.5–30) |
| `-f, --format` | `mp3` | Audio format: `mp3`, `pcm`, `ulaw` |

### List Voices

```sh
elevencli voices
elevencli voices --search "aria"
```

Flags:

| Flag | Description |
|------|-------------|
| `-s, --search` | Filter voices by name (case-insensitive) |

### Audiobook

Generate a complete audiobook from a JSON script that combines narration, sound effects, and silence:

```sh
elevencli audiobook examples/story.json --output story.mp3
```

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output` | `audiobook.mp3` | Output MP3 file path |
| `--keep-blocks` | `false` | Save individual block audio files |

#### Script Format

The script is a JSON file with an array of blocks. Each block has a `type` — one of `tts`, `sfx`, or `silence`:

```json
{
  "blocks": [
    {
      "type": "tts",
      "voice": "JBFqnCBsd6RMkjVDRZzb",
      "text": "Once upon a time...",
      "stability": 0.5,
      "similarity_boost": 0.75
    },
    {
      "type": "sfx",
      "text": "gentle ocean waves",
      "background": true,
      "duration": 8.0
    },
    {
      "type": "silence",
      "duration": 1.5
    }
  ]
}
```

Setting `"background": true` on an SFX block mixes it with the next TTS block instead of playing sequentially.

Print the full JSON Schema for the script format:

```sh
elevencli audiobook schema
```

See [`examples/story.json`](examples/story.json) for a complete example.

## License

[MIT](LICENSE)

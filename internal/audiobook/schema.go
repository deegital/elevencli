package audiobook

// Schema is the JSON Schema for an audiobook script file.
const Schema = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Audiobook Script",
  "description": "A sequence of narration, sound effects, and silence blocks rendered into a single audio file.",
  "type": "object",
  "required": ["blocks"],
  "additionalProperties": false,
  "properties": {
    "blocks": {
      "type": "array",
      "minItems": 1,
      "description": "Ordered list of audio blocks to render.",
      "items": {
        "type": "object",
        "required": ["type"],
        "oneOf": [
          { "$ref": "#/$defs/tts" },
          { "$ref": "#/$defs/sfx" },
          { "$ref": "#/$defs/silence" }
        ]
      }
    }
  },
  "$defs": {
    "tts": {
      "type": "object",
      "description": "Text-to-speech narration block.",
      "required": ["type", "voice", "text"],
      "additionalProperties": false,
      "properties": {
        "type": { "const": "tts" },
        "voice": {
          "type": "string",
          "description": "ElevenLabs voice ID."
        },
        "text": {
          "type": "string",
          "minLength": 1,
          "description": "Text to synthesize."
        },
        "model": {
          "type": "string",
          "default": "eleven_multilingual_v2",
          "description": "ElevenLabs model ID."
        },
        "stability": {
          "type": "number",
          "minimum": 0,
          "maximum": 1,
          "description": "Voice stability (0.0–1.0)."
        },
        "similarity_boost": {
          "type": "number",
          "minimum": 0,
          "maximum": 1,
          "description": "Voice similarity boost (0.0–1.0)."
        },
        "style": {
          "type": "number",
          "minimum": 0,
          "maximum": 1,
          "description": "Style exaggeration (0.0–1.0)."
        },
        "speed": {
          "type": "number",
          "minimum": 0.5,
          "maximum": 2.0,
          "description": "Playback speed multiplier (0.5–2.0)."
        }
      }
    },
    "sfx": {
      "type": "object",
      "description": "Sound effect block generated from a text prompt.",
      "required": ["type", "text"],
      "additionalProperties": false,
      "properties": {
        "type": { "const": "sfx" },
        "text": {
          "type": "string",
          "minLength": 1,
          "description": "Text prompt describing the sound effect."
        },
        "background": {
          "type": "boolean",
          "default": false,
          "description": "When true, the SFX is mixed (overlaid) onto the next TTS block instead of playing sequentially."
        },
        "duration": {
          "type": "number",
          "minimum": 0.5,
          "maximum": 22,
          "description": "Duration in seconds (0.5–22.0)."
        }
      }
    },
    "silence": {
      "type": "object",
      "description": "Silent pause.",
      "required": ["type", "duration"],
      "additionalProperties": false,
      "properties": {
        "type": { "const": "silence" },
        "duration": {
          "type": "number",
          "exclusiveMinimum": 0,
          "description": "Duration in seconds (must be positive)."
        }
      }
    }
  }
}`

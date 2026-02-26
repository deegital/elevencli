# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.1.0] - 2026-02-26

### Added

- Text-to-speech command (`tts`) with voice selection, model choice, and file output
- Sound effects generation command (`sfx`) with duration and prompt-based generation
- Voice listing command (`voices`) to browse available ElevenLabs voices
- Audiobook generation command (`audiobook`) from structured JSON input
- ASCII banner display on root command
- Configuration via environment variable, flag, or `~/.elevencli.yaml`
- Cross-platform builds (macOS, Linux, Windows) via goreleaser
- Homebrew tap installation support

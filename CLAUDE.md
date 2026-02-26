# Project Guidelines

## Git Commits

- Never add "Co-Authored-By" lines or any other Claude/AI attribution to commit messages.
- Follow the Conventional Commits specification (https://www.conventionalcommits.org):
  - Format: `<type>(<optional scope>): <description>`
  - Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`
  - Use lowercase for type and description.
  - Keep the subject line under 72 characters.
  - Add a body separated by a blank line when more context is needed.
  - Use `!` after the type/scope for breaking changes (e.g., `feat!: remove legacy API`).

## Changelog

- Keep `CHANGELOG.md` up to date with every user-facing change.
- Follow the [Keep a Changelog](https://keepachangelog.com) format.
- Use the standard change types: Added, Changed, Deprecated, Removed, Fixed, Security.
- Place new entries under the `[Unreleased]` section.
- Never remove or rewrite existing entries.

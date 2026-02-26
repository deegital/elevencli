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

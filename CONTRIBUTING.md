# Contributing to sz (essenz)

Thanks for your interest in improving sz — a sharp Unix tool to distill the web. This guide covers setup, workflow, and PR expectations.

## Ways to contribute
- Bug reports and minimal repros
- Feature requests with concrete use cases
- Documentation improvements and examples
- Code changes that are small and focused

## Project setup

Requirements:
- Go 1.21+
- Chrome/Chromium installed (for JavaScript‑heavy pages)
- Optional: asdf to pin tool versions (`.tool-versions`)

Recommended steps:
```bash
# Clone
git clone https://github.com/jewell-lgtm/essenz
cd essenz

# Optional: install pinned tool versions
asdf install  # if you use asdf

# Install pre-commit hooks (fmt, vet, lint, test, conventional commits)
make setup-pre-commit

# Fetch deps and run the full quality gate
go mod download
make check  # fmt, vet, lint, test

# Build the CLI
make build
./sz --help
```

Useful make targets:
- `make build` — builds `sz`
- `make test` — runs tests
- `make test-sandbox` — runs tests using local caches for restricted envs
- `make lint` — runs golangci-lint (configured via `.golangci.yml`)
- `make check` — fmt, vet, lint, test

## Development workflow

```bash
# Create a feature branch
git checkout -b feature/your-change

# Make small, focused changes and keep tests green
make check

# Commit with Conventional Commits
# Examples: feat:, fix:, docs:, refactor:, test:, chore:
git commit -m "feat: add reader-view flag to fetch"

# Push and open a PR
git push -u origin feature/your-change
```

Guidelines:
- Prefer opening an issue first for non-trivial changes
- Keep PRs small; split large changes into increments
- Update docs/examples when behavior changes
- Add/adjust tests where it makes sense (see `specs/` and Go tests)

## Running locally

```bash
# Distill a page
./sz https://example.com/article > article.md

# Use raw mode (original HTML)
./sz --raw https://example.com

# TL;DR summarization (optional LLM)
export OPENAI_API_KEY=sk-...
./sz tldr https://example.com/article
./sz tldr --summary-length short /path/to/local.html
```

Advanced TL;DR flags: `--model`, `--base-url`, `--timeout`. The tool sends only distilled content to your provider.

## Commit messages

This repo enforces Conventional Commits via a commit‑msg hook:
- `feat: ...` new user‑facing functionality
- `fix: ...` bug fixes
- `docs: ...` documentation only
- `refactor: ...` code change without behavior change
- `test: ...` adding or fixing tests
- `chore: ...` tooling/infra

## Opening issues

- Issues: https://github.com/jewell-lgtm/essenz/issues
- Include environment, steps to reproduce, expected vs actual behavior, and logs if relevant.

## License

By contributing, you agree your contributions are licensed under the MIT license of this repository.

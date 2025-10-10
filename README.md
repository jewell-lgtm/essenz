# sz — a sharp Unix tool to distill the web

Get the content of any web page as clean, readable Markdown. Optionally ask an LLM for a TL;DR. Designed to behave like any other Unix tool.

Why “sz”? It’s short for essenz — the essence of a page.

## Features

- Skips cookie banners and many soft paywall overlays
- Handles JavaScript-heavy pages and SPA dynamic content
- Uses Google Chrome for rendering when needed, or direct HTTP fetch for static pages

## Philosophy

- Do one thing well: extract what matters from web pages
- Unix-friendly: stdout by default, composable with pipes and redirects
- Private by default: runs locally; LLM use is opt‑in for `tldr`
- Predictable: clear flags, sensible defaults, no surprises

## Quick Start

Extract the essence of a page (Markdown to stdout):

```bash
sz https://example.com/article
```

Save to a file or pipe into your tools:

```bash
sz https://example.com/article > article.md
sz https://news.ycombinator.com | rg -i "privacy|security"
sz /path/to/local.html | less -R
```

Prefer the original HTML? Use raw mode:

```bash
sz --raw https://example.com
```

## How It Works

- Render: Starts or reuses a lightweight Google Chrome daemon to load the page, run JavaScript, and wait for content readiness.
- Extract: Builds a clean text-node tree and filters boilerplate (nav, ads, cookie overlays) to focus on primary content.
- Rank: Scores blocks using simple heuristics (tag weight, length, link density, position) to surface what matters first.
- Render Markdown: Converts the result to readable Markdown and prints to stdout (or `--raw` to output original HTML).
- Optional TL;DR: `sz tldr` sends only the distilled content to your configured LLM for a concise summary.

## TL;DR Summaries (optional LLM)

Provide an API key via `OPENAI_API_KEY` (or pass `--api-key`). Then:

```bash
# One‑liner summary of an article
OPENAI_API_KEY=sk-... sz tldr https://example.com/article

# Control length (short|medium|long)
sz tldr --summary-length short https://example.com/article

# Summarize a local file
sz tldr /path/to/article.html
```

Advanced flags: `--model`, `--base-url`, `--timeout`. By default, only the distilled content is sent to your provider for summarization.

## Installation

Requirements:
- Chrome/Chromium installed for JavaScript‑heavy sites (falls back to HTTP fetch otherwise)
- Go 1.24+ (only needed for building from source)

### Homebrew (Recommended)

```bash
brew tap jewell-lgtm/essenz
brew install sz
```

### Other Methods

```bash
# Manual from release
curl -L https://github.com/jewell-lgtm/essenz/archive/refs/tags/v0.1.0.tar.gz -o essenz-0.1.0.tar.gz
tar -xzf essenz-0.1.0.tar.gz
cd essenz-0.1.0
go build -o sz ./cmd/essenz
sudo mv sz /usr/local/bin/

# From source (latest development)
git clone https://github.com/jewell-lgtm/essenz
cd essenz
make build && sudo mv sz /usr/local/bin/
```

See [INSTALLATION.md](INSTALLATION.md) for detailed installation instructions and troubleshooting.

## Links

- Issues: https://github.com/jewell-lgtm/essenz/issues
- Contributing guidelines: CONTRIBUTING.md

## Contributing

We welcome PRs! A good PR:
- Discusses the change in an issue first (when in doubt)
- Includes a small, focused diff
- Passes `make check` locally

For setup, workflow, and commit conventions, see CONTRIBUTING.md.

## License

MIT — see `LICENSE`.

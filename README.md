# sz — a sharp Unix tool to distill the web

Get the content of any web page as clean, readable Markdown. Optionally ask an LLM for a TL;DR. Designed to behave like any other Unix tool.

Why “sz”? It’s short for essenz — the essence of a page.

## Features

- Skips cookie banners and many soft paywall overlays
- Handles JavaScript-heavy pages and SPA dynamic content
- Tiered fetch engines: plain HTTP + readability for static pages, the lightweight [Lightpanda](https://lightpanda.io) browser for JavaScript, and Google Chrome as a heavy fallback

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

## Fetch Engines

`sz fetch` chooses how to retrieve a URL with `--engine` (default `auto`):

| Engine | JavaScript | Footprint | Use when |
|---|---|---|---|
| `http` | no | tiny | static articles, news, docs |
| `lightpanda` | yes | ~80 MB | SPAs / JS-rendered content |
| `chrome` | yes | heavy | sites Lightpanda can't yet render |
| `auto` | as needed | adaptive | default — escalates http → lightpanda → chrome |

```bash
sz fetch https://example.com                  # auto (default)
sz fetch --engine http https://example.com    # force the light path
sz fetch --engine lightpanda https://app.dev  # force JS rendering
ESSENZ_ENGINE=chrome sz fetch https://app.dev  # via environment
```

`auto` starts with the lightest engine and escalates only when a page looks
JavaScript-dependent, so most pages never spin up a browser.

The Lightpanda binary is downloaded and cached automatically on first use. To use
your own build, set `ESSENZ_LIGHTPANDA_PATH=/path/to/lightpanda`.

## How It Works

- Render: Picks a fetch engine (HTTP, Lightpanda, or Chrome) to load the page, run JavaScript if needed, and wait for content readiness.
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
- Nothing extra for most pages — static pages use a built-in HTTP fetch, and the lightweight Lightpanda browser is auto-downloaded on first use for JavaScript-heavy sites
- Chrome/Chromium only needed for `--engine chrome` (the heavy fallback)
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

# sz — a sharp Unix tool to distill the web

Get the content of any web page as clean, readable Markdown. Optionally ask an LLM for a TL;DR. Designed to behave like any other Unix tool.

Why “sz”? It’s short for essenz — the essence of a page.

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
- Go 1.21+ (for building from source)
- Chrome/Chromium installed for JavaScript‑heavy sites (falls back to HTTP fetch otherwise)

Choose one:

```bash
# Prebuilt binary (macOS/Linux)
curl -L "https://github.com/jewell-lgtm/essenz/releases/latest/download/sz-$(uname -s)-$(uname -m)" -o sz
chmod +x sz && sudo mv sz /usr/local/bin/

# Go install (installs binary named 'essenz')
go install github.com/jewell-lgtm/essenz/cmd/essenz@latest
# Optional: symlink to 'sz'
sudo ln -sf "$(go env GOPATH)/bin/essenz" /usr/local/bin/sz

# From source
git clone https://github.com/jewell-lgtm/essenz
cd essenz
make build && sudo mv sz /usr/local/bin/
```

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

# Installation Guide for sz

## Quick Start (Recommended)

### Using Homebrew

The easiest way to install `sz` is via Homebrew:

```bash
brew tap jewell-lgtm/essenz
brew install sz
```

Verify the installation:

```bash
sz version
# Output: sz version 0.1.0
```

## Alternative Installation Methods

### Manual Installation from Release

1. Download the latest release:
```bash
curl -L -o essenz-0.1.0.tar.gz https://github.com/jewell-lgtm/essenz/archive/refs/tags/v0.1.0.tar.gz
```

2. Extract:
```bash
tar -xzf essenz-0.1.0.tar.gz
cd essenz-0.1.0
```

3. Build:
```bash
go build -o sz ./cmd/essenz
```

4. Install (choose one):

   **System-wide (requires sudo):**
   ```bash
   sudo mv sz /usr/local/bin/
   ```

   **User directory:**
   ```bash
   mkdir -p ~/bin
   mv sz ~/bin/
   # Add ~/bin to PATH if not already there
   echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc  # or ~/.bashrc
   source ~/.zshrc
   ```

5. Verify:
```bash
sz version
```

### Building from Source (Latest Development)

1. Clone the repository:
```bash
git clone https://github.com/jewell-lgtm/essenz.git
cd essenz
```

2. Build:
```bash
make build
```

3. The binary `sz` will be created in the current directory. Install it as shown above.

## Requirements

- **For Building**: Go 1.24.0 or later
- **For Running**: Chrome or Chromium (for headless browser functionality)

## Quick Test

After installation, test with a simple example:

```bash
# Extract content from a URL
sz https://example.com

# Show help
sz --help

# Show version
sz version
```

## Updating

### Homebrew

```bash
brew update
brew upgrade sz
```

### Manual

Download and install the latest release following the manual installation steps above.

## Uninstalling

### Homebrew

```bash
brew uninstall sz
brew untap jewell-lgtm/essenz
```

### Manual

```bash
# Remove the binary
rm /usr/local/bin/sz
# or
rm ~/bin/sz
```

## Troubleshooting

### sz: command not found

Make sure the installation directory is in your PATH:

```bash
echo $PATH
```

For user directory installations, ensure `~/bin` is in your PATH:

```bash
export PATH="$HOME/bin:$PATH"
```

Add this to your shell configuration file (`~/.zshrc` or `~/.bashrc`) to make it permanent.

### Chrome/Chromium not found

Install Chrome or Chromium:

**macOS (Homebrew):**
```bash
brew install --cask google-chrome
# or
brew install chromium
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install chromium-browser
# or
sudo apt-get install google-chrome-stable
```

## More Information

- **Repository**: https://github.com/jewell-lgtm/essenz
- **Issues**: https://github.com/jewell-lgtm/essenz/issues
- **Documentation**: See [README.md](README.md) for full documentation

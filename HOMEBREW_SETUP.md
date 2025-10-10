# Homebrew Installation Setup for sz

This document contains the complete setup for installing `sz` via Homebrew.

## GitHub Release ✅

**Release Created**: https://github.com/jewell-lgtm/essenz/releases/tag/v0.1.0

- **Tag**: v0.1.0
- **Source Tarball**: https://github.com/jewell-lgtm/essenz/archive/refs/tags/v0.1.0.tar.gz
- **SHA256**: `cbd80a23d48fa40e44ca1765aeb7fe574103032f8dad4cc438b5a90d78a615dc`

## Homebrew Tap Repository Setup

### 1. Create the Tap Repository

You need to create a new GitHub repository named `homebrew-essenz` in your account:

```bash
# Create the repository on GitHub
gh repo create homebrew-essenz --public --description "Homebrew tap for essenz (sz) CLI tool"

# Clone it locally
git clone git@github.com:jewell-lgtm/homebrew-essenz.git
cd homebrew-essenz

# Create the Formula directory
mkdir -p Formula

# Copy the formula file (provided below)
# Place sz.rb in Formula/sz.rb
```

### 2. Homebrew Formula

Create `Formula/sz.rb` with the following content:

```ruby
class Sz < Formula
  desc "CLI web browser that captures the essence of web pages"
  homepage "https://github.com/jewell-lgtm/essenz"
  url "https://github.com/jewell-lgtm/essenz/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "cbd80a23d48fa40e44ca1765aeb7fe574103032f8dad4cc438b5a90d78a615dc"
  license "MIT"
  head "https://github.com/jewell-lgtm/essenz.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(output: bin/"sz", ldflags: "-s -w"), "./cmd/essenz"
  end

  test do
    assert_match "sz version 0.1.0", shell_output("#{bin}/sz version")
  end
end
```

### 3. Commit and Push the Formula

```bash
cd homebrew-essenz
git add Formula/sz.rb
git commit -m "Add sz formula v0.1.0"
git push origin main
```

### 4. Optional: Add a README

Create `README.md`:

```markdown
# Homebrew Tap for essenz (sz)

This is the official Homebrew tap for [essenz](https://github.com/jewell-lgtm/essenz) - a CLI web browser that captures the essence of web pages.

## Installation

\`\`\`bash
brew tap jewell-lgtm/essenz
brew install sz
\`\`\`

## Usage

\`\`\`bash
# Extract content from URL
sz https://example.com

# Extract content from local file
sz /path/to/file.html

# Get raw HTML without processing
sz --raw https://example.com

# Show version
sz version
\`\`\`

## Documentation

See the [main repository](https://github.com/jewell-lgtm/essenz) for full documentation.
\`\`\`

Then commit:

```bash
git add README.md
git commit -m "Add README"
git push origin main
```

## Installation Instructions for Users

### Using Homebrew (Recommended)

```bash
# Add the tap
brew tap jewell-lgtm/essenz

# Install sz
brew install sz

# Verify installation
sz version
```

### Manual Installation from Source

```bash
# Download and extract
curl -L -o essenz-0.1.0.tar.gz https://github.com/jewell-lgtm/essenz/archive/refs/tags/v0.1.0.tar.gz
tar -xzf essenz-0.1.0.tar.gz
cd essenz-0.1.0

# Build
go build -o sz ./cmd/essenz

# Install (requires sudo for /usr/local/bin)
sudo mv sz /usr/local/bin/

# Or install to user directory
mkdir -p ~/bin
mv sz ~/bin/
# Add ~/bin to PATH if not already there
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc  # or ~/.bashrc
source ~/.zshrc  # or source ~/.bashrc

# Verify
sz version
```

## Testing the Formula

Before making it public, test the formula locally:

```bash
# Test the formula audit
brew audit --new-formula Formula/sz.rb

# Test installation locally
brew install --build-from-source Formula/sz.rb

# Test that it works
sz version

# Uninstall test
brew uninstall sz
```

## Updating the Formula (For Future Releases)

When releasing a new version:

1. Create a new GitHub release with a new tag (e.g., v0.2.0)
2. Get the new SHA256 checksum:
   ```bash
   curl -L https://github.com/jewell-lgtm/essenz/archive/refs/tags/v0.2.0.tar.gz | shasum -a 256
   ```
3. Update `Formula/sz.rb`:
   - Change the `url` to the new tag
   - Update the `sha256` with the new checksum
   - Update the test assertion if the version string changed
4. Commit and push the changes
5. Users can update with:
   ```bash
   brew update
   brew upgrade sz
   ```

## Requirements

- **For Building**: Go 1.24.0 or later
- **For Running**: Chrome/Chromium (for headless browser functionality)

## Verification

The formula has been tested with:
- ✅ Source tarball extracts correctly
- ✅ Build command completes successfully
- ✅ Binary runs and shows correct version
- ✅ SHA256 checksum verified
- ✅ LICENSE file included in source distribution

## Next Steps

1. Create the `homebrew-essenz` repository on GitHub
2. Add the formula file as shown above
3. Test the installation locally
4. Share the installation instructions with users:
   ```bash
   brew tap jewell-lgtm/essenz
   brew install sz
   ```

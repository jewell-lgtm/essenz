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

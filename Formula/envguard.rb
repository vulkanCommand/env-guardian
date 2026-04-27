class Envguard < Formula
  desc "Validate, secure, encrypt, and ship environment files safely"
  homepage "https://github.com/vulkanCommand/env-guardian"
  url "https://github.com/vulkanCommand/env-guardian.git",
      tag:      "v0.1.13"
  license "MIT"
  head "https://github.com/vulkanCommand/env-guardian.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", bin/"envguard", "./cmd/envguard"
  end

  test do
    assert_match "0.1.13", shell_output("#{bin}/envguard version")
  end
end

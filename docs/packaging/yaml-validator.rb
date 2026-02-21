# Homebrew formula для yaml-validator
# Использование: brew install --formula yaml-validator.rb

class YamlValidator < Formula
  desc "YAML validation tool for DevOps"
  homepage "https://github.com/qmish/yaml-validator"
  version "1.50.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/qmish/yaml-validator/releases/download/v1.50.0/yaml-validator-v1.50.0-darwin-amd64.tar.gz"
      sha256 ""
    end
    on_arm do
      url "https://github.com/qmish/yaml-validator/releases/download/v1.50.0/yaml-validator-v1.50.0-darwin-arm64.tar.gz"
      sha256 ""
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/qmish/yaml-validator/releases/download/v1.50.0/yaml-validator-v1.50.0-linux-amd64.tar.gz"
      sha256 ""
    end
    on_arm do
      url "https://github.com/qmish/yaml-validator/releases/download/v1.50.0/yaml-validator-v1.50.0-linux-arm64.tar.gz"
      sha256 ""
    end
  end

  def install
    bin.install Dir["yaml-validator-*"].first => "yaml-validator"
  end

  test do
    system "#{bin}/yaml-validator", "version"
  end
end

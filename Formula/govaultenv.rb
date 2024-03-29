# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Govaultenv < Formula
  desc ""
  homepage ""
  version "1.2.10"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/jamhed/govaultenv/releases/download/v1.2.10/govaultenv_1.2.10_darwin_amd64.tar.gz"
      sha256 "a91808ef507c381cdde94969eedc123724d7bed7021c048d8069aafd22da1ab8"

      def install
        bin.install "govaultenv"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/jamhed/govaultenv/releases/download/v1.2.10/govaultenv_1.2.10_darwin_arm64.tar.gz"
      sha256 "90217dae1bdce955688d6534e39bbc05fe4b6bb28f27c7ce44789dff09f65944"

      def install
        bin.install "govaultenv"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/jamhed/govaultenv/releases/download/v1.2.10/govaultenv_1.2.10_linux_arm64.tar.gz"
      sha256 "589526af4180293b50871481d352ad3b986e91dba7496a3ce44343b2e537d548"

      def install
        bin.install "govaultenv"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/jamhed/govaultenv/releases/download/v1.2.10/govaultenv_1.2.10_linux_amd64.tar.gz"
      sha256 "25cdfacc8847aaca165462a99bb1b67f24412e803a88b541cd57fefc6ddf911e"

      def install
        bin.install "govaultenv"
      end
    end
  end
end

#!/usr/bin/env bash
set -euo pipefail
if [ "$#" -ne 3 ]; then echo 'usage: generate-formula.sh VERSION ASSET_DIR OUTPUT' >&2; exit 2; fi
version=$1
assets=$2
output=$3
[[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo 'Invalid stable version' >&2; exit 1; }
for platform in darwin-amd64 darwin-arm64 linux-amd64 linux-arm64; do
  name="infra-$platform"
  line=$(awk -v n="$name" '$2 == n {print}' "$assets/checksums.txt")
  [ "$(printf '%s\n' "$line" | wc -l | tr -d ' ')" = 1 ] && [ -n "$line" ] || { echo "Missing or duplicate checksum for $name" >&2; exit 1; }
  (cd "$assets" && printf '%s\n' "$line" | shasum -a 256 -c -)
done
sha() { awk -v n="infra-$1" '$2 == n {print $1}' "$assets/checksums.txt"; }
mkdir -p "$(dirname "$output")"
cat > "$output" <<FORMULA
class Infra < Formula
  desc "Secure infrastructure access CLI"
  homepage "https://github.com/ktsstudio/infra"
  version "$version"
  license "Elastic-2.0"

  on_macos do
    on_arm do
      url "https://github.com/ktsstudio/infra/releases/download/v$version/infra-darwin-arm64"
      sha256 "$(sha darwin-arm64)"
    end
    on_intel do
      url "https://github.com/ktsstudio/infra/releases/download/v$version/infra-darwin-amd64"
      sha256 "$(sha darwin-amd64)"
    end
  end
  on_linux do
    on_arm do
      url "https://github.com/ktsstudio/infra/releases/download/v$version/infra-linux-arm64"
      sha256 "$(sha linux-arm64)"
    end
    on_intel do
      url "https://github.com/ktsstudio/infra/releases/download/v$version/infra-linux-amd64"
      sha256 "$(sha linux-amd64)"
    end
  end

  def install
    bin.install Dir["infra-*"].first => "infra"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/infra version")
  end
end
FORMULA

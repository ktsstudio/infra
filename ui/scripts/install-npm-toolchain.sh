#!/bin/sh
set -eu

npm install --global npm@12.2.1

# npm 12.2.1 still bundles vulnerable ip-address/undici releases. Replace only
# these bundles with patched compatible releases until npm updates its bundle.
# undici must retain major 6 because npm's bundled callers use that API.
patch_dir=$(mktemp -d)
trap 'rm -rf "$patch_dir"' EXIT HUP INT TERM
npm install --prefix "$patch_dir" --ignore-scripts --no-audit --no-fund \
  --package-lock=false ip-address@10.7.5 undici@6.29.1
npm_modules="$(npm root --global)/npm/node_modules"
for package in ip-address undici; do
  rm -rf "$npm_modules/$package"
  cp -R "$patch_dir/node_modules/$package" "$npm_modules/$package"
done
npm --version

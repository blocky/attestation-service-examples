#!/usr/bin/env bash

set -e

bin=$1
commit=$2
os=$3
arch=$4

mkdir -p $bin

echo -n "Getting commit of current bky-as..."
if [[ -x "$bin/bky-as" ]]; then
    tmp_out=$(mktemp)
    trap 'rm -f "$tmp_out"' EXIT
    if "$bin/bky-as" inspect >"$tmp_out"; then
        current_version_commit=$(jq -r .Build.Commit < "$tmp_out")
        echo "'$current_version_commit'"
    else
        echo "Failed to inspect bky-as binary" >&2
        cat "$tmp_out" >&2
        exit 1
    fi
else
    echo "no current version"
fi

if [[ $commit == "latest" ]]; then
  echo -n "Getting the latest commit..."
  commit=$(gh api repos/blocky/delphi/commits --jq '.[0].sha')
  echo "'$commit'"
fi

if [[ $current_version_commit != "$commit" ]]; then
    echo "Versions differ ...updating"
    aws s3 cp "s3://blocky-internal-release/delphi/cli/${commit}/${os}_${arch}" "${bin}/bky-as"
    chmod +x "$bin/bky-as"
else
    echo "Version up to date"
fi

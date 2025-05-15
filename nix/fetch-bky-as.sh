#!/usr/bin/env bash

set -e

bin=$1
commit=$2
os=$3
arch=$4

mkdir -p $bin

echo -n "Getting commit of current bky-as..."
current_version_commit=""
if [[ -e $bin/bky-as ]]; then
    current_version_commit=$($bin/bky-as inspect | jq -r .Build.Commit)
    echo "'$current_version_commit'"
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

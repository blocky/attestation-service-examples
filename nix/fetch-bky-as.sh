#!/usr/bin/env bash

set -e

bin=$1
os=$2
arch=$3

mkdir -p $bin

echo -n "Getting commit of current bky-as..."
current_version_commit=""
if [[ -e $bin/bky-as ]]; then
    current_version_commit=$($bin/bky-as inspect | jq -r .Build.Commit)
    echo "'$current_version_commit'"
else
    echo "no current version"
fi

echo -n "Getting the latest commit..."
latest_commit=$(gh api repos/blocky/delphi/commits --jq '.[0].sha')
echo "'$latest_commit'"

if [[ $current_version_commit != $latest_commit ]]; then
    echo "Versions differ ...updating"
    aws s3 cp s3://blocky-internal-release/delphi/cli/$latest_commit/${os}_${arch} ${bin}/bky-as
    chmod +x "$bin/bky-as"
else
    echo "Version up to date"
fi

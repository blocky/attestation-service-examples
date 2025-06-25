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

if [[ -z $AWS_ACCESS_KEY_ID || -z $AWS_SECRET_ACCESS_KEY || -z $AWS_SESSION_TOKEN ]]; then
    echo "Error: AWS credentials are missing!"
    exit 1
fi

if [[ $current_version_commit != "$commit" ]]; then
    echo "Versions differ ...updating"
    aws s3 cp "s3://blocky-internal-release/delphi/cli/${commit}/${os}_${arch}" "${bin}/bky-as" --region us-west-2
    chmod +x "$bin/bky-as"
else
    echo "Version up to date"
fi

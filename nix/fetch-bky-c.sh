#!/usr/bin/env bash

set -e

bin=$1
version=$2
os=$3
arch=$4
#arch="amd64"

REPO="attestation-service-cli"
APP="bky-as"

mkdir -p $bin

if [[ "${version}" == latest ]]; then
  release=$(curl -s \
                  -H "Accept: application/vnd.github+json" \
                  -H "X-GitHub-Api-Version: 2022-11-28" \
                  "https://api.github.com/repos/blocky/${REPO}/releases" \
            | jq '
                map(select(.draft == false))
                | sort_by(.published_at)
                | reverse
                | .[0]
            ')
  artifact="${APP}_${os}_${arch}"
  echo "Wanted artifact: ${artifact}"
  url=$(echo "$release" | jq -r --arg name "${artifact}" '
    .assets[] | select(.name == $name) | .browser_download_url
  ')
  echo $url
else
   base="https://github.com/blocky/${REPO}/releases/download"
   artifact="${APP}_${os}_${arch}"
   url="${base}/${version}/${artifact}"
  echo $url
fi

echo "Downloading cli from: ${url}"
if ! curl --silent --location --fail --show-error "${url}" -o "${bin}/bky-c"; then
    exit 1
fi
echo "Downloaded bky-c"
chmod +x "${bin}/bky-c"
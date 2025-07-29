#!/usr/bin/env bash

set -e

bin=$1
version=$2
os=$3
arch=$4

REPO="compiler"
APP="bky-c"

mkdir -p $bin

if [[ "${version}" == latest ]]; then
  echo "Downloading latest ${APP} release"
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
  url=$(echo "$release" | jq -r --arg name "${artifact}" '
    .assets[] | select(.name == $name) | .browser_download_url
  ')
else
  echo "Downloading tagged ${APP} release: ${version}"
   base="https://github.com/blocky/${REPO}/releases/download"
   artifact="${APP}_${os}_${arch}"
   url="${base}/${version}/${artifact}"
fi

echo "Downloading ${APP} cli from: ${url}"
if ! curl --silent --location --fail --show-error "${url}" -o "${bin}/${APP}"; then
    exit 1
fi

echo "Downloaded ${APP}"
chmod +x "${bin}/${APP}"
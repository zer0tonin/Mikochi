#!/bin/bash

get_architecture() {
    case $(uname -m) in
        x86_64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        armv7l|armhf)
            echo "armv7"
            ;;
        *)
            # fallback to dpkg if unsure
            dpkg --print-architecture
            ;;
    esac
}

zip_file="mikochi-linux-$(get_architecture).tar.gz"
github_repo="zer0tonin/mikochi"

apt-get update
apt-get install -y curl
apt-get install -y jq

cd /
tar -czvf "$zip_file" /app

upload_url=$(curl -X GET -H "Authorization: token $GITHUB_TOKEN" "https://api.github.com/repos/$github_repo/releases/$RELEASE" | jq -r '.upload_url')

upload_url="${upload_url%\{*}"
echo "Upload URL: $upload_url"

curl -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Content-Type: application/tar+gzip" \
  --data-binary @"$zip_file" \
  "$upload_url?name=$zip_file"

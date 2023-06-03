#!/bin/bash
version=$1
token=$2

git tag $version
git push origin main --tags
docker buildx build --platform linux/amd64,linux/arm64 -t zer0tonin/mikochi:latest -t zer0tonin/mikochi:$version --push .

release_id=$(curl -X POST -H "Authorization: token $token" -H "Content-Type: application/json" -d '{"tag_name": "'"$version"'"}' -s "https://api.github.com/repos/zer0tonin/mikochi/releases" | jq -r '.id')
docker run --env RELEASE=$release_id --env GITHUB_TOKEN=$token --platform linux/arm64 -v $(PWD)/zip-upload.sh:/zip-upload.sh -it zer0tonin/mikochi:latest /bin/bash /zip-upload.sh
docker run --env RELEASE=$release_id --env GITHUB_TOKEN=$token --platform linux/amd64 -v $(PWD)/zip-upload.sh:/zip-upload.sh -it zer0tonin/mikochi:latest /bin/bash /zip-upload.sh

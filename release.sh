#!/bin/bash
version=$1
token=$2

git tag $version
git push origin main --tags
docker buildx build --platform linux/amd64,linux/arm64 -t zer0tonin/mikochi:latest -t zer0tonin/mikochi:$version --push .
#docker buildx build --platform linux/arm/v7 -t zer0tonin/mikochi:latest-armv7 -t zer0tonin/mikochi:$version-armv7 --push -f ./Dockerfile-armv7 .

release_id=$(curl -X POST -H "Authorization: token $token" -H "Content-Type: application/json" -d '{"tag_name": "'"$version"'"}' -s "https://api.github.com/repos/zer0tonin/mikochi/releases" | jq -r '.id')
echo $release_id
#docker run --env RELEASE=$release_id --env GITHUB_TOKEN=$token --platform linux/arm/v7 -v $(PWD)/tar-upload.sh:/tar-upload.sh -it zer0tonin/mikochi:latest-armv7 /bin/bash /tar-upload.sh
docker run --env RELEASE=$release_id --env GITHUB_TOKEN=$token --platform linux/arm64 -v $(PWD)/tar-upload.sh:/tar-upload.sh -it zer0tonin/mikochi:latest /bin/bash /tar-upload.sh
docker run --env RELEASE=$release_id --env GITHUB_TOKEN=$token --platform linux/amd64 -v $(PWD)/tar-upload.sh:/tar-upload.sh -it zer0tonin/mikochi:latest /bin/bash /tar-upload.sh

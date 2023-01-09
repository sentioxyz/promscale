TAG=us-west1-docker.pkg.dev/sentio-352722/sentio/promscale:0.17.0-patch.1

docker build -t $TAG -f build/Dockerfile --build-arg TARGETARCH=amd64 .
docker push $TAG
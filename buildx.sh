set -e
docker buildx use localremote_builder

# docker buildx build --push \
# --platform linux/amd64,linux/arm64 \
# --tag pedromol/catoso:base -f Dockerfile.base .

# docker buildx build --push \
# --platform linux/amd64,linux/arm64 \
# --tag pedromol/catoso:latest .

# docker buildx build --push \
# --platform linux/amd64 \
# --tag pedromol/catoso:base-cuda -f Dockerfile.base-cuda .

docker buildx build --push \
--platform linux/amd64 \
--tag pedromol/catoso:cuda -f Dockerfile.cuda .
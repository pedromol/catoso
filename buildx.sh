sudo docker buildx build --push \
--platform linux/amd64,linux/arm64 \
--tag pedromol/catoso:latest .

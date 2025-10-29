build:
	docker buildx build --platform linux/amd64,linux/arm64 -t o9yst03/odtimagereplacer:1.0.1  -t o9yst03/odtimagereplacer:latest --push  .
.PHONY: build
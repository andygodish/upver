IMAGE_NAME := upver
VERSION := latest

# Build the Docker image (single architecture)
.PHONY: docker-build
docker-build:
	@echo "Building $(IMAGE_NAME):$(VERSION)..."
	docker build \
		-t $(IMAGE_NAME):$(VERSION) \
		.
	@echo "Build complete: $(IMAGE_NAME):$(VERSION)"
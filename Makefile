# Define target platforms, image builder and the fully qualified image name.
TARGET_PLATFORMS ?= linux/amd64,linux/arm64

TAG ?= dev
REPO ?= rancher
IMAGE ?= rancher-ai-llm-mock
IMAGE_NAME ?= $(REPO)/$(IMAGE):$(TAG)

.PHONY: push-image
push-image:
	docker buildx build \
	  $(IID_FILE_FLAG) \
	  $(BUILDX_ARGS) \
	  --platform=$(TARGET_PLATFORMS) \
	  --tag $(IMAGE_NAME) \
	  --push \
	  .

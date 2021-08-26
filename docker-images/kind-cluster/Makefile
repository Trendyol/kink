VERSION ?= $(shell cat VERSION)

IMG_REPO ?= trendyoltech
IMG_TAG ?= $(VERSION)
BASE_IMAGES ?= buster 

BUILD_TARGETS := $(BASE_IMAGES:%=build-%)
PUSH_TARGETS := $(BASE_IMAGES:%=push-%)

all: $(BUILD_TARGETS)

push: $(PUSH_TARGETS)

$(BUILD_TARGETS):
	docker build -t $(IMG_REPO)/kind-cluster-$(@:build-%=%):$(IMG_TAG) -f Dockerfile.$(@:build-%=%) .

$(PUSH_TARGETS):
	docker push $(IMG_REPO)/kind-cluster-$(@:push-%=%):$(IMG_TAG)

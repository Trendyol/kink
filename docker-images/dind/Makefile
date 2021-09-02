VERSION ?= $(shell cat VERSION)

IMG_REPO ?= trendyoltech
IMG_TAG ?= $(VERSION)
BASE_IMAGES ?= buster buster-golang centos7 centos7-golang

BUILD_TARGETS := $(BASE_IMAGES:%=build-%)
PUSH_TARGETS := $(BASE_IMAGES:%=push-%)

all: $(BUILD_TARGETS)

push: $(PUSH_TARGETS)

$(BUILD_TARGETS):
	docker build -t $(IMG_REPO)/dind-$(@:build-%=%):$(IMG_TAG) -f Dockerfile.$(@:build-%=%) .

$(PUSH_TARGETS):
	docker push $(IMG_REPO)/dind-$(@:push-%=%):$(IMG_TAG)

VERSION ?= $(shell cat VERSION)

IMG_REPO ?= trendyoltech
IMG_TAG ?= $(VERSION)

all: kind-node

kind-node:
	docker build -t $(IMG_REPO)/kind-node:$(IMG_TAG) .

push:
	docker push $(IMG_REPO)/kind-node:$(IMG_TAG)

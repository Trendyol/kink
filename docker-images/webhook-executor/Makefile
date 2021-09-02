VERSION ?= $(shell cat VERSION)

IMG_REPO ?= trendyoltech
IMG_NAME ?= webhook-executor
IMG_TAG ?= $(VERSION)

all: build

build:
	docker build -t $(IMG_REPO)/$(IMG_NAME):$(IMG_TAG) .

push:
	docker push $(IMG_REPO)/$(IMG_NAME):$(IMG_TAG)

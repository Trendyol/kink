# [docker-images](https://github.com/jieyu/docker-images)

A collection of Docker images by [@jieyu](https://github.com/jieyu)

Article: [Running kind inside a kubernetes cluster for continuous integration](https://d2iq.com/blog/running-kind-inside-a-kubernetes-cluster-for-continuous-integration)

# Build your own kind-cluster

```bash
$ cd kind-cluster
$ docker build -t $(IMG_REPO)/kind-cluster-$(@:build-%=%):$(IMG_TAG) -f Dockerfile.$(@:build-%=%) .    
      - --platform=linux/amd64
      - --label=org.opencontainers.image.title=kink
      - --label=org.opencontainers.image.description={{ .ProjectName }}
      - --label=org.opencontainers.image.url=https://github.com/Trendyol/kink
      - --label=org.opencontainers.image.source=https://github.com/Trendyol/kink
      - --label=org.opencontainers.image.version={{ .Version }}
      - --label=org.opencontainers.image.created={{ time "2006-01-02T15:04:05Z07:00" }}
      - --label=org.opencontainers.image.revision={{ .FullCommit }}
      - --label=org.opencontainers.image.licenses=MIT
$ docker push $(IMG_REPO)/kind-cluster-$(@:push-%=%):$(IMG_TAG)
```
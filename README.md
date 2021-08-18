<div align="center">
  <img width="500" height="500" src="https://github.com/marcusolsson/gophers/raw/master/viking.png">
</div>


# kink (KinD in Kubernetes)

A helper CLI that facilitates to manage KinD clusters as Kubernetes pods

## Introduction

Before getting started, we should talk about a bit [KinD](https://kind.sigs.k8s.io) who is not familiar with this project. **_KinD_** is a tool for running local Kubernetes clusters using Docker container **_nodes_**. **_KinD_** was primarily designed for testing Kubernetes itself, but may be used for local development or CI.

So, what is **_kink_** then, where does this idea come from?

This CLI is basically an application that facilitates to run KinD cluster in Kubernetes Pod. There is a very detailed guide about how you can run KinD cluster in a Pod, for more detail, please [see](https://d2iq.com/blog/running-kind-inside-a-kubernetes-cluster-for-continuous-integration). Because this project is developed based on that blog post. The idea is that when you want to run ephemeral clusters by using projects like KinD in your CI/CD system instead of having test Kubernetes clusters, because it might cost more, you might want to run your KinD cluster in a Pod, especially if you are using Gitlab as a CI/CD solution and running your jobs as Kubernetes Pod. This project specifically aims to solve that problem. By using **_kink_**, you can easily manage whole lifecycle of your KinD cluster no matter how many they are as Kubernetes Pod.

## Usage

GOPRIVATE=“gitlab.trendyol.com” \
go install gitlab.trendyol.com/platform/base/poc/kink@latest

# if you want to build it from source code

Build kink and move kink binary to a file location on your system PATH

go build && mv ./kink /usr/local/bin/kink


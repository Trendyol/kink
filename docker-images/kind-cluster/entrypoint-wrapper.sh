#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# Startup Docker daemon and wait for it to be ready.
/entrypoint-original.sh bash -c "touch /dockerd-ready && sleep infinity" &
while [ ! -f /dockerd-ready ]; do sleep 1; done

echo "Setting up KIND cluster"

# Startup a KIND cluster with given configurations
API_SERVER_ADDRESS=${API_SERVER_ADDRESS:-"127.0.0.1"}
sed -i "s/apiServerAddress:$/apiServerAddress: ${API_SERVER_ADDRESS}/" kind-config.yaml

CERT_SANS=(${CERT_SANS:-""})
CERT_SANS+=(${API_SERVER_ADDRESS})
CERT_SANS+=($(hostname -i))
CERT_SANS+=(localhost)
CERT_SANS+=(127.0.0.1)

UNIQUE_CERT_SANS=($(echo "${CERT_SANS[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))

for hostname in "${UNIQUE_CERT_SANS[@]}"; do
cat <<EOF >> kind-config.yaml
- group: kubeadm.k8s.io
  version: v1beta2
  kind: ClusterConfiguration
  patch: |
    - op: add
      path: /apiServer/certSANs/-
      value: ${hostname}
EOF
done

kind create cluster --name=${KIND_CLUSTER_NAME:-""} --config=kind-config.yaml --image=${KIND_NODE_IMAGE-"registry.trendyol.com/platform/base/image/kind-node:v1.21.2"} --wait=900s

exec "$@"

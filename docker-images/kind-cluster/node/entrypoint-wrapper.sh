#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

CURRENT_CGROUP=$(grep systemd /proc/self/cgroup | cut -d: -f3)
CGROUP_SUBSYSTEMS=$(findmnt -lun -o source,target -t cgroup | grep "${CURRENT_CGROUP}" | awk '{print $2}')

/usr/local/bin/entrypoint-original echo "KIND entrypoint done"

mount --make-rprivate /sys/fs/cgroup

echo "${CGROUP_SUBSYSTEMS}" |
while IFS= read -r SUBSYSTEM; do
  # This is because we set Kubelet's cgroup-root to `/kubelet` by
  # default. We have to do that because otherwise, it'll collide
  # with the cgroups used by the Kubelet running on the host if we
  # run Konvoy docker cluster within a Kubernetes pod, resulting
  # random processes to be killed.
  mkdir -p "${SUBSYSTEM}/kubelet"
  if [ "${SUBSYSTEM}" == "/sys/fs/cgroup/cpuset" ]; then
    # This is needed. Otherwise, assigning process to the cgroup
    # (or any nested cgroup) would result in ENOSPC.
    cat "${SUBSYSTEM}/cpuset.cpus" > "${SUBSYSTEM}/kubelet/cpuset.cpus"
    cat "${SUBSYSTEM}/cpuset.mems" > "${SUBSYSTEM}/kubelet/cpuset.mems"
  fi
  # We need to perform a self bind mount here because otherwise,
  # systemd might delete the cgroup unintentionally before the
  # kubelet starts.
  mount --bind "${SUBSYSTEM}/kubelet" "${SUBSYSTEM}/kubelet"
done

exec "$@"

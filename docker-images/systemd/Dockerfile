# This Dockerfile is derived from:
# https://github.com/kubernetes-sigs/kind/blob/master/images/base/Dockerfile

# For systemd + docker configuration used below, see the following
# references:
# https://www.freedesktop.org/wiki/Software/systemd/ContainerInterface/
# https://developers.redhat.com/blog/2014/05/05/running-systemd-within-docker-container/
# https://developers.redhat.com/blog/2016/09/13/running-systemd-in-a-non-privileged-container/

ARG BASE_IMAGE="centos:7.6.1810"
FROM ${BASE_IMAGE}

# Install dependencies.
RUN yum -y update && \
    yum -y install systemd openssh-server openssh-clients libseccomp nfs-utils && \
    yum clean all && \
    find /lib/systemd/system/sysinit.target.wants/ -name "systemd-tmpfiles-setup.service" -delete && \
    rm -f /lib/systemd/system/multi-user.target.wants/* && \
    rm -f /etc/systemd/system/*.wants/* && \
    rm -f /lib/systemd/system/local-fs.target.wants/* && \
    rm -f /lib/systemd/system/sockets.target.wants/*udev* && \
    rm -f /lib/systemd/system/sockets.target.wants/*initctl* && \
    rm -f /lib/systemd/system/basic.target.wants/* && \
    rm -f /lib/systemd/system/anaconda.target.wants/*

# Tell systemd that it is running in docker (it will check for the
# container env). See details in:
# https://www.freedesktop.org/wiki/Software/systemd/ContainerInterface/
ENV container docker

# Systemd exits on SIGRTMIN+3, not SIGTERM (which re-executes it)
# https://bugzilla.redhat.com/show_bug.cgi?id=1201657
STOPSIGNAL SIGRTMIN+3

# Wrap systemd with a special entrypoint.This lets us set up some
# things before continuing on to systemd while preserving that systemd
# is PID 1.
COPY [ "entrypoint", "/usr/local/bin/" ]

# We need systemd to be PID1 to run the various services (docker,
# kubelet, etc.).
ENTRYPOINT [ "/usr/local/bin/entrypoint", "/sbin/init" ]

# TODO(bentheelder): deal with systemd MAC address assignment
# https://github.com/systemd/systemd/issues/3374#issuecomment-288882355
# https://github.com/systemd/systemd/issues/3374#issuecomment-339258483

#!/bin/bash
# This script runs some very basic commands to ensure that the newly build
# images are working correctly. Invoke as:
# ./image-checks.sh <image-tag> <registry-name> <arch>
# Kill with Ctrl + C once sidecar starts up successfully.
#
# To host different arch images on x86_64 you can run
# docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
TAG=$1
REGISTRY=${2:-gcr.io/google-containers}
ARCH=${3:-amd64}
echo "Verifying that iptables exists in node-cache image"
docker run --rm -it --entrypoint=iptables "${REGISTRY}"/k8s-dns-node-cache-"${ARCH}":"${TAG}"
echo "Verifying that node-cache binary exists in node-cache image"
docker run --rm -it --entrypoint=/node-cache "${REGISTRY}"/k8s-dns-node-cache-"${ARCH}":"${TAG}"
echo "Verifying dnsmasq-nanny startup"
docker run --rm -it --entrypoint=/dnsmasq-nanny "${REGISTRY}"/k8s-dns-dnsmasq-nanny-"${ARCH}":"${TAG}"
echo "Verifying kube-dns startup"
docker run --rm -it --entrypoint=/kube-dns "${REGISTRY}"/k8s-dns-kube-dns-"${ARCH}":"${TAG}"
echo "Verifying sidecar startup"
docker run --rm -it --entrypoint=/sidecar "${REGISTRY}"/k8s-dns-sidecar-"${ARCH}":"${TAG}"
